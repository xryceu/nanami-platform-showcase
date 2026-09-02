package gatewayruntime

import "testing"

func TestCanonicalWorkerModelDefinesSharedGatewayWorker(t *testing.T) {
	model := CanonicalWorkerModel()

	if model.WorkerKind != WorkerKindGatewayWorker {
		t.Fatalf("worker kind must be %q, got %q", WorkerKindGatewayWorker, model.WorkerKind)
	}
	if model.CanonicalUnit != CanonicalUnitWorkload {
		t.Fatalf("canonical runtime unit must be %q, got %q", CanonicalUnitWorkload, model.CanonicalUnit)
	}
	if model.DefaultSharingModel != SharingModelSharedRouteDomain {
		t.Fatalf("default sharing model must be %q, got %q", SharingModelSharedRouteDomain, model.DefaultSharingModel)
	}
	if model.DefaultIsolationMode != SharingModelSharedRouteDomain {
		t.Fatalf("compat isolation field must report %q, got %q", SharingModelSharedRouteDomain, model.DefaultIsolationMode)
	}
	if model.SupportsAdvancedIsolation {
		t.Fatalf("worker model must not claim shared route-domain isolation before it exists")
	}
	if !model.RouteDomainIsolationRequired || model.RouteDomainIsolationAvailable {
		t.Fatalf("worker model must require route-domain isolation without claiming it is available")
	}
	if !model.SourceDestinationACLRequired || model.SourceDestinationACLEnforced {
		t.Fatalf("worker model must require source+destination ACL without claiming it is enforced")
	}
	if !model.LegacySingleRuntimeCompatibility {
		t.Fatalf("worker model must keep legacy single-runtime compatibility explicit")
	}
	if model.GRE.Supported || model.GRE.Verified || model.GRE.Ready || model.GRE.SchedulerEligible {
		t.Fatalf("GRE must default to fail-closed unsupported readiness, got %#v", model.GRE)
	}
	if model.GRE.Reason != "gateway_gre_not_supported" {
		t.Fatalf("GRE unsupported reason drifted: %#v", model.GRE)
	}
	if model.Protocols[ProtocolWireGuard].Protocol != ProtocolWireGuard || !model.Protocols[ProtocolWireGuard].Supported {
		t.Fatalf("WireGuard readiness must be present as supported but runtime-unverified, got %#v", model.Protocols)
	}
	if model.Protocols[ProtocolGRE].Reason != GatewayReasonGREUnsupported || model.Protocols[ProtocolGRE].SchedulerEligible {
		t.Fatalf("GRE readiness map must fail closed, got %#v", model.Protocols[ProtocolGRE])
	}
}

func TestCanonicalWorkerModelOwnsTransportRoutingAndDNSResponsibilities(t *testing.T) {
	model := CanonicalWorkerModel()
	required := []string{
		ResponsibilityTransport,
		ResponsibilityRouting,
		ResponsibilityDNS,
	}

	for _, responsibility := range required {
		if !containsString(model.Responsibilities, responsibility) {
			t.Fatalf("gateway worker responsibility %q missing from %#v", responsibility, model.Responsibilities)
		}
	}
}

func TestWireGuardInterfaceRoleIsRuntimeObservationOnly(t *testing.T) {
	if WireGuardInterfaceRole != "gateway_transport" {
		t.Fatalf("wireguard interface role drifted: %q", WireGuardInterfaceRole)
	}
}

func TestProtocolRegistryDefinesWireGuardAndGREFacts(t *testing.T) {
	wireguard, ok := ProtocolEntry(ProtocolWireGuard)
	if !ok {
		t.Fatalf("WireGuard registry entry missing")
	}
	if wireguard.ProductionStatus != ProtocolProductionStatusAvailable || !wireguard.Encrypted || !wireguard.ClientDeviceCapable || !wireguard.NATFriendly {
		t.Fatalf("WireGuard registry entry drifted: %#v", wireguard)
	}

	gre, ok := ProtocolEntry(ProtocolGRE)
	if !ok {
		t.Fatalf("GRE registry entry missing")
	}
	if gre.ProductionStatus != ProtocolProductionStatusInternalUnavailable || gre.Encrypted || gre.ClientDeviceCapable || !gre.RequiresPublicEndpoint {
		t.Fatalf("GRE registry entry must remain internal/unavailable and non-client-capable, got %#v", gre)
	}
}

func TestNormalizeWorkerModelProtocolsKeepsLegacyGREInSync(t *testing.T) {
	model := CanonicalWorkerModel()
	model.Protocols[ProtocolGRE] = ProtocolReadiness{
		Protocol:          ProtocolGRE,
		Supported:         true,
		Verified:          true,
		Ready:             true,
		SchedulerEligible: true,
	}
	model.GRE = ProtocolReadiness{}

	normalized := NormalizeWorkerModelProtocols(model)
	if !normalized.GRE.Supported || !normalized.Protocols[ProtocolGRE].Supported {
		t.Fatalf("expected GRE readiness to sync between legacy field and protocol map, got legacy=%#v map=%#v", normalized.GRE, normalized.Protocols[ProtocolGRE])
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
