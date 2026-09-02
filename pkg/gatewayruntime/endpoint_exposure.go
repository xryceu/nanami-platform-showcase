package gatewayruntime

import (
	"context"
	"net"
	"strconv"
	"strings"
)

const (
	EndpointExposureDirectOrUnknown       = "direct_or_unresolved"
	EndpointExposureCloudflareProxyLikely = "cloudflare_proxy_suspected"
)

// EndpointExposure summarizes how clients reach a WireGuard endpoint.
type EndpointExposure struct {
	Endpoint                 string   `json:"endpoint,omitempty"`
	Host                     string   `json:"host,omitempty"`
	ResolvedIPs              []string `json:"resolvedIPs,omitempty"`
	Mode                     string   `json:"mode,omitempty"`
	CloudflareProxySuspected bool     `json:"cloudflareProxySuspected,omitempty"`
	Warning                  string   `json:"warning,omitempty"`
	ResolutionError          string   `json:"resolutionError,omitempty"`
}

type EndpointResolver func(context.Context, string) ([]net.IP, error)

// ExtractEndpointHost returns the host component from a WireGuard endpoint.
func ExtractEndpointHost(endpoint string) string {
	value := strings.TrimSpace(endpoint)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "[") {
		if closing := strings.Index(value, "]"); closing > 1 {
			return value[1:closing]
		}
	}
	if host, port, err := net.SplitHostPort(value); err == nil && strings.TrimSpace(port) != "" {
		return strings.TrimSpace(host)
	}
	if strings.Count(value, ":") == 1 {
		parts := strings.SplitN(value, ":", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}
	return value
}

// AnalyzeEndpointExposure distinguishes a direct origin from a proxied address.
func AnalyzeEndpointExposure(ctx context.Context, endpoint string, resolver EndpointResolver) EndpointExposure {
	raw := strings.TrimSpace(endpoint)
	out := EndpointExposure{
		Endpoint: raw,
		Host:     ExtractEndpointHost(raw),
		Mode:     EndpointExposureDirectOrUnknown,
	}
	if raw == "" || out.Host == "" {
		return out
	}
	ip := net.ParseIP(out.Host)
	if ip != nil {
		if IsCloudflareEdgeIP(ip) {
			out.CloudflareProxySuspected = true
			out.Mode = EndpointExposureCloudflareProxyLikely
			out.ResolvedIPs = []string{ip.String()}
			out.Warning = "WireGuard endpoint appears to resolve through Cloudflare proxy; WireGuard UDP requires a DNS-only/direct endpoint unless a UDP-capable proxy is configured."
		}
		return out
	}
	if resolver == nil {
		resolver = defaultEndpointResolver
	}
	ips, err := resolver(ctx, out.Host)
	if err != nil {
		out.ResolutionError = err.Error()
		return out
	}
	seen := make(map[string]struct{}, len(ips))
	for _, resolved := range ips {
		if resolved == nil {
			continue
		}
		normalized := resolved.String()
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out.ResolvedIPs = append(out.ResolvedIPs, normalized)
		if IsCloudflareEdgeIP(resolved) {
			out.CloudflareProxySuspected = true
		}
	}
	if out.CloudflareProxySuspected {
		out.Mode = EndpointExposureCloudflareProxyLikely
		out.Warning = "WireGuard endpoint appears to resolve through Cloudflare proxy; WireGuard UDP requires a DNS-only/direct endpoint unless a UDP-capable proxy is configured."
	}
	return out
}

func defaultEndpointResolver(ctx context.Context, host string) ([]net.IP, error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	out := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		out = append(out, addr.IP)
	}
	return out, nil
}

// IsCloudflareEdgeIP returns true for Cloudflare public edge prefixes used by
// proxied DNS records. These ranges are published by Cloudflare and stable
// enough for diagnostics; direct origin IPs are not expected to fall inside
// them.
func IsCloudflareEdgeIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, prefix := range cloudflareEdgeCIDRs {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

var cloudflareEdgeCIDRs = mustParseCIDRs([]string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
	"2400:cb00::/32",
	"2606:4700::/32",
	"2803:f800::/32",
	"2405:b500::/32",
	"2405:8100::/32",
	"2a06:98c0::/29",
	"2c0f:f248::/32",
})

func mustParseCIDRs(values []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		_, ipnet, err := net.ParseCIDR(value)
		if err != nil {
			panic("invalid Cloudflare CIDR " + strconv.Quote(value) + ": " + err.Error())
		}
		out = append(out, ipnet)
	}
	return out
}
