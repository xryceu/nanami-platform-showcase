package urlutil

import (
	"fmt"
	"net/url"
	"strings"
)

// Build composes a URL from protocol, host, and port.
// If port is empty or a default for the protocol, it is omitted.
func Build(protocol, host, port string) (string, error) {
	scheme := strings.ToLower(strings.TrimSpace(protocol))
	if scheme == "" {
		return "", fmt.Errorf("protocol is required")
	}
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("protocol must be http or https, got %q", protocol)
	}
	targetHost := strings.TrimSpace(host)
	if targetHost == "" {
		return "", fmt.Errorf("host is required")
	}
	p := strings.TrimSpace(port)
	if (scheme == "https" && p == "443") || (scheme == "http" && p == "80") {
		p = ""
	}
	if p != "" {
		targetHost = targetHost + ":" + p
	}
	u := url.URL{
		Scheme: scheme,
		Host:   targetHost,
	}
	return u.String(), nil
}
