package internaltransport

import (
	"path/filepath"
	"strings"
)

const (
	// InternalMTLSDirKey configures the directory containing service certificates.
	InternalMTLSDirKey = "INTERNAL_MTLS_DIR"
	// DefaultInternalMTLSDir is the default container mount path.
	DefaultInternalMTLSDir = "/etc/nanami/internal-certs"
)

// MTLSBundlePaths contains the certificate files used by internal services.
type MTLSBundlePaths struct {
	Dir                    string
	CAFile                 string
	ControlPlaneCertFile   string
	ControlPlaneKeyFile    string
	GatewayManagerCertFile string
	GatewayManagerKeyFile  string
	GatewayDaemonCertFile  string
	GatewayDaemonKeyFile   string
}

// ResolveMTLSBundlePaths expands a bundle directory into service-specific paths.
func ResolveMTLSBundlePaths(dir string) MTLSBundlePaths {
	resolvedDir := strings.TrimSpace(dir)
	if resolvedDir == "" {
		resolvedDir = DefaultInternalMTLSDir
	}
	resolvedDir = filepath.Clean(resolvedDir)

	return MTLSBundlePaths{
		Dir:                    resolvedDir,
		CAFile:                 filepath.Join(resolvedDir, "ca.pem"),
		ControlPlaneCertFile:   filepath.Join(resolvedDir, "control-plane.pem"),
		ControlPlaneKeyFile:    filepath.Join(resolvedDir, "control-plane.key"),
		GatewayManagerCertFile: filepath.Join(resolvedDir, "gateway-manager.pem"),
		GatewayManagerKeyFile:  filepath.Join(resolvedDir, "gateway-manager.key"),
		GatewayDaemonCertFile:  filepath.Join(resolvedDir, "gateway-daemon.pem"),
		GatewayDaemonKeyFile:   filepath.Join(resolvedDir, "gateway-daemon.key"),
	}
}
