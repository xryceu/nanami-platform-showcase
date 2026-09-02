package gatewayruntime

import (
	"context"
	"errors"
	"net"
	"testing"
)

func TestResolveEndpointDNSMode(t *testing.T) {
	if got := ResolveEndpointDNSMode(EndpointDNSModeManual); got != EndpointDNSModeManual {
		t.Fatalf("manual mode should remain manual, got %q", got)
	}
	if got := ResolveEndpointDNSMode(EndpointDNSModeCloudflare); got != EndpointDNSModeCloudflare {
		t.Fatalf("cloudflare mode should remain cloudflare, got %q", got)
	}
	if got := ResolveEndpointDNSMode(""); got != EndpointDNSModeManual {
		t.Fatalf("empty mode should default to manual, got %q", got)
	}
	if got := ResolveEndpointDNSMode("unsupported"); got != EndpointDNSModeManual {
		t.Fatalf("unsupported value should fall back to manual, got %q", got)
	}
}

func TestBuildGatewayRuntimeEndpointManualFlatTemplateDNSOnly(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:         EndpointDNSModeManual,
		BaseDomain:   "example.test",
		GatewayID:    "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		Environment:  "s-nanami",
		PublicIPv4:   "203.0.113.10",
		EndpointPort: 51820,
	}, resolverFor("203.0.113.10"))

	if state.Mode != EndpointDNSModeManual {
		t.Fatalf("expected manual mode, got %q", state.Mode)
	}
	if state.Managed {
		t.Fatalf("manual mode must not be managed")
	}
	if state.Hostname != "gateway-daemon-9dd1cd6b312f-s-nanami.example.test" {
		t.Fatalf("unexpected generated hostname %q", state.Hostname)
	}
	if state.DNSRecordName != "gateway-daemon-9dd1cd6b312f-s-nanami" {
		t.Fatalf("unexpected DNS record name %q", state.DNSRecordName)
	}
	if state.Status != EndpointStatusDNSOnlyDirect {
		t.Fatalf("expected DNS-only direct status, got %#v", state)
	}
}

func TestBuildGatewayRuntimeEndpointManualMissingPublicIPReportsExactIssue(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:        EndpointDNSModeManual,
		BaseDomain:  "example.test",
		GatewayID:   "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		Environment: "s-nanami",
	}, resolverShouldNotRun(t))

	if state.Status != EndpointStatusPublicIPMissing {
		t.Fatalf("expected public IP missing status, got %#v", state)
	}
	if state.StatusMessage != "Gateway public IP is missing; cannot create runtime endpoint DNS record" {
		t.Fatalf("unexpected status message %q", state.StatusMessage)
	}
}

func TestBuildGatewayRuntimeEndpointAllowsValidatedManualPublicIPWithoutDNS(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		PublicIPv4: "203.0.113.10",
	}, resolverShouldNotRun(t))

	if state.Status != EndpointStatusDNSOnlyDirect || state.Hostname != "203.0.113.10" {
		t.Fatalf("expected direct validated public-IP endpoint, got %#v", state)
	}
	if state.DNSRecordName != "" || state.Managed {
		t.Fatalf("direct public-IP endpoint must not claim managed DNS, got %#v", state)
	}
}

func TestBuildGatewayRuntimeEndpointTreatsLiteralIPTemplateAsDirectEndpoint(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:         EndpointDNSModeManual,
		BaseDomain:   "nanami.test",
		HostTemplate: "203.0.113.10",
		PublicIPv4:   "203.0.113.10",
	}, resolverShouldNotRun(t))

	if state.Status != EndpointStatusDNSOnlyDirect || state.Hostname != "203.0.113.10" {
		t.Fatalf("expected literal IP template to remain direct, got %#v", state)
	}
	if state.DNSRecordName != "" {
		t.Fatalf("literal IP endpoint must not produce a DNS record name, got %q", state.DNSRecordName)
	}
}

func TestBuildGatewayRuntimeEndpointRejectsPrivateManualOriginWithoutDNS(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		PublicIPv4: "10.0.0.10",
	}, resolverShouldNotRun(t))

	if state.Status != EndpointStatusPublicIPMissing || state.Hostname != "" {
		t.Fatalf("expected private direct endpoint to fail closed, got %#v", state)
	}
}

func TestBuildGatewayRuntimeEndpointRejectsPrivateManualOriginIP(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		BaseDomain: "example.test",
		GatewayID:  "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		PublicIPv4: "10.0.0.10",
	}, resolverShouldNotRun(t))

	if state.Status != EndpointStatusPublicIPMissing {
		t.Fatalf("expected public IP missing/rejected status, got %#v", state)
	}
	if state.StatusMessage == "" {
		t.Fatalf("expected private IP rejection message")
	}
}

