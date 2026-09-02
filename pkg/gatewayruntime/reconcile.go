package gatewayruntime

import (
	"sort"
	"strings"
	"time"
)

const (
	ReconcileSafetyReady       = "ready"
	ReconcileSafetyReconciling = "reconciling"
	ReconcileSafetyDegraded    = "degraded"
	ReconcileSafetyUnsafe      = "unsafe"

	ReconcileStatusReconciling    = "reconciling"
	ReconcileStatusApplied        = "applied"
	ReconcileStatusFailed         = "failed"
	ReconcileStatusStale          = "stale"
	ReconcileStatusCleanupPending = "cleanup_pending"

	ReconcilePhaseInterface   = "interface"
	ReconcilePhasePeer        = "peer"
	ReconcilePhaseRouteDomain = "route_domain"
	ReconcilePhaseACL         = "acl"
	ReconcilePhaseRouting     = "routing"
	ReconcilePhaseDNS         = "dns"
	ReconcilePhaseHealth      = "health"

	ReconcilePhaseStatusApplied = "applied"
	ReconcilePhaseStatusFailed  = "failed"

	ReconcileReasonApplyFailed                = "apply_failed"
	ReconcileReasonCleanupFailed              = "cleanup_failed"
	ReconcileReasonInterfaceApplyFailed       = "interface_apply_failed"
	ReconcileReasonWireGuardPeerApplyFailed   = "wireguard_peer_apply_failed"
	ReconcileReasonStaleWireGuardPeer         = "stale_wireguard_peer"
	ReconcileReasonStaleRouteDomain           = "stale_route_domain"
	ReconcileReasonStaleRoute                 = "stale_route"
	ReconcileReasonStaleACLRule               = "stale_acl_rule"
	ReconcileReasonACLUnverified              = "acl_unverified"
	ReconcileReasonACLDefaultDenyMissing      = "acl_default_deny_missing"
	ReconcileReasonStaleDNSState              = "stale_dns_state"
	ReconcileReasonRouteDomainApplyFailed     = "route_domain_apply_failed"
	ReconcileReasonSourceDestinationACLFailed = "source_destination_acl_failed"
	ReconcileReasonRoutingApplyFailed         = "routing_apply_failed"
	ReconcileReasonDNSRuntimeApplyFailed      = "dns_runtime_apply_failed"
	ReconcileReasonAppliedVersionDrift        = "applied_version_drift"
	ReconcileReasonRuntimeVerificationFailed  = "runtime_verification_failed"
	ReconcileResourceWireGuardPeer            = "wireguard_peer"
	ReconcileResourceInterface                = "interface"
	ReconcileResourceRouteDomain              = "route_domain"
	ReconcileResourceRoute                    = "route"
	ReconcileResourceACLRule                  = "acl_rule"
	ReconcileResourceDNSState                 = "dns_state"

	RuntimeAccessStateActive              = "active"
	RuntimeAccessStateRevoked             = "revoked"
	RuntimeAccessStatePendingRemoval      = "pending_removal"
	RuntimeAccessStateRemoved             = "removed"
	RuntimeAccessStateStaleRuntimePresent = "stale_runtime_present"
	RuntimeAccessStateUnknown             = "unknown"

	RuntimePeerStateExpected     = "expected"
	RuntimePeerStateMissing      = "missing"
	RuntimePeerStateStale        = "stale"
	RuntimePeerStateUnauthorized = "unauthorized"
	RuntimePeerStateUnknown      = "unknown"

	RuntimeCleanupStateNotRequired    = "not_required"
	RuntimeCleanupStatePending        = "pending"
	RuntimeCleanupStateCompleted      = "completed"
	RuntimeCleanupStateFailed         = "failed"
	RuntimeCleanupStateManualRequired = "manual_required"
	RuntimeCleanupStateUnknown        = "unknown"
)

// RuntimeReconcileStatus tracks one desired-state apply attempt.
type RuntimeReconcileStatus struct {
	DesiredRevision int64                   `json:"desiredRevision,omitempty"`
	AppliedRevision int64                   `json:"appliedRevision,omitempty"`
	DesiredStateID  string                  `json:"desiredStateId,omitempty"`
	Status          string                  `json:"status,omitempty"`
	SafetyState     string                  `json:"safetyState,omitempty"`
	DataplaneReady  bool                    `json:"dataplaneReady"`
	Reasons         []string                `json:"reasons,omitempty"`
	Phases          []RuntimeReconcilePhase `json:"phases,omitempty"`
	StaleResources  []RuntimeStaleResource  `json:"staleResources,omitempty"`
	Cleanup         *RuntimeCleanupStatus   `json:"cleanup,omitempty"`
	LastError       *RuntimeReconcileError  `json:"lastError,omitempty"`
	UpdatedAt       int64                   `json:"updatedAt,omitempty"`
}

