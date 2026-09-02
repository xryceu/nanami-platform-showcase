package internaltransport

import (
	"strings"
	"testing"
)

func FuzzResolveAndValidateControlPlaneInternalURL(f *testing.F) {
	seeds := []struct {
		url                 string
		protocol            string
		host                string
		port                string
		authMode            string
		requireExplicitPort bool
	}{
		{url: "https://control-plane.internal:18443", authMode: "mtls", requireExplicitPort: true},
		{url: "https://control-plane.internal", authMode: "mtls", requireExplicitPort: false},
		{protocol: "https", host: "control-plane.internal", port: "18443", authMode: "mtls", requireExplicitPort: true},
		{protocol: "http", host: "127.0.0.1", port: "8080", authMode: "token", requireExplicitPort: true},
		{url: "https://control-plane.internal/", authMode: "mtls", requireExplicitPort: false},
		{url: "https://cp.internal?query=nope", authMode: "mtls", requireExplicitPort: false},
	}

	for _, seed := range seeds {
		f.Add(seed.url, seed.protocol, seed.host, seed.port, seed.authMode, seed.requireExplicitPort)
	}

	f.Fuzz(func(t *testing.T, rawURL, protocol, host, port, authMode string, requireExplicitPort bool) {
		resolved, err := ResolveControlPlaneInternalURL(ResolveControlPlaneInternalURLInput{
			URL:      rawURL,
			Protocol: protocol,
			Host:     host,
			Port:     port,
		})
		if err != nil {
			return
		}

		normalized, err := NormalizeControlPlaneInternalURL(resolved)
		if err != nil {
			t.Fatalf("resolved URL %q must remain normalizable: %v", resolved, err)
		}
		if normalized != resolved {
			t.Fatalf("resolved URL must already be normalized: got %q want %q", resolved, normalized)
		}

		parsed, err := ParseControlPlaneInternalURL(resolved)
		if err != nil {
			t.Fatalf("resolved URL %q must remain parseable: %v", resolved, err)
		}

		validateErr := ValidateControlPlaneInternalURLForAuthMode(resolved, authMode, requireExplicitPort)
		if validateErr != nil {
			mode := strings.ToLower(strings.TrimSpace(authMode))
			allowsError := (mode == "mtls" && !strings.EqualFold(parsed.Scheme, "https")) ||
				(requireExplicitPort && parsed.Port() == "")
			if !allowsError {
				t.Fatalf("unexpected validation error for %q authMode=%q requireExplicitPort=%t: %v", resolved, authMode, requireExplicitPort, validateErr)
			}
		}

		listenAddr, listenErr := ControlPlaneInternalListenAddrFromURL(resolved)
		if parsed.Port() == "" {
			if listenErr == nil {
				t.Fatalf("expected listen addr to require explicit port for %q", resolved)
			}
			return
		}
		if listenErr != nil {
			t.Fatalf("expected explicit port to yield listen addr for %q: %v", resolved, listenErr)
		}
		if want := ":" + parsed.Port(); listenAddr != want {
			t.Fatalf("unexpected listen addr for %q: got %q want %q", resolved, listenAddr, want)
		}
	})
}
