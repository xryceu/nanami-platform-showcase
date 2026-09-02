package gatewayruntime

import "testing"

func TestRuntimeCleanupStatusForStalePeer(t *testing.T) {
	reconcile := NewRuntimeReconcileStatus(7, "state-7")
	reconcile.AddStaleResource(RuntimeStaleResource{
		Type:            ReconcileResourceWireGuardPeer,
		Interface:       "wg-gateway",
		PublicKeyPrefix: "abcdef12",
	})

	cleanup := RuntimeCleanupStatusForReconcile(reconcile, 2, 3)
	if cleanup.AccessState != RuntimeAccessStateStaleRuntimePresent {
		t.Fatalf("expected stale runtime access state, got %q", cleanup.AccessState)
	}
	if cleanup.RuntimePeerState != RuntimePeerStateStale {
		t.Fatalf("expected stale peer state, got %q", cleanup.RuntimePeerState)
	}
	if cleanup.CleanupState != RuntimeCleanupStateManualRequired {
		t.Fatalf("expected manual cleanup state, got %q", cleanup.CleanupState)
	}
	if cleanup.DesiredPeerCount != 2 || cleanup.AppliedPeerCount != 3 || cleanup.UndesiredPeerCount != 1 {
		t.Fatalf("unexpected peer counts: %#v", cleanup)
	}
	if !cleanup.ManualRemediationAvailable {
		t.Fatal("expected manual remediation to be available")
	}
	if cleanup.RecommendedAction != "run_gateway_daemon_purge_stale_dry_run_before_apply" {
		t.Fatalf("unexpected recommended action: %q", cleanup.RecommendedAction)
	}
}

func TestRuntimeCleanupStatusForAppliedState(t *testing.T) {
	reconcile := NewRuntimeReconcileStatus(7, "state-7")
	reconcile.MarkApplied(7)

	cleanup := RuntimeCleanupStatusForReconcile(reconcile, 2, 2)
	if cleanup.CleanupState != RuntimeCleanupStateNotRequired {
		t.Fatalf("expected no cleanup required, got %q", cleanup.CleanupState)
	}
	if cleanup.RuntimePeerState != RuntimePeerStateExpected {
		t.Fatalf("expected expected peer state, got %q", cleanup.RuntimePeerState)
	}
	if cleanup.UndesiredPeerCount != 0 || cleanup.StaleResourceCount != 0 {
		t.Fatalf("expected clean counts, got %#v", cleanup)
	}
}