type RuntimeReconcilePhase struct {
	Name          string `json:"name"`
	Scope         string `json:"scope,omitempty"`
	Status        string `json:"status,omitempty"`
	ReasonCode    string `json:"reasonCode,omitempty"`
	ReasonMessage string `json:"reasonMessage,omitempty"`
}

type RuntimeStaleResource struct {
	Type            string `json:"type"`
	NetworkID       string `json:"networkId,omitempty"`
	DeviceID        string `json:"deviceId,omitempty"`
	Interface       string `json:"interface,omitempty"`
	PublicKeyPrefix string `json:"publicKeyPrefix,omitempty"`
	ReasonCode      string `json:"reasonCode,omitempty"`
}

type RuntimeReconcileError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// RuntimeCleanupStatus summarizes stale resources and required cleanup.
type RuntimeCleanupStatus struct {
	AccessState                string                 `json:"accessState,omitempty"`
	RuntimePeerState           string                 `json:"runtimePeerState,omitempty"`
	CleanupState               string                 `json:"cleanupState,omitempty"`
	DesiredPeerCount           int                    `json:"desiredPeerCount"`
	AppliedPeerCount           int                    `json:"appliedPeerCount"`
	UndesiredPeerCount         int                    `json:"undesiredPeerCount"`
	StaleResourceCount         int                    `json:"staleResourceCount"`
	LastObservedAt             int64                  `json:"lastObservedAt,omitempty"`
	LastCleanupAt              int64                  `json:"lastCleanupAt,omitempty"`
	LastCleanupError           *RuntimeReconcileError `json:"lastCleanupError,omitempty"`
	ManualRemediationAvailable bool                   `json:"manualRemediationAvailable"`
	RecommendedAction          string                 `json:"recommendedAction,omitempty"`
	Reasons                    []string               `json:"reasons,omitempty"`
	StaleResources             []RuntimeStaleResource `json:"staleResources,omitempty"`
}

func NewRuntimeReconcileStatus(desiredRevision int64, desiredStateID string) RuntimeReconcileStatus {
	return RuntimeReconcileStatus{
		DesiredRevision: desiredRevision,
		DesiredStateID:  strings.TrimSpace(desiredStateID),
		Status:          ReconcileStatusReconciling,
		SafetyState:     ReconcileSafetyReconciling,
		DataplaneReady:  false,
		UpdatedAt:       time.Now().UTC().Unix(),
	}
}

func (s *RuntimeReconcileStatus) MarkApplied(appliedRevision int64) {
	if s == nil {
		return
	}
	if appliedRevision > 0 {
		s.AppliedRevision = appliedRevision
	}
	s.Status = ReconcileStatusApplied
	s.SafetyState = ReconcileSafetyReady
	s.DataplaneReady = true
	s.LastError = nil
	s.UpdatedAt = time.Now().UTC().Unix()
}

func (s *RuntimeReconcileStatus) AttachCleanupStatus(desiredPeerCount, appliedPeerCount int) {
	if s == nil {
		return
	}
	cleanup := RuntimeCleanupStatusForReconcile(*s, desiredPeerCount, appliedPeerCount)
	s.Cleanup = &cleanup
}

