package gatewayruntime

const (
	WorkerKindGatewayWorker       = "gateway_worker"
	CanonicalUnitWorkload         = "gateway_workload"
	SharingModelSharedRouteDomain = "shared_route_domain"

	ResponsibilityTransport = "transport"
	ResponsibilityRouting   = "routing"
	ResponsibilityDNS       = "dns"

	WireGuardInterfaceRole = "gateway_transport"
)

// ProtocolReadiness reports availability for one transport protocol.
type ProtocolReadiness struct {
	Protocol          string `json:"protocol,omitempty"`
	Supported         bool   `json:"supported"`
	Verified          bool   `json:"verified"`
	Ready             bool   `json:"ready"`
	SchedulerEligible bool   `json:"schedulerEligible"`
	Reason            string `json:"reason,omitempty"`
	LastVerifiedAt    *int64 `json:"lastVerifiedAt,omitempty"`
	Stale             bool   `json:"stale,omitempty"`
}

// WorkerModel describes one gateway worker and its advertised capabilities.
type WorkerModel struct {
	WorkerKind                       string                       `json:"workerKind"`
	CanonicalUnit                    string                       `json:"canonicalUnit"`
	DefaultSharingModel              string                       `json:"defaultSharingModel"`
	DefaultIsolationMode             string                       `json:"defaultIsolationMode,omitempty"`
	SupportsAdvancedIsolation        bool                         `json:"supportsAdvancedIsolation"`
	RouteDomainIsolationRequired     bool                         `json:"routeDomainIsolationRequired"`
	RouteDomainIsolationAvailable    bool                         `json:"routeDomainIsolationAvailable"`
	RouteDomainIsolationVerified     bool                         `json:"routeDomainIsolationVerified"`
	RouteDomainOverlapSafe           bool                         `json:"routeDomainOverlapSafe"`
	StaleRouteDomains                int                          `json:"staleRouteDomains,omitempty"`
	SourceDestinationACLRequired     bool                         `json:"sourceDestinationAclRequired"`
	SourceDestinationACLSupported    bool                         `json:"sourceDestinationAclSupported,omitempty"`
	SourceDestinationACLEnforced     bool                         `json:"sourceDestinationAclEnforced"`
	SourceDestinationACLVerified     bool                         `json:"sourceDestinationAclVerified,omitempty"`
	SourceDestinationACLDefaultDeny  bool                         `json:"sourceDestinationAclDefaultDeny,omitempty"`
	SourceDestinationACLPolicyBacked bool                         `json:"sourceDestinationAclPolicyBacked,omitempty"`
	SourceDestinationACLBaseline     bool                         `json:"sourceDestinationAclBaseline,omitempty"`
	StaleACLRules                    int                          `json:"staleAclRules,omitempty"`
	MaxRouteDomains                  int                          `json:"maxRouteDomains,omitempty"`
	ActiveRouteDomains               int                          `json:"activeRouteDomains,omitempty"`
	MaxPeers                         int                          `json:"maxPeers,omitempty"`
	ActivePeers                      int                          `json:"activePeers,omitempty"`
	MaxProtocolInterfaces            int                          `json:"maxProtocolInterfaces,omitempty"`
	ActiveProtocolInterfaces         int                          `json:"activeProtocolInterfaces,omitempty"`
	LoadScore                        *float64                     `json:"loadScore,omitempty"`
	CapacityState                    string                       `json:"capacityState,omitempty"`
	Responsibilities                 []string                     `json:"responsibilities,omitempty"`
	LegacySingleRuntimeCompatibility bool                         `json:"legacySingleRuntimeCompatibility"`
	Protocols                        map[string]ProtocolReadiness `json:"protocols,omitempty"`
	GRE                              ProtocolReadiness            `json:"gre"`
}

// WireGuardInterfaceState reports observed runtime state for one managed
// WireGuard interface workload.
type WireGuardInterfaceState struct {
	Name            string  `json:"name"`
	Role            string  `json:"role,omitempty"`
	Address         string  `json:"address,omitempty"`
	DesiredAddress  string  `json:"desiredAddress,omitempty"`
	HealthTarget    string  `json:"healthTarget,omitempty"`
	ListenPort      *int    `json:"listenPort,omitempty"`
	ListenPortKnown bool    `json:"listenPortKnown"`
	Error           *string `json:"error,omitempty"`
}

// CanonicalWorkerModel returns the default worker capabilities.
func CanonicalWorkerModel() WorkerModel {
	return WorkerModel{
		WorkerKind:                       WorkerKindGatewayWorker,
		CanonicalUnit:                    CanonicalUnitWorkload,
		DefaultSharingModel:              SharingModelSharedRouteDomain,
		DefaultIsolationMode:             SharingModelSharedRouteDomain,
		SupportsAdvancedIsolation:        false,
		RouteDomainIsolationRequired:     true,
		RouteDomainIsolationAvailable:    false,
		SourceDestinationACLRequired:     true,
		SourceDestinationACLEnforced:     false,
		CapacityState:                    "unknown",
		Responsibilities:                 []string{ResponsibilityTransport, ResponsibilityRouting, ResponsibilityDNS},
		LegacySingleRuntimeCompatibility: true,
		Protocols:                        CanonicalProtocolReadiness(),
		GRE: ProtocolReadiness{
			Protocol: ProtocolGRE,
			Reason:   GatewayReasonGREUnsupported,
		},
	}
}

// NormalizeWorkerModelProtocols fills missing entries and keeps legacy GRE data in sync.
func NormalizeWorkerModelProtocols(model WorkerModel) WorkerModel {
	out := CloneWorkerModel(model)
	if len(out.Protocols) == 0 {
		out.Protocols = CanonicalProtocolReadiness()
	}
	for _, id := range ProtocolIDs() {
		readiness, ok := out.Protocols[id]
		if !ok {
			readiness = CanonicalProtocolReadiness()[id]
		}
		if readiness.Protocol == "" {
			readiness.Protocol = id
		}
		if id == ProtocolGRE && readiness.Reason == "" && !readiness.SchedulerEligible {
			readiness.Reason = GatewayReasonGREUnsupported
		}
		if id == ProtocolWireGuard && readiness.Reason == "" && !readiness.SchedulerEligible {
			readiness.Reason = GatewayReasonWireGuardUnverified
		}
		out.Protocols[id] = readiness
	}
	if out.GRE.Protocol == "" {
		out.GRE.Protocol = ProtocolGRE
	}
	if out.GRE == (ProtocolReadiness{Protocol: ProtocolGRE}) || (out.GRE.Reason == "" && !out.GRE.Supported && !out.GRE.Verified && !out.GRE.Ready && !out.GRE.SchedulerEligible) {
		out.GRE = out.Protocols[ProtocolGRE]
	}
	if out.GRE.Reason == "" && !out.GRE.SchedulerEligible {
		out.GRE.Reason = GatewayReasonGREUnsupported
	}
	out.Protocols[ProtocolGRE] = out.GRE
	return out
}

// CloneWorkerModel returns a deep copy.
func CloneWorkerModel(model WorkerModel) WorkerModel {
	cloned := model
	cloned.Responsibilities = append([]string(nil), model.Responsibilities...)
	cloned.Protocols = CloneProtocolReadinessMap(model.Protocols)
	return cloned
}
