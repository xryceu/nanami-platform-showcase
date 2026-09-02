package gatewayruntime

import (
	"context"
	"net"
	"testing"
)

func TestAnalyzeEndpointExposureFlagsCloudflareEdgeDNS(t *testing.T) {
	exposure := AnalyzeEndpointExposure(
		context.Background(),
		"gateway-daemon-test.example.com:51820",
		func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("188.114.96.6")}, nil
		},
	)

	if !exposure.CloudflareProxySuspected {
		t.Fatalf("expected Cloudflare proxy suspicion, got %#v", exposure)
	}
	if exposure.Mode != EndpointExposureCloudflareProxyLikely {
		t.Fatalf("expected proxied exposure mode, got %q", exposure.Mode)
	}
	if exposure.Warning == "" {
		t.Fatalf("expected diagnostic warning for proxied endpoint")
	}
}

func TestAnalyzeEndpointExposureKeepsDirectOriginAddress(t *testing.T) {
	exposure := AnalyzeEndpointExposure(context.Background(), "198.51.100.10:51820", nil)
	if exposure.CloudflareProxySuspected {
		t.Fatalf("did not expect direct test-net endpoint to be flagged: %#v", exposure)
	}
	if exposure.Mode != EndpointExposureDirectOrUnknown {
		t.Fatalf("expected direct exposure mode, got %q", exposure.Mode)
	}
}

func TestExtractEndpointHost(t *testing.T) {
	cases := map[string]string{
		"gateway.example.com:51820": "gateway.example.com",
		"[2001:db8::1]:51820":       "2001:db8::1",
		"198.51.100.10:51820":       "198.51.100.10",
		"gateway.example.com":       "gateway.example.com",
	}
	for input, want := range cases {
		if got := ExtractEndpointHost(input); got != want {
			t.Fatalf("ExtractEndpointHost(%q) = %q, want %q", input, got, want)
		}
	}
}