func (s *RuntimeReconcileStatus) AddPhase(name, scope, status, reasonCode, message string) {
	if s == nil {
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	s.Phases = append(s.Phases, RuntimeReconcilePhase{
		Name:          name,
		Scope:         strings.TrimSpace(scope),
		Status:        strings.TrimSpace(status),
		ReasonCode:    strings.TrimSpace(reasonCode),
		ReasonMessage: strings.TrimSpace(message),
	})
	s.UpdatedAt = time.Now().UTC().Unix()
}

func (s *RuntimeReconcileStatus) MarkUnsafe(reasonCode, message string) {
	if s == nil {
		return
	}
	reasonCode = strings.TrimSpace(reasonCode)
	message = strings.TrimSpace(message)
	if reasonCode == "" {
		reasonCode = ReconcileReasonApplyFailed
	}
	s.Status = reconcileStatusForReason(reasonCode)
	s.SafetyState = ReconcileSafetyUnsafe
	s.DataplaneReady = false
	s.addReason(reasonCode)
	if message != "" {
		s.LastError = &RuntimeReconcileError{Code: reasonCode, Message: message}
	}
	s.UpdatedAt = time.Now().UTC().Unix()
}

func (s *RuntimeReconcileStatus) AddStaleResource(resource RuntimeStaleResource) {
	if s == nil {
		return
	}
	resource.Type = strings.TrimSpace(resource.Type)
	if resource.Type == "" {
		return
	}
	resource.NetworkID = strings.TrimSpace(resource.NetworkID)
	resource.DeviceID = strings.TrimSpace(resource.DeviceID)
	resource.Interface = strings.TrimSpace(resource.Interface)
	resource.PublicKeyPrefix = strings.TrimSpace(resource.PublicKeyPrefix)
	resource.ReasonCode = strings.TrimSpace(resource.ReasonCode)
	if resource.ReasonCode == "" {
		resource.ReasonCode = reasonForStaleResource(resource.Type)
	}
	s.StaleResources = append(s.StaleResources, resource)
	s.MarkUnsafe(resource.ReasonCode, "stale Nanami-owned runtime resource remains")
}

func (s *RuntimeReconcileStatus) addReason(reason string) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return
	}
	for _, existing := range s.Reasons {
		if existing == reason {
			return
		}
	}
	s.Reasons = append(s.Reasons, reason)
	sort.Strings(s.Reasons)
}

func reconcileStatusForReason(reason string) string {
	switch reason {
	case ReconcileReasonStaleWireGuardPeer,
		ReconcileReasonStaleRouteDomain,
		ReconcileReasonStaleRoute,
		ReconcileReasonStaleACLRule,
		ReconcileReasonStaleDNSState,
		ReconcileReasonCleanupFailed:
		return ReconcileStatusCleanupPending
	case ReconcileReasonAppliedVersionDrift:
		return ReconcileStatusStale
	default:
		return ReconcileStatusFailed
	}
}

func reasonForStaleResource(resourceType string) string {
	switch resourceType {
	case ReconcileResourceWireGuardPeer:
		return ReconcileReasonStaleWireGuardPeer
	case ReconcileResourceRouteDomain:
		return ReconcileReasonStaleRouteDomain
	case ReconcileResourceRoute:
		return ReconcileReasonStaleRoute
	case ReconcileResourceACLRule:
		return ReconcileReasonStaleACLRule
	case ReconcileResourceDNSState:
		return ReconcileReasonStaleDNSState
	default:
		return ReconcileReasonCleanupFailed
	}
}

