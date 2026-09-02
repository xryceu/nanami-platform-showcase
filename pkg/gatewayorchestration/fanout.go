package gatewayorchestration

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	FanoutConflictBehaviorRejectStale = "reject_stale_observations"
	FanoutIdempotencyScopeDaemonApply = "daemon_apply"
)

// FanoutContract identifies a desired-state delivery from the control plane to
// one gateway daemon.
type FanoutContract struct {
	DesiredStateID   string `json:"desiredStateId"`
	IdempotencyKey   string `json:"idempotencyKey"`
	IdempotencyScope string `json:"idempotencyScope"`

	Delivery           string `json:"delivery"`
	ConflictBehavior   string `json:"conflictBehavior"`
	ManagerAgentID     string `json:"managerAgentId,omitempty"`
	GatewayDeviceID    string `json:"gatewayDeviceId,omitempty"`
	TargetDaemonID     string `json:"targetDaemonId,omitempty"`
	TargetWorkerID     string `json:"targetWorkerId,omitempty"`
	ProtocolType       string `json:"protocolType"`
	Revision           string `json:"revision"`
	DesiredRevision    int64  `json:"desiredRevision"`
	DirectControlPlane string `json:"directControlPlaneFanout"`
	ObservedStatePath  string `json:"observedStatePath"`
}

// BuildDesiredStateID creates a stable identifier without exposing config data.
func BuildDesiredStateID(gatewayDeviceID, protocolType, revision string, desiredRevision int64) string {
	gatewayDeviceID = strings.TrimSpace(gatewayDeviceID)
	protocolType = strings.ToLower(strings.TrimSpace(protocolType))
	revision = strings.TrimSpace(revision)
	if revision == "" && desiredRevision > 0 {
		revision = fmt.Sprintf("%d", desiredRevision)
	}
	raw := strings.Join([]string{"gateway-desired-state", gatewayDeviceID, protocolType, revision}, "|")
	sum := sha256.Sum256([]byte(raw))
	return "gds_" + hex.EncodeToString(sum[:10])
}

// BuildFanoutContract creates the manager-to-daemon delivery envelope.
func BuildFanoutContract(managerAgentID, gatewayDeviceID, protocolType, revision string, desiredRevision int64) FanoutContract {
	model := CanonicalModel()
	desiredStateID := BuildDesiredStateID(gatewayDeviceID, protocolType, revision, desiredRevision)
	return FanoutContract{
		DesiredStateID:     desiredStateID,
		IdempotencyKey:     desiredStateID,
		IdempotencyScope:   FanoutIdempotencyScopeDaemonApply,
		Delivery:           model.NodeConfigDelivery,
		ConflictBehavior:   FanoutConflictBehaviorRejectStale,
		ManagerAgentID:     strings.TrimSpace(managerAgentID),
		GatewayDeviceID:    strings.TrimSpace(gatewayDeviceID),
		ProtocolType:       strings.ToLower(strings.TrimSpace(protocolType)),
		Revision:           strings.TrimSpace(revision),
		DesiredRevision:    desiredRevision,
		DirectControlPlane: model.DirectControlPlaneFanout,
		ObservedStatePath:  model.DaemonObservedStatePath,
	}
}

// ForDaemon scopes the delivery to one daemon and worker.
func (c FanoutContract) ForDaemon(targetDaemonID, workerID string) FanoutContract {
	c.TargetDaemonID = strings.TrimSpace(targetDaemonID)
	c.TargetWorkerID = strings.TrimSpace(workerID)
	c.IdempotencyKey = strings.Join(
		[]string{
			strings.TrimSpace(c.DesiredStateID),
			strings.TrimSpace(c.TargetDaemonID),
			strings.TrimSpace(c.TargetWorkerID),
		},
		":",
	)
	return c
}