func TestBuildGatewayRuntimeEndpointFlagsCloudflareProxiedEndpoint(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:         EndpointDNSModeManual,
		BaseDomain:   "example.test",
		GatewayID:    "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		Environment:  "s-nanami",
		PublicIPv4:   "203.0.113.10",
		EndpointPort: 51820,
	}, resolverFor("188.114.96.6"))

	if state.Status != EndpointStatusProxiedProtocolUnsupported {
		t.Fatalf("expected proxied protocol unsupported status, got %#v", state)
	}
	if state.ExposureMode != EndpointExposureProxied {
		t.Fatalf("expected proxied exposure, got %q", state.ExposureMode)
	}
}

func TestBuildGatewayRuntimeEndpointIPv6OnlyWhenEnabled(t *testing.T) {
	disabled := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		BaseDomain: "example.test",
		GatewayID:  "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		PublicIPv6: "2001:db8::10",
	}, resolverShouldNotRun(t))
	if disabled.Status != EndpointStatusPublicIPMissing {
		t.Fatalf("expected IPv6 to be ignored unless enabled, got %#v", disabled)
	}

	enabled := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		BaseDomain: "example.test",
		GatewayID:  "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		PublicIPv6: "2001:4860:4860::8888",
		AllowIPv6:  true,
	}, resolverFor("2001:4860:4860::8888"))
	if enabled.Status != EndpointStatusDNSOnlyDirect {
		t.Fatalf("expected enabled IPv6 endpoint to validate, got %#v", enabled)
	}
	if enabled.ExpectedOriginIPv6 != "2001:4860:4860::8888" {
		t.Fatalf("expected IPv6 origin to be kept when enabled, got %q", enabled.ExpectedOriginIPv6)
	}
}

func TestBuildGatewayRuntimeEndpointDetectsWrongOriginIP(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		BaseDomain: "example.test",
		GatewayID:  "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		PublicIPv4: "203.0.113.10",
	}, resolverFor("203.0.113.11"))

	if state.Status != EndpointStatusWrongIP {
		t.Fatalf("expected wrong IP status, got %#v", state)
	}
}

func TestBuildGatewayRuntimeEndpointHandlesMissingRecord(t *testing.T) {
	state := BuildGatewayRuntimeEndpoint(context.Background(), EndpointDNSConfig{
		Mode:       EndpointDNSModeManual,
		BaseDomain: "example.test",
		GatewayID:  "9dd1cd6b-312f-4cb2-9ec3-60fd78c30758",
		PublicIPv4: "203.0.113.10",
	}, func(context.Context, string) ([]net.IP, error) {
		return nil, errors.New("no such host")
	})

	if state.Status != EndpointStatusMissingRecord {
		t.Fatalf("expected missing record, got %#v", state)
	}
	if state.ExposureMode != EndpointExposureUnresolved {
		t.Fatalf("expected unresolved exposure, got %q", state.ExposureMode)
	}
}

func TestAttachProtocolEndpointsAddsWireGuardAndGRE(t *testing.T) {
	state := AttachProtocolEndpoints(GatewayRuntimeEndpoint{
		Hostname: "gateway-daemon-test.example.test",
	}, 51820, "10.0.112.4")

	if len(state.Protocols) != 2 {
		t.Fatalf("expected wireguard and GRE protocol endpoints, got %#v", state.Protocols)
	}
	if state.Protocols[0].Protocol != "wireguard" || state.Protocols[0].Endpoint != "gateway-daemon-test.example.test:51820" {
		t.Fatalf("unexpected WireGuard protocol endpoint: %#v", state.Protocols[0])
	}
	if state.Protocols[0].HealthMethod != "ICMP" {
		t.Fatalf("expected ICMP health method for IP target, got %q", state.Protocols[0].HealthMethod)
	}
	if state.Protocols[1].Protocol != "gre" || state.Protocols[1].Port != nil {
		t.Fatalf("unexpected GRE protocol endpoint: %#v", state.Protocols[1])
	}
}

func resolverFor(values ...string) EndpointResolver {
	return func(context.Context, string) ([]net.IP, error) {
		out := make([]net.IP, 0, len(values))
		for _, value := range values {
			out = append(out, net.ParseIP(value))
		}
		return out, nil
	}
}

func resolverShouldNotRun(t *testing.T) EndpointResolver {
	t.Helper()
	return func(context.Context, string) ([]net.IP, error) {
		t.Fatalf("resolver should not run")
		return nil, nil
	}
}