func RuntimeCleanupStatusForReconcile(reconcile RuntimeReconcileStatus, desiredPeerCount, appliedPeerCount int) RuntimeCleanupStatus {
	if desiredPeerCount < 0 {
		desiredPeerCount = 0
	}
	if appliedPeerCount < 0 {
		appliedPeerCount = 0
	}
	status := RuntimeCleanupStatus{
		AccessState:        RuntimeAccessStateActive,
		RuntimePeerState:   RuntimePeerStateExpected,
		CleanupState:       RuntimeCleanupStateNotRequired,
		DesiredPeerCount:   desiredPeerCount,
		AppliedPeerCount:   appliedPeerCount,
		LastObservedAt:     reconcile.UpdatedAt,
		Reasons:            uniqueSortedStrings(reconcile.Reasons),
		StaleResources:     append([]RuntimeStaleResource(nil), reconcile.StaleResources...),
		StaleResourceCount: len(reconcile.StaleResources),
	}
	status.UndesiredPeerCount = staleResourceCount(reconcile.StaleResources, ReconcileResourceWireGuardPeer)
	if status.UndesiredPeerCount == 0 && appliedPeerCount > desiredPeerCount {
		status.UndesiredPeerCount = appliedPeerCount - desiredPeerCount
	}
	if reconcile.LastError != nil {
		status.LastCleanupError = &RuntimeReconcileError{
			Code:    strings.TrimSpace(reconcile.LastError.Code),
			Message: strings.TrimSpace(reconcile.LastError.Message),
		}
	}
	if reconcile.Status == "" && reconcile.SafetyState == "" {
		status.AccessState = RuntimeAccessStateUnknown
		status.RuntimePeerState = RuntimePeerStateUnknown
		status.CleanupState = RuntimeCleanupStateUnknown
		status.RecommendedAction = "wait_for_gateway_runtime_observation"
		return NormalizeRuntimeCleanupStatus(status)
	}
	if status.StaleResourceCount > 0 || status.UndesiredPeerCount > 0 {
		status.AccessState = RuntimeAccessStateStaleRuntimePresent
		status.RuntimePeerState = RuntimePeerStateStale
		status.CleanupState = RuntimeCleanupStateManualRequired
		status.ManualRemediationAvailable = true
		status.RecommendedAction = "run_gateway_daemon_purge_stale_dry_run_before_apply"
		return NormalizeRuntimeCleanupStatus(status)
	}
	switch strings.TrimSpace(reconcile.Status) {
	case ReconcileStatusCleanupPending:
		status.AccessState = RuntimeAccessStatePendingRemoval
		status.RuntimePeerState = RuntimePeerStateUnknown
		status.CleanupState = RuntimeCleanupStatePending
		status.RecommendedAction = "wait_for_reconcile_or_inspect_gateway_daemon"
	case ReconcileStatusFailed:
		status.RuntimePeerState = RuntimePeerStateUnknown
		status.CleanupState = RuntimeCleanupStateFailed
		status.ManualRemediationAvailable = true
		status.RecommendedAction = "inspect_gateway_runtime_apply_error"
	case ReconcileStatusStale:
		status.AccessState = RuntimeAccessStateUnknown
		status.RuntimePeerState = RuntimePeerStateUnknown
		status.CleanupState = RuntimeCleanupStatePending
		status.RecommendedAction = "wait_for_current_desired_state_apply"
	case ReconcileStatusApplied:
		status.CleanupState = RuntimeCleanupStateNotRequired
		status.LastCleanupAt = reconcile.UpdatedAt
		if desiredPeerCount == 0 && appliedPeerCount == 0 {
			status.AccessState = RuntimeAccessStateRemoved
		}
	default:
		if appliedPeerCount < desiredPeerCount {
			status.RuntimePeerState = RuntimePeerStateMissing
			status.CleanupState = RuntimeCleanupStatePending
			status.RecommendedAction = "wait_for_gateway_peer_apply"
		}
	}
	if appliedPeerCount > desiredPeerCount && status.RuntimePeerState == RuntimePeerStateExpected {
		status.AccessState = RuntimeAccessStateStaleRuntimePresent
		status.RuntimePeerState = RuntimePeerStateUnauthorized
		status.CleanupState = RuntimeCleanupStateManualRequired
		status.ManualRemediationAvailable = true
		status.RecommendedAction = "inspect_gateway_peer_inventory"
	}
	return NormalizeRuntimeCleanupStatus(status)
}

func NormalizeRuntimeCleanupStatus(status RuntimeCleanupStatus) RuntimeCleanupStatus {
	status.AccessState = strings.TrimSpace(status.AccessState)
	status.RuntimePeerState = strings.TrimSpace(status.RuntimePeerState)
	status.CleanupState = strings.TrimSpace(status.CleanupState)
	if status.DesiredPeerCount < 0 {
		status.DesiredPeerCount = 0
	}
	if status.AppliedPeerCount < 0 {
		status.AppliedPeerCount = 0
	}
	if status.UndesiredPeerCount < 0 {
		status.UndesiredPeerCount = 0
	}
	if status.StaleResourceCount < 0 {
		status.StaleResourceCount = 0
	}
	if status.AccessState == "" {
		if status.UndesiredPeerCount > 0 || status.StaleResourceCount > 0 {
			status.AccessState = RuntimeAccessStateStaleRuntimePresent
		} else {
			status.AccessState = RuntimeAccessStateUnknown
		}
	}
	if status.RuntimePeerState == "" {
		if status.UndesiredPeerCount > 0 || status.StaleResourceCount > 0 {
			status.RuntimePeerState = RuntimePeerStateStale
		} else {
			status.RuntimePeerState = RuntimePeerStateUnknown
		}
	}
	if status.CleanupState == "" {
		if status.UndesiredPeerCount > 0 || status.StaleResourceCount > 0 {
			status.CleanupState = RuntimeCleanupStateManualRequired
		} else {
			status.CleanupState = RuntimeCleanupStateUnknown
		}
	}
	status.RecommendedAction = strings.TrimSpace(status.RecommendedAction)
	if status.RecommendedAction == "" && status.ManualRemediationAvailable {
		status.RecommendedAction = "inspect_gateway_runtime_cleanup"
	}
	status.Reasons = uniqueSortedStrings(status.Reasons)
	status.StaleResources = append([]RuntimeStaleResource(nil), status.StaleResources...)
	return status
}

func staleResourceCount(resources []RuntimeStaleResource, resourceType string) int {
	count := 0
	for _, resource := range resources {
		if strings.TrimSpace(resource.Type) == resourceType {
			count++
		}
	}
	return count
}

func uniqueSortedStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
