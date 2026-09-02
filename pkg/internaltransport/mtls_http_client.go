package internaltransport

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var ErrMissingMTLSConfig = errors.New("internal mTLS configuration is missing")

// MTLSHTTPClientConfig contains TLS material for internal mTLS HTTP clients.
type MTLSHTTPClientConfig struct {
	CAFile         string
	ClientCertFile string
	ClientKeyFile  string
	ServerName     string
	Timeout        time.Duration
}

// MaxInternalClientCertLifetime enforces short-lived internal mTLS client certificates.
const MaxInternalClientCertLifetime = 30 * 24 * time.Hour

// NewMTLSHTTPClient constructs hardened mTLS HTTP client used by internal control-plane clients.
func NewMTLSHTTPClient(cfg MTLSHTTPClientConfig) (*http.Client, error) {
	caPath := strings.TrimSpace(cfg.CAFile)
	certPath := strings.TrimSpace(cfg.ClientCertFile)
	keyPath := strings.TrimSpace(cfg.ClientKeyFile)
	if caPath == "" || certPath == "" || keyPath == "" {
		return nil, ErrMissingMTLSConfig
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	transport, err := newReloadingMTLSRoundTripper(MTLSHTTPClientConfig{
		CAFile:         caPath,
		ClientCertFile: certPath,
		ClientKeyFile:  keyPath,
		ServerName:     strings.TrimSpace(cfg.ServerName),
		Timeout:        timeout,
	})
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}, nil
}

type mtlsBundleSnapshot struct {
	caModTime   time.Time
	caSize      int64
	certModTime time.Time
	certSize    int64
	keyModTime  time.Time
	keySize     int64
}

func snapshotMTLSBundle(cfg MTLSHTTPClientConfig) (mtlsBundleSnapshot, error) {
	caInfo, err := os.Stat(cfg.CAFile)
	if err != nil {
		return mtlsBundleSnapshot{}, fmt.Errorf("stat internal mTLS CA: %w", err)
	}
	certInfo, err := os.Stat(cfg.ClientCertFile)
	if err != nil {
		return mtlsBundleSnapshot{}, fmt.Errorf("stat internal mTLS client certificate: %w", err)
	}
	keyInfo, err := os.Stat(cfg.ClientKeyFile)
	if err != nil {
		return mtlsBundleSnapshot{}, fmt.Errorf("stat internal mTLS client key: %w", err)
	}
	return mtlsBundleSnapshot{
		caModTime:   caInfo.ModTime().UTC(),
		caSize:      caInfo.Size(),
		certModTime: certInfo.ModTime().UTC(),
		certSize:    certInfo.Size(),
		keyModTime:  keyInfo.ModTime().UTC(),
		keySize:     keyInfo.Size(),
	}, nil
}

func (s mtlsBundleSnapshot) equal(other mtlsBundleSnapshot) bool {
	return s.caModTime.Equal(other.caModTime) &&
		s.caSize == other.caSize &&
		s.certModTime.Equal(other.certModTime) &&
		s.certSize == other.certSize &&
		s.keyModTime.Equal(other.keyModTime) &&
		s.keySize == other.keySize
}

func loadInternalMTLSRoots(path string) (*x509.CertPool, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read internal mTLS CA: %w", err)
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(pemData); !ok {
		return nil, fmt.Errorf("parse internal mTLS CA: no certificates found")
	}
	return roots, nil
}

func loadInternalMTLSClientCertificate(certPath, keyPath string) (tls.Certificate, error) {
	clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("load internal mTLS client keypair: %w", err)
	}
	if len(clientCert.Certificate) == 0 {
		return tls.Certificate{}, fmt.Errorf("load internal mTLS client keypair: certificate chain is empty")
	}
	leaf, err := x509.ParseCertificate(clientCert.Certificate[0])
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("parse internal mTLS client certificate: %w", err)
	}
	if !leaf.NotAfter.After(leaf.NotBefore) {
		return tls.Certificate{}, fmt.Errorf("internal mTLS client certificate has invalid lifetime")
	}
	if leaf.NotAfter.Sub(leaf.NotBefore) > MaxInternalClientCertLifetime {
		return tls.Certificate{}, fmt.Errorf("internal mTLS client certificate lifetime exceeds max allowed duration")
	}
	clientCert.Leaf = leaf
	return clientCert, nil
}

func buildMTLSTransport(cfg MTLSHTTPClientConfig) (*http.Transport, error) {
	roots, err := loadInternalMTLSRoots(cfg.CAFile)
	if err != nil {
		return nil, err
	}
	clientCert, err := loadInternalMTLSClientCertificate(cfg.ClientCertFile, cfg.ClientKeyFile)
	if err != nil {
		return nil, err
	}

	tlsCfg := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      roots,
		Certificates: []tls.Certificate{clientCert},
	}
	if cfg.ServerName != "" {
		tlsCfg.ServerName = cfg.ServerName
	}

	return &http.Transport{
		TLSClientConfig:       tlsCfg,
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}

type reloadingMTLSTransportState struct {
	snapshot  mtlsBundleSnapshot
	transport *http.Transport
}

type reloadingMTLSRoundTripper struct {
	cfg MTLSHTTPClientConfig

	mu    sync.RWMutex
	state reloadingMTLSTransportState
}

func newReloadingMTLSRoundTripper(cfg MTLSHTTPClientConfig) (*reloadingMTLSRoundTripper, error) {
	state, err := buildReloadingMTLSTransportState(cfg)
	if err != nil {
		return nil, err
	}
	return &reloadingMTLSRoundTripper{
		cfg:   cfg,
		state: state,
	}, nil
}

func buildReloadingMTLSTransportState(cfg MTLSHTTPClientConfig) (reloadingMTLSTransportState, error) {
	snapshot, err := snapshotMTLSBundle(cfg)
	if err != nil {
		return reloadingMTLSTransportState{}, err
	}
	transport, err := buildMTLSTransport(cfg)
	if err != nil {
		return reloadingMTLSTransportState{}, err
	}
	return reloadingMTLSTransportState{
		snapshot:  snapshot,
		transport: transport,
	}, nil
}

func (r *reloadingMTLSRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport, err := r.currentTransport()
	if err != nil {
		return nil, err
	}
	return transport.RoundTrip(req)
}

func (r *reloadingMTLSRoundTripper) CloseIdleConnections() {
	r.mu.RLock()
	transport := r.state.transport
	r.mu.RUnlock()
	if transport != nil {
		transport.CloseIdleConnections()
	}
}

func (r *reloadingMTLSRoundTripper) currentTransport() (*http.Transport, error) {
	snapshot, err := snapshotMTLSBundle(r.cfg)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	current := r.state
	r.mu.RUnlock()
	if current.transport != nil && current.snapshot.equal(snapshot) {
		return current.transport, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.state.transport != nil && r.state.snapshot.equal(snapshot) {
		return r.state.transport, nil
	}

	nextState, err := buildReloadingMTLSTransportState(r.cfg)
	if err != nil {
		return nil, err
	}
	previous := r.state.transport
	r.state = nextState
	if previous != nil {
		previous.CloseIdleConnections()
	}
	return r.state.transport, nil
}
