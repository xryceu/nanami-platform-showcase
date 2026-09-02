package gatewayruntime

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	EndpointDNSModeManual     = "manual"
	EndpointDNSModeCloudflare = "cloudflare"

	EndpointDNSProviderManual     = "manual"
	EndpointDNSProviderCloudflare = "cloudflare"

	EndpointStatusDNSOnlyDirect              = "endpoint_dns_only_direct"
	EndpointStatusProxiedProtocolUnsupported = "endpoint_proxied_protocol_unsupported"
	EndpointStatusMissingRecord              = "endpoint_missing_record"
	EndpointStatusWrongIP                    = "endpoint_wrong_ip"
	EndpointStatusPublicIPMissing            = "endpoint_public_ip_missing"
	EndpointStatusDNSSyncFailed              = "endpoint_dns_sync_failed"
	EndpointStatusIPv6Mismatch               = "endpoint_ipv6_mismatch"
	EndpointStatusNotConfigured              = "endpoint_not_configured"
	EndpointStatusUnknown                    = "endpoint_unknown"

	ProtocolStatusWireGuardUDPBlocked = "wireguard_udp_blocked"
	ProtocolStatusGREBlocked          = "gre_protocol_blocked"
	ProtocolStatusTCPPortBlocked      = "tcp_port_blocked"
	ProtocolStatusHealthTargetFailed  = "health_target_failed"

	EndpointExposureDNSOnlyDirect = "dns_only_direct"
	EndpointExposureProxied       = "proxied"
	EndpointExposureUnresolved    = "unresolved"
	EndpointExposureUnknown       = "unknown"

	DefaultEndpointDNSTTL              = 60
	DefaultGatewayEndpointHostTemplate = "gateway-daemon-{gateway_id_short}-{env}"
)

// EndpointDNSConfig describes how a gateway endpoint should be published.
type EndpointDNSConfig struct {
	Mode           string
	BaseDomain     string
	HostTemplate   string
	TTL            int
	AllowIPv6      bool
	DeleteOnDelete bool
	PublicIPv4     string
	PublicIPv6     string
	EndpointPort   int
	Environment    string
	GatewayID      string
	TenantID       string
	NetworkID      string
	Region         string
}

