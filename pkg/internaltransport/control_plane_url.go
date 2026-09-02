package internaltransport

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/xryceu/nanami-platform-showcase/pkg/urlutil"
)

const (
	// CanonicalControlPlaneInternalURLKey names the preferred URL setting.
	CanonicalControlPlaneInternalURLKey = "CONTROL_PLANE_INTERNAL_URL"
	// LegacyControlPlaneInternalPrefix is kept for split protocol/host/port settings.
	LegacyControlPlaneInternalPrefix = "CONTROL_PLANE_INTERNAL"
)

// ResolveControlPlaneInternalURLInput accepts either one URL or its legacy parts.
type ResolveControlPlaneInternalURLInput struct {
	URL      string
	Protocol string
	Host     string
	Port     string
}

// ResolveControlPlaneInternalURL rejects ambiguous mixed configuration.
func ResolveControlPlaneInternalURL(input ResolveControlPlaneInternalURLInput) (string, error) {
	rawURL := strings.TrimSpace(input.URL)
	protocol := strings.TrimSpace(input.Protocol)
	host := strings.TrimSpace(input.Host)
	port := strings.TrimSpace(input.Port)

	if rawURL != "" {
		if protocol != "" || host != "" || port != "" {
			return "", fmt.Errorf("%s must not be combined with %s_PROTOCOL/%s_HOST/%s_PORT", CanonicalControlPlaneInternalURLKey, LegacyControlPlaneInternalPrefix, LegacyControlPlaneInternalPrefix, LegacyControlPlaneInternalPrefix)
		}
		return NormalizeControlPlaneInternalURL(rawURL)
	}

	if protocol == "" && host == "" && port == "" {
		return "", fmt.Errorf("%s is required", CanonicalControlPlaneInternalURLKey)
	}
	if protocol == "" || host == "" {
		return "", fmt.Errorf("%s_PROTOCOL and %s_HOST are required when %s is not set", LegacyControlPlaneInternalPrefix, LegacyControlPlaneInternalPrefix, CanonicalControlPlaneInternalURLKey)
	}

	built, err := urlutil.Build(protocol, host, port)
	if err != nil {
		return "", fmt.Errorf("%s_* invalid: %w", LegacyControlPlaneInternalPrefix, err)
	}
	return NormalizeControlPlaneInternalURL(built)
}

// NormalizeControlPlaneInternalURL validates the URL and removes a trailing slash.
func NormalizeControlPlaneInternalURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("internal url is empty")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid internal url: %w", err)
	}
	if parsed.Scheme == "" {
		return "", fmt.Errorf("internal url scheme is required")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("internal url host is required")
	}
	if parsed.User != nil {
		return "", fmt.Errorf("internal url must not include user info")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("internal url must not include query or fragment")
	}
	if path := strings.TrimSpace(parsed.Path); path != "" && path != "/" {
		return "", fmt.Errorf("internal url must not include path")
	}
	parsed.Path = ""
	parsed.RawPath = ""

	return strings.TrimRight(parsed.String(), "/"), nil
}

// ParseControlPlaneInternalURL parses a validated internal URL.
func ParseControlPlaneInternalURL(raw string) (*url.URL, error) {
	normalized, err := NormalizeControlPlaneInternalURL(raw)
	if err != nil {
		return nil, err
	}
	parsed, err := url.Parse(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid internal url: %w", err)
	}
	return parsed, nil
}

// ValidateControlPlaneInternalURLForAuthMode applies transport requirements.
func ValidateControlPlaneInternalURLForAuthMode(rawURL, authMode string, requireExplicitPort bool) error {
	parsed, err := ParseControlPlaneInternalURL(rawURL)
	if err != nil {
		return err
	}
	mode := strings.ToLower(strings.TrimSpace(authMode))
	if mode == "mtls" && !strings.EqualFold(parsed.Scheme, "https") {
		return fmt.Errorf("INTERNAL_AUTH_MODE=mtls requires %s scheme=https (actual: %s)", CanonicalControlPlaneInternalURLKey, parsed.Scheme)
	}
	if requireExplicitPort && parsed.Port() == "" {
		return fmt.Errorf("%s must include explicit port", CanonicalControlPlaneInternalURLKey)
	}
	return nil
}

// ControlPlaneInternalListenAddrFromURL returns the configured listener port.
func ControlPlaneInternalListenAddrFromURL(rawURL string) (string, error) {
	parsed, err := ParseControlPlaneInternalURL(rawURL)
	if err != nil {
		return "", err
	}
	port := strings.TrimSpace(parsed.Port())
	if port == "" {
		return "", fmt.Errorf("%s must include explicit port", CanonicalControlPlaneInternalURLKey)
	}
	return ":" + port, nil
}
