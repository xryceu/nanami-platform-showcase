package internaltransport

import (
	"errors"
	"testing"
)

func TestResolveControlPlaneInternalURLPrefersCanonical(t *testing.T) {
	resolved, err := ResolveControlPlaneInternalURL(ResolveControlPlaneInternalURLInput{
		URL: "https://cp.internal:18443/",
	})
	if err != nil {
		t.Fatalf("ResolveControlPlaneInternalURL returned error: %v", err)
	}
	if resolved != "https://cp.internal:18443" {
		t.Fatalf("expected canonical URL, got %q", resolved)
	}
}

func TestResolveControlPlaneInternalURLErrorsOnAmbiguousInputs(t *testing.T) {
	_, err := ResolveControlPlaneInternalURL(ResolveControlPlaneInternalURLInput{
		URL:      "https://cp.internal:18443",
		Protocol: "https",
		Host:     "cp.internal",
		Port:     "18443",
	})
	if err == nil {
		t.Fatalf("expected error for ambiguous canonical and legacy inputs")
	}
}

func TestResolveControlPlaneInternalURLFromParts(t *testing.T) {
	resolved, err := ResolveControlPlaneInternalURL(ResolveControlPlaneInternalURLInput{
		Protocol: "https",
		Host:     "cp.internal",
		Port:     "18443",
	})
	if err != nil {
		t.Fatalf("ResolveControlPlaneInternalURL returned error: %v", err)
	}
	if resolved != "https://cp.internal:18443" {
		t.Fatalf("expected URL built from parts, got %q", resolved)
	}
}

func TestValidateControlPlaneInternalURLForAuthModeRequiresHTTPSForMTLS(t *testing.T) {
	err := ValidateControlPlaneInternalURLForAuthMode("http://cp.internal:18080", "mtls", false)
	if err == nil {
		t.Fatalf("expected error for http scheme in mtls mode")
	}
}

func TestControlPlaneInternalListenAddrFromURLRequiresPort(t *testing.T) {
	_, err := ControlPlaneInternalListenAddrFromURL("https://cp.internal")
	if err == nil {
		t.Fatalf("expected error for missing explicit port")
	}
}

func TestNewMTLSHTTPClientMissingConfig(t *testing.T) {
	_, err := NewMTLSHTTPClient(MTLSHTTPClientConfig{})
	if !errors.Is(err, ErrMissingMTLSConfig) {
		t.Fatalf("expected ErrMissingMTLSConfig, got %v", err)
	}
}