// EndpointProtocol describes one protocol exposed through a gateway endpoint.
type EndpointProtocol struct {
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Transport    string `json:"transport,omitempty"`
	Port         *int   `json:"port,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
	HealthTarget string `json:"healthTarget,omitempty"`
	HealthMethod string `json:"healthMethod,omitempty"`
	Status       string `json:"status,omitempty"`
}

// GatewayRuntimeEndpoint contains status data that can be returned to clients.
// Provider credentials never enter this model.
type GatewayRuntimeEndpoint struct {
	Mode               string             `json:"mode"`
	Owner              string             `json:"owner,omitempty"`
	Managed            bool               `json:"managed"`
	Provider           string             `json:"provider,omitempty"`
	Hostname           string             `json:"hostname,omitempty"`
	DNSRecordName      string             `json:"dnsRecordName,omitempty"`
	BaseDomain         string             `json:"baseDomain,omitempty"`
	HostTemplate       string             `json:"hostTemplate,omitempty"`
	TTL                int                `json:"ttl,omitempty"`
	ExpectedOriginIPv4 string             `json:"expectedOriginIPv4,omitempty"`
	ExpectedOriginIPv6 string             `json:"expectedOriginIPv6,omitempty"`
	ResolvedIPv4       []string           `json:"resolvedIPv4,omitempty"`
	ResolvedIPv6       []string           `json:"resolvedIPv6,omitempty"`
	Status             string             `json:"status,omitempty"`
	StatusMessage      string             `json:"statusMessage,omitempty"`
	ExposureMode       string             `json:"exposureMode,omitempty"`
	Action             string             `json:"action,omitempty"`
	SyncStatus         string             `json:"dnsSyncStatus,omitempty"`
	SyncAction         string             `json:"dnsSyncAction,omitempty"`
	LastSyncAt         string             `json:"dnsLastSyncAt,omitempty"`
	LastSyncError      string             `json:"dnsLastSyncError,omitempty"`
	LastCheckedAt      string             `json:"lastCheckedAt,omitempty"`
	Protocols          []EndpointProtocol `json:"protocols,omitempty"`
}

// ResolveEndpointDNSMode falls back to manual mode for unknown values.
func ResolveEndpointDNSMode(explicit string) string {
	mode := strings.ToLower(strings.TrimSpace(explicit))
	switch mode {
	case EndpointDNSModeManual, EndpointDNSModeCloudflare:
		return mode
	case "":
	default:
		return EndpointDNSModeManual
	}
	return EndpointDNSModeManual
}

// BuildGatewayRuntimeEndpoint compares configured origins with current DNS.
func BuildGatewayRuntimeEndpoint(ctx context.Context, cfg EndpointDNSConfig, resolver EndpointResolver) GatewayRuntimeEndpoint {
	mode := ResolveEndpointDNSMode(cfg.Mode)
	ttl := cfg.TTL
	if ttl <= 0 {
		ttl = DefaultEndpointDNSTTL
	}
	provider := endpointProviderForMode(mode)
	out := GatewayRuntimeEndpoint{
		Mode:          mode,
		Owner:         endpointOwnerForMode(mode),
		Managed:       mode == EndpointDNSModeCloudflare,
		Provider:      provider,
		BaseDomain:    normalizeDomain(cfg.BaseDomain),
		HostTemplate:  strings.TrimSpace(cfg.HostTemplate),
		TTL:           ttl,
		LastCheckedAt: time.Now().UTC().Format(time.RFC3339),
		ExposureMode:  EndpointExposureUnknown,
	}
	if out.HostTemplate == "" {
		out.HostTemplate = DefaultGatewayEndpointHostTemplate
	}
	if out.BaseDomain == "" && mode == EndpointDNSModeManual {
		ipv4, err := normalizePublicOriginIPString(cfg.PublicIPv4)
		if err != nil {
			out.Status = EndpointStatusPublicIPMissing
			out.StatusMessage = err.Error()
			out.Action = missingGatewayPublicIPMessage()
			return out
		}
		ipv6 := ""
		if cfg.AllowIPv6 {
			ipv6, err = normalizePublicOriginIPString(cfg.PublicIPv6)
			if err != nil {
				out.Status = EndpointStatusIPv6Mismatch
				out.StatusMessage = err.Error()
				out.Action = "Set NANAMI_GATEWAY_PUBLIC_IPV6 to a public IPv6 address or disable IPv6 endpoint publication."
				return out
			}
		}
		out.Hostname = ipv4
		if out.Hostname == "" {
			out.Hostname = ipv6
		}
		if out.Hostname == "" {
			out.Status = EndpointStatusNotConfigured
			out.StatusMessage = "NANAMI_DNS_BASE_DOMAIN or a public gateway origin IP is required"
			out.Action = "Configure a DNS base domain or NANAMI_GATEWAY_PUBLIC_IP."
			return out
		}
		out.ExpectedOriginIPv4 = ipv4
		out.ExpectedOriginIPv6 = ipv6
		if ipv4 != "" {
			out.ResolvedIPv4 = appendUniqueString(out.ResolvedIPv4, ipv4)
		}
		if ipv6 != "" {
			out.ResolvedIPv6 = appendUniqueString(out.ResolvedIPv6, ipv6)
		}
		out.Status = EndpointStatusDNSOnlyDirect
		out.ExposureMode = EndpointExposureDNSOnlyDirect
		out.Action = "Gateway endpoint uses the configured public origin address directly."
		return out
	}

	host, hostErr := endpointHostname(cfg, out.HostTemplate, out.BaseDomain)
	if hostErr != nil {
		out.Status = EndpointStatusNotConfigured
		out.StatusMessage = hostErr.Error()
		out.Action = "Fix gateway endpoint host/base-domain configuration."
		return out
	}
	out.Hostname = host
	if net.ParseIP(host) == nil {
		out.DNSRecordName = DNSRecordName(host, out.BaseDomain)
	}

	var originErr error
	out.ExpectedOriginIPv4, originErr = normalizePublicOriginIPString(cfg.PublicIPv4)
	if originErr != nil {
		out.Status = EndpointStatusPublicIPMissing
		out.StatusMessage = originErr.Error()
		out.Action = missingGatewayPublicIPMessage()
		return out
	}
	if cfg.AllowIPv6 {
		out.ExpectedOriginIPv6, originErr = normalizePublicOriginIPString(cfg.PublicIPv6)
		if originErr != nil {
			out.Status = EndpointStatusIPv6Mismatch
			out.StatusMessage = originErr.Error()
			out.Action = "Set NANAMI_GATEWAY_PUBLIC_IPV6 to the gateway origin public IPv6 address or disable NANAMI_DNS_ALLOW_IPV6."
			return out
		}
	}
	if ip := net.ParseIP(out.Hostname); ip != nil {
		if !isPublicOriginIP(ip) {
			out.Status = EndpointStatusPublicIPMissing
			out.StatusMessage = fmt.Sprintf("%s must be a public gateway origin IP address", ip.String())
			out.Action = missingGatewayPublicIPMessage()
			return out
		}
		if ip.To4() != nil {
			if out.ExpectedOriginIPv4 == "" {
				out.ExpectedOriginIPv4 = ip.String()
			}
			out.ResolvedIPv4 = appendUniqueString(out.ResolvedIPv4, ip.String())
		} else {
			if out.ExpectedOriginIPv6 == "" {
				out.ExpectedOriginIPv6 = ip.String()
			}
			out.ResolvedIPv6 = appendUniqueString(out.ResolvedIPv6, ip.String())
		}
		out.Status = EndpointStatusDNSOnlyDirect
		out.ExposureMode = EndpointExposureDNSOnlyDirect
		out.Action = "Gateway endpoint uses the configured public origin address."
		return out
	}
	if out.ExpectedOriginIPv4 == "" && out.ExpectedOriginIPv6 == "" {
		out.Status = EndpointStatusPublicIPMissing
		out.StatusMessage = missingGatewayPublicIPMessage()
		out.Action = missingGatewayPublicIPMessage()
		return out
	}

	if resolver == nil {
		resolver = defaultEndpointResolver
	}
	resolved, err := resolver(ctx, out.Hostname)
	if err != nil || len(resolved) == 0 {
		out.Status = EndpointStatusMissingRecord
		if err != nil {
			out.StatusMessage = err.Error()
		}
		out.ExposureMode = EndpointExposureUnresolved
		out.Action = manualOrManagedAction(mode, "Create a DNS-only A record for the gateway hostname.")
		return out
	}
	for _, ip := range resolved {
		if ip == nil {
			continue
		}
		if IsCloudflareEdgeIP(ip) {
			out.ExposureMode = EndpointExposureProxied
			out.Status = EndpointStatusProxiedProtocolUnsupported
			out.Action = "Move the gateway endpoint to DNS-only/direct DNS or a UDP-capable proxy."
		}
		if ip.To4() != nil {
			out.ResolvedIPv4 = appendUniqueString(out.ResolvedIPv4, ip.String())
		} else {
			out.ResolvedIPv6 = appendUniqueString(out.ResolvedIPv6, ip.String())
		}
	}
	if out.Status == EndpointStatusProxiedProtocolUnsupported {
		return out
	}
	out.ExposureMode = EndpointExposureDNSOnlyDirect
	if out.ExpectedOriginIPv4 != "" && !containsEndpointString(out.ResolvedIPv4, out.ExpectedOriginIPv4) {
		out.Status = EndpointStatusWrongIP
		out.Action = manualOrManagedAction(mode, "Update gateway DNS A record to the expected origin IP.")
		return out
	}
	if cfg.AllowIPv6 && out.ExpectedOriginIPv6 != "" && !containsEndpointString(out.ResolvedIPv6, out.ExpectedOriginIPv6) {
		out.Status = EndpointStatusIPv6Mismatch
		out.Action = manualOrManagedAction(mode, "Update gateway DNS AAAA record to the expected origin IPv6.")
		return out
	}
	out.Status = EndpointStatusDNSOnlyDirect
	out.Action = "Gateway endpoint DNS resolves directly to the expected origin."
	return out
}

func AttachProtocolEndpoints(endpoint GatewayRuntimeEndpoint, port int, healthTarget string) GatewayRuntimeEndpoint {
	if port <= 0 || port > 65535 {
		port = 51820
	}
	host := strings.TrimSpace(endpoint.Hostname)
	wgEndpoint := ""
	if host != "" {
		wgEndpoint = net.JoinHostPort(host, strconv.Itoa(port))
	}
	healthMethod := "ICMP"
	if strings.Contains(strings.TrimSpace(healthTarget), "://") || strings.Contains(strings.TrimSpace(healthTarget), ":") {
		healthMethod = "TCP/HTTP"
	}
	endpoint.Protocols = append(endpoint.Protocols, EndpointProtocol{
		Name:         "wireguard",
		Protocol:     "wireguard",
		Transport:    "udp",
		Port:         intPtr(port),
		Endpoint:     wgEndpoint,
		HealthTarget: strings.TrimSpace(healthTarget),
		HealthMethod: healthMethod,
	})
	endpoint.Protocols = append(endpoint.Protocols, EndpointProtocol{
		Name:      "gre",
		Protocol:  "gre",
		Transport: "ip-protocol-47",
	})
	return endpoint
}

func DNSRecordName(hostname, baseDomain string) string {
	host := normalizeDomain(hostname)
	base := normalizeDomain(baseDomain)
	if host == "" || base == "" {
		return host
	}
	if host == base {
		return "@"
	}
	suffix := "." + base
	if strings.HasSuffix(host, suffix) {
		return strings.TrimSuffix(host, suffix)
	}
	return host
}

func endpointHostname(cfg EndpointDNSConfig, template string, baseDomain string) (string, error) {
	base := normalizeDomain(baseDomain)
	if base == "" {
		return "", errors.New("NANAMI_DNS_BASE_DOMAIN is required for generated gateway endpoint hosts")
	}
	label := ApplyEndpointHostTemplate(template, EndpointTemplateContext{
		GatewayID: cfg.GatewayID,
		TenantID:  cfg.TenantID,
		NetworkID: cfg.NetworkID,
		Region:    cfg.Region,
		Env:       cfg.Environment,
	})
	label = strings.Trim(label, ".")
	if label == "" {
		return "", errors.New("gateway endpoint host template resolved to an empty host")
	}
	if ip := net.ParseIP(label); ip != nil {
		return ip.String(), nil
	}
	if strings.HasSuffix(normalizeDomain(label), "."+base) || normalizeDomain(label) == base {
		return normalizeDomain(label), validateHostnameInsideBase(label, base)
	}
	return normalizeDomain(label + "." + base), nil
}

// EndpointTemplateContext contains values available to endpoint host templates.
type EndpointTemplateContext struct {
	GatewayID string
	TenantID  string
	NetworkID string
	Region    string
	Env       string
}

func ApplyEndpointHostTemplate(template string, ctx EndpointTemplateContext) string {
	value := strings.TrimSpace(template)
	if value == "" {
		value = DefaultGatewayEndpointHostTemplate
	}
	replacements := map[string]string{
		"{gateway_id}":       sanitizeDNSSegment(ctx.GatewayID),
		"{gateway_id_short}": gatewayIDShort(ctx.GatewayID),
		"{tenant_id}":        sanitizeDNSSegment(ctx.TenantID),
		"{network_id}":       sanitizeDNSSegment(ctx.NetworkID),
		"{region}":           sanitizeDNSSegment(ctx.Region),
		"{env}":              sanitizeDNSSegment(ctx.Env),
		"{id}":               gatewayIDShort(ctx.GatewayID),
		"<id>":               gatewayIDShort(ctx.GatewayID),
	}
	for token, replacement := range replacements {
		if replacement == "" {
			replacement = "unknown"
		}
		value = strings.ReplaceAll(value, token, replacement)
	}
	parts := strings.Split(value, ".")
	for i := range parts {
		parts[i] = sanitizeDNSSegment(parts[i])
	}
	return strings.Trim(strings.Join(parts, "."), ".")
}

func validateHostnameInsideBase(hostname, baseDomain string) error {
	host := normalizeDomain(hostname)
	base := normalizeDomain(baseDomain)
	if host == "" {
		return errors.New("gateway endpoint hostname is empty")
	}
	if net.ParseIP(host) != nil || base == "" {
		return nil
	}
	if host == base || strings.HasSuffix(host, "."+base) {
		return nil
	}
	return fmt.Errorf("gateway endpoint hostname %q must be inside base domain %q", host, base)
}

func normalizeDomain(value string) string {
	return strings.Trim(strings.ToLower(strings.TrimSpace(value)), ".")
}

func sanitizeDNSSegment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 63 {
		out = strings.Trim(out[:63], "-")
	}
	return out
}

func gatewayIDShort(value string) string {
	segment := strings.ReplaceAll(sanitizeDNSSegment(value), "-", "")
	if len(segment) > 12 {
		return segment[:12]
	}
	return segment
}

func endpointProviderForMode(mode string) string {
	if mode == EndpointDNSModeCloudflare {
		return EndpointDNSProviderCloudflare
	}
	return EndpointDNSProviderManual
}

func endpointOwnerForMode(mode string) string {
	if mode == EndpointDNSModeCloudflare {
		return "control-plane"
	}
	return "operator"
}

func manualOrManagedAction(mode string, manual string) string {
	if mode == EndpointDNSModeCloudflare {
		return "Cloudflare DNS automation will repair the DNS-only gateway record when provider credentials are valid."
	}
	return manual
}

func normalizeIPString(value string) string {
	ip := net.ParseIP(strings.TrimSpace(value))
	if ip == nil {
		return ""
	}
	return ip.String()
}

func normalizePublicOriginIPString(value string) (string, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return "", nil
	}
	ip := net.ParseIP(raw)
	if ip == nil {
		return "", fmt.Errorf("%s must be a public gateway origin IP address", raw)
	}
	if !isPublicOriginIP(ip) {
		return "", fmt.Errorf("%s must be a public gateway origin IP address", ip.String())
	}
	return ip.String(), nil
}

func isPublicOriginIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return !(ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.IsMulticast())
}

func missingGatewayPublicIPMessage() string {
	return "Gateway public IP is missing; cannot create runtime endpoint DNS record"
}

func appendUniqueString(values []string, next string) []string {
	next = strings.TrimSpace(next)
	if next == "" || containsEndpointString(values, next) {
		return values
	}
	return append(values, next)
}

func containsEndpointString(values []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, value := range values {
		if strings.TrimSpace(value) == target {
			return true
		}
	}
	return false
}

func intPtr(value int) *int {
	return &value
}
