package internaltransport

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewMTLSHTTPClientRejectsLongLivedClientCertificate(t *testing.T) {
	tmpDir := t.TempDir()
	caFile, certFile, keyFile := mustWriteClientBundle(t, tmpDir, 90*24*time.Hour)

	_, err := NewMTLSHTTPClient(MTLSHTTPClientConfig{
		CAFile:         caFile,
		ClientCertFile: certFile,
		ClientKeyFile:  keyFile,
	})
	if err == nil {
		t.Fatalf("expected long-lived certificate validation error")
	}
	if !strings.Contains(err.Error(), "lifetime exceeds max allowed duration") {
		t.Fatalf("expected lifetime validation error, got %v", err)
	}
}

func TestNewMTLSHTTPClientAcceptsShortLivedClientCertificate(t *testing.T) {
	tmpDir := t.TempDir()
	caFile, certFile, keyFile := mustWriteClientBundle(t, tmpDir, 24*time.Hour)

	client, err := NewMTLSHTTPClient(MTLSHTTPClientConfig{
		CAFile:         caFile,
		ClientCertFile: certFile,
		ClientKeyFile:  keyFile,
	})
	if err != nil {
		t.Fatalf("expected short-lived certificate to pass validation: %v", err)
	}
	if client == nil {
		t.Fatalf("expected non-nil http client")
	}
}

func TestNewMTLSHTTPClientReloadsRotatedBundle(t *testing.T) {
	tmpDir := t.TempDir()
	caFile := filepath.Join(tmpDir, "ca.pem")
	clientCertFile := filepath.Join(tmpDir, "client.pem")
	clientKeyFile := filepath.Join(tmpDir, "client.key")
	serverCertFile := filepath.Join(tmpDir, "server.pem")
	serverKeyFile := filepath.Join(tmpDir, "server.key")

	bundle1 := mustCreateMTLSBundle(t, time.Now().UTC(), 101)
	mustWriteMTLSBundleFiles(t, caFile, clientCertFile, clientKeyFile, serverCertFile, serverKeyFile, bundle1, time.Now().UTC())

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			t.Fatalf("expected peer certificate on request")
		}
		_, _ = io.WriteString(w, r.TLS.PeerCertificates[0].SerialNumber.Text(16))
	}))
	srv.TLS = &tls.Config{
		MinVersion: tls.VersionTLS12,
		GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
			serverKeypair, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
			if err != nil {
				return nil, err
			}
			caPEM, err := os.ReadFile(caFile)
			if err != nil {
				return nil, err
			}
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(caPEM); !ok {
				return nil, io.ErrUnexpectedEOF
			}
			return &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{serverKeypair},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    pool,
			}, nil
		},
	}
	srv.StartTLS()
	defer srv.Close()

	client, err := NewMTLSHTTPClient(MTLSHTTPClientConfig{
		CAFile:         caFile,
		ClientCertFile: clientCertFile,
		ClientKeyFile:  clientKeyFile,
		ServerName:     "localhost",
	})
	if err != nil {
		t.Fatalf("NewMTLSHTTPClient failed: %v", err)
	}

	firstSerial := mustReadClientSerial(t, client, srv.URL)
	if firstSerial != bundle1.clientSerial.Text(16) {
		t.Fatalf("expected first client serial %s, got %s", bundle1.clientSerial.Text(16), firstSerial)
	}

	bundle2 := mustCreateMTLSBundle(t, time.Now().UTC().Add(time.Minute), 202)
	mustWriteMTLSBundleFiles(t, caFile, clientCertFile, clientKeyFile, serverCertFile, serverKeyFile, bundle2, time.Now().UTC().Add(2*time.Second))

	secondSerial := mustReadClientSerial(t, client, srv.URL)
	if secondSerial != bundle2.clientSerial.Text(16) {
		t.Fatalf("expected rotated client serial %s, got %s", bundle2.clientSerial.Text(16), secondSerial)
	}
}

func mustWriteClientBundle(t *testing.T, dir string, certLifetime time.Duration) (caPath, certPath, keyPath string) {
	t.Helper()

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate CA key: %v", err)
	}
	now := time.Now().UTC()
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Nanami Test CA"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse CA certificate: %v", err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate client key: %v", err)
	}
	notBefore := now.Add(-time.Minute)
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "Nanami Test Client"},
		NotBefore:    notBefore,
		NotAfter:     notBefore.Add(certLifetime),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create client certificate: %v", err)
	}

	caPath = filepath.Join(dir, "ca.pem")
	certPath = filepath.Join(dir, "client.pem")
	keyPath = filepath.Join(dir, "client.key")

	mustWritePEMFile(t, caPath, "CERTIFICATE", caDER, 0o644)
	mustWritePEMFile(t, certPath, "CERTIFICATE", clientDER, 0o644)
	clientKeyDER := x509.MarshalPKCS1PrivateKey(clientKey)
	mustWritePEMFile(t, keyPath, "RSA PRIVATE KEY", clientKeyDER, 0o600)

	return caPath, certPath, keyPath
}

func mustWritePEMFile(t *testing.T, path, pemType string, der []byte, mode os.FileMode) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: pemType, Bytes: der}); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

type testMTLSBundle struct {
	caPEM         []byte
	clientCertPEM []byte
	clientKeyPEM  []byte
	serverCertPEM []byte
	serverKeyPEM  []byte
	clientSerial  *big.Int
}

func mustCreateMTLSBundle(t *testing.T, now time.Time, serialBase int64) testMTLSBundle {
	t.Helper()

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate CA key: %v", err)
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(serialBase),
		Subject:               pkix.Name{CommonName: "Nanami Reload Test CA"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse CA certificate: %v", err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate client key: %v", err)
	}
	clientSerial := big.NewInt(serialBase + 1)
	clientTemplate := &x509.Certificate{
		SerialNumber: clientSerial,
		Subject:      pkix.Name{CommonName: "Nanami Reload Test Client"},
		NotBefore:    now.Add(-time.Minute),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create client certificate: %v", err)
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate server key: %v", err)
	}
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(serialBase + 2),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    now.Add(-time.Minute),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create server certificate: %v", err)
	}

	return testMTLSBundle{
		caPEM:         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}),
		clientCertPEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER}),
		clientKeyPEM:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)}),
		serverCertPEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER}),
		serverKeyPEM:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)}),
		clientSerial:  clientSerial,
	}
}

func mustWriteMTLSBundleFiles(t *testing.T, caFile, clientCertFile, clientKeyFile, serverCertFile, serverKeyFile string, bundle testMTLSBundle, modTime time.Time) {
	t.Helper()

	mustWriteRawFile(t, caFile, bundle.caPEM, 0o644, modTime)
	mustWriteRawFile(t, clientCertFile, bundle.clientCertPEM, 0o644, modTime)
	mustWriteRawFile(t, clientKeyFile, bundle.clientKeyPEM, 0o600, modTime)
	mustWriteRawFile(t, serverCertFile, bundle.serverCertPEM, 0o644, modTime)
	mustWriteRawFile(t, serverKeyFile, bundle.serverKeyPEM, 0o600, modTime)
}

func mustWriteRawFile(t *testing.T, path string, data []byte, mode os.FileMode, modTime time.Time) {
	t.Helper()
	if err := os.WriteFile(path, data, mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes %s: %v", path, err)
	}
}

func mustReadClientSerial(t *testing.T, client *http.Client, url string) string {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("GET %s failed: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}
