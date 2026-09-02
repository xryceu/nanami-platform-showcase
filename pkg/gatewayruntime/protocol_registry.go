package gatewayruntime

const (
	ProtocolWireGuard = "wireguard"
	ProtocolGRE       = "gre"

	ProtocolProductionStatusAvailable           = "available"
	ProtocolProductionStatusInternalUnavailable = "internal_unavailable"

	ProtocolGatewayCostClassNone        = "none"
	ProtocolGatewayCostClassUnsupported = "unsupported"

	EndpointKindUserDevice           = "user_device"
	EndpointKindServer               = "server"
	EndpointKindSite                 = "site"
	EndpointKindRelay                = "relay"
	EndpointKindGateway              = "gateway"
	EndpointKindMobileDevice         = "mobile_device"
	EndpointKindDesktopCLI           = "desktop_cli"
	EndpointKindTenantOwnedGateway   = "tenant_owned_gateway"
	EndpointKindNanamiSharedGateway  = "nanami_shared_gateway"
	GatewayClassTenantOwned          = "tenant_owned"
	GatewayClassNanamiShared         = "nanami_shared"
	GatewayClassDedicated            = "dedicated"
	GatewayReasonWireGuardUnverified = "gateway_wireguard_runtime_unverified"
	GatewayReasonGREUnsupported      = "gateway_gre_not_supported"
)

// ProtocolRegistryEntry describes static properties of a transport protocol.
type ProtocolRegistryEntry struct {
	ID                           string   `json:"id"`
	DisplayName                  string   `json:"displayName"`
	ProductionStatus             string   `json:"productionStatus"`
	Encrypted                    bool     `json:"encrypted"`
	ClientDeviceCapable          bool     `json:"clientDeviceCapable"`
	SiteToSiteCapable            bool     `json:"siteToSiteCapable"`
	RequiresPublicEndpoint       bool     `json:"requiresPublicEndpoint"`
	NATFriendly                  bool     `json:"natFriendly"`
	RequiresKernelSupport        bool     `json:"requiresKernelSupport"`
	RequiresPrivilegedRuntime    bool     `json:"requiresPrivilegedRuntime"`
	SupportsRoaming              bool     `json:"supportsRoaming"`
	SupportsRouteDomainIsolation bool     `json:"supportsRouteDomainIsolation"`
	SupportsACLVerification      bool     `json:"supportsAclVerification"`
	SupportsHealthProbe          bool     `json:"supportsHealthProbe"`
	SupportsMobile               bool     `json:"supportsMobile"`
	SupportsDesktop              bool     `json:"supportsDesktop"`
	RuntimeReadinessRequired     bool     `json:"runtimeReadinessRequired"`
	DefaultForEndpointKinds      []string `json:"defaultForEndpointKinds"`
	AllowedForGatewayClasses     []string `json:"allowedForGatewayClasses"`
	GatewayCostClass             string   `json:"gatewayCostClass"`
}

// ProtocolRegistry returns a fresh protocol catalog.
func ProtocolRegistry() map[string]ProtocolRegistryEntry {
	return map[string]ProtocolRegistryEntry{
		ProtocolWireGuard: {
			ID:                           ProtocolWireGuard,
			DisplayName:                  "WireGuard",
			ProductionStatus:             ProtocolProductionStatusAvailable,
			Encrypted:                    true,
			ClientDeviceCapable:          true,
			SiteToSiteCapable:            true,
			RequiresPublicEndpoint:       false,
			NATFriendly:                  true,
			RequiresKernelSupport:        true,
			RequiresPrivilegedRuntime:    true,
			SupportsRoaming:              true,
			SupportsRouteDomainIsolation: true,
			SupportsACLVerification:      true,
			SupportsHealthProbe:          true,
			SupportsMobile:               true,
			SupportsDesktop:              true,
			RuntimeReadinessRequired:     true,
			DefaultForEndpointKinds:      []string{EndpointKindUserDevice, EndpointKindMobileDevice, EndpointKindDesktopCLI, EndpointKindServer, EndpointKindGateway, EndpointKindTenantOwnedGateway, EndpointKindNanamiSharedGateway},
			AllowedForGatewayClasses:     []string{GatewayClassTenantOwned, GatewayClassNanamiShared, GatewayClassDedicated},
			GatewayCostClass:             ProtocolGatewayCostClassNone,
		},
		ProtocolGRE: {
			ID:                           ProtocolGRE,
			DisplayName:                  "GRE",
			ProductionStatus:             ProtocolProductionStatusInternalUnavailable,
			Encrypted:                    false,
			ClientDeviceCapable:          false,
			SiteToSiteCapable:            true,
			RequiresPublicEndpoint:       true,
			NATFriendly:                  false,
			RequiresKernelSupport:        true,
			RequiresPrivilegedRuntime:    true,
			SupportsRoaming:              false,
			SupportsRouteDomainIsolation: true,
			SupportsACLVerification:      true,
			SupportsHealthProbe:          false,
			SupportsMobile:               false,
			SupportsDesktop:              false,
			RuntimeReadinessRequired:     true,
			DefaultForEndpointKinds:      []string{EndpointKindSite, EndpointKindGateway, EndpointKindTenantOwnedGateway},
			AllowedForGatewayClasses:     []string{},
			GatewayCostClass:             ProtocolGatewayCostClassUnsupported,
		},
	}
}

// ProtocolEntry returns a registry entry by id.
func ProtocolEntry(id string) (ProtocolRegistryEntry, bool) {
	entry, ok := ProtocolRegistry()[id]
	return entry, ok
}

// ProtocolIDs returns protocol IDs in display order.
func ProtocolIDs() []string {
	return []string{ProtocolWireGuard, ProtocolGRE}
}

// CanonicalProtocolReadiness returns conservative default readiness.
func CanonicalProtocolReadiness() map[string]ProtocolReadiness {
	return map[string]ProtocolReadiness{
		ProtocolWireGuard: {
			Protocol:  ProtocolWireGuard,
			Supported: true,
			Reason:    GatewayReasonWireGuardUnverified,
		},
		ProtocolGRE: {
			Protocol: ProtocolGRE,
			Reason:   GatewayReasonGREUnsupported,
		},
	}
}

// CloneProtocolReadinessMap returns a deep copy.
func CloneProtocolReadinessMap(values map[string]ProtocolReadiness) map[string]ProtocolReadiness {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]ProtocolReadiness, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}
