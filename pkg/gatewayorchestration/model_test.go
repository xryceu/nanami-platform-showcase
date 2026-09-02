package gatewayorchestration

import "testing"

func TestCanonicalModelKeepsGatewayManagerMandatory(t *testing.T) {
	model := CanonicalModel()

	if model.CommunityTopology != TopologyManagerMandatory {
		t.Fatalf("community topology must require gateway-manager, got %q", model.CommunityTopology)
	}
	if model.SaaSTopology != TopologyManagerMandatory {
		t.Fatalf("saas topology must require gateway-manager, got %q", model.SaaSTopology)
	}
	if !model.ManagerMandatoryInCommunity {
		t.Fatalf("gateway-manager must be mandatory in Community")
	}
	if !model.ManagerMandatoryInSaaS {
		t.Fatalf("gateway-manager must be mandatory in SaaS")
	}
}

func TestCanonicalModelKeepsManagerInPathRuntimeFlow(t *testing.T) {
	model := CanonicalModel()

	if model.DesiredStateDelivery != DesiredStateDeliveryControlPlaneToManager {
		t.Fatalf("desired state must flow control-plane -> gateway-manager, got %q", model.DesiredStateDelivery)
	}
	if model.NodeConfigDelivery != NodeConfigDeliveryHeartbeatFanout {
		t.Fatalf("node config must be delivered through manager heartbeat fan-out, got %q", model.NodeConfigDelivery)
	}
	if model.DaemonHeartbeatPath != DaemonHeartbeatPathManager {
		t.Fatalf("daemon heartbeat must target gateway-manager, got %q", model.DaemonHeartbeatPath)
	}
	if model.DaemonObservedStatePath != DaemonObservedStatePathControlPlane {
		t.Fatalf("daemon observed state must return to control-plane, got %q", model.DaemonObservedStatePath)
	}
	if model.ManagerObservedStatePath != ManagerObservedStatePathControlPlane {
		t.Fatalf("manager observed state must return to control-plane, got %q", model.ManagerObservedStatePath)
	}
	if model.DirectControlPlaneFanout != DirectControlPlaneFanoutFutureOnly {
		t.Fatalf("direct control-plane fan-out must stay future-only, got %q", model.DirectControlPlaneFanout)
	}
}

func TestCanonicalModelStatesManagerResponsibilities(t *testing.T) {
	model := CanonicalModel()
	required := []string{
		ResponsibilityDesiredStatePolling,
		ResponsibilityHeartbeatFanout,
		ResponsibilityNodeScopedOverrides,
		ResponsibilityObservedAggregation,
		ResponsibilityFleetCoordinationAPI,
	}

	for _, responsibility := range required {
		if !containsString(model.ManagerResponsibilities, responsibility) {
			t.Fatalf("gateway-manager responsibility %q missing from %#v", responsibility, model.ManagerResponsibilities)
		}
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
