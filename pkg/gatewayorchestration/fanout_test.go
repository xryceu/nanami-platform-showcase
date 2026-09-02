package gatewayorchestration

import "testing"

func TestBuildFanoutContractIsStableAndNodeScoped(t *testing.T) {
	t.Parallel()

	contract := BuildFanoutContract("manager-1", "gateway-1", "WireGuard", "rev-1", 42)
	again := BuildFanoutContract("manager-1", "gateway-1", "wireguard", "rev-1", 42)
	if contract.DesiredStateID == "" {
		t.Fatalf("expected desired state id")
	}
	if contract.DesiredStateID != again.DesiredStateID {
		t.Fatalf("expected stable desired state id, got %q vs %q", contract.DesiredStateID, again.DesiredStateID)
	}
	if contract.Delivery != NodeConfigDeliveryHeartbeatFanout {
		t.Fatalf("expected heartbeat fan-out delivery, got %q", contract.Delivery)
	}
	if contract.DirectControlPlane != DirectControlPlaneFanoutFutureOnly {
		t.Fatalf("expected direct fanout to remain future-only, got %q", contract.DirectControlPlane)
	}

	nodeScoped := contract.ForDaemon("daemon-1", "gateway_worker")
	if nodeScoped.TargetDaemonID != "daemon-1" {
		t.Fatalf("expected daemon target, got %q", nodeScoped.TargetDaemonID)
	}
	if nodeScoped.TargetWorkerID != "gateway_worker" {
		t.Fatalf("expected worker target, got %q", nodeScoped.TargetWorkerID)
	}
	if nodeScoped.IdempotencyKey == contract.IdempotencyKey {
		t.Fatalf("expected daemon-scoped idempotency key")
	}
}
