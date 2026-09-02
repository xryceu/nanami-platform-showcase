package gatewayorchestration

const (
	TopologyManagerMandatory = "manager_mandatory"

	DesiredStateDeliveryControlPlaneToManager = "control_plane_to_manager"
	NodeConfigDeliveryHeartbeatFanout         = "manager_heartbeat_fanout"
	DaemonHeartbeatPathManager                = "daemon_to_manager"
	DaemonObservedStatePathControlPlane       = "daemon_to_control_plane"
	ManagerObservedStatePathControlPlane      = "manager_to_control_plane"
	DirectControlPlaneFanoutFutureOnly        = "future_only"

	ResponsibilityDesiredStatePolling  = "desired_state_polling"
	ResponsibilityHeartbeatFanout      = "heartbeat_fanout"
	ResponsibilityNodeScopedOverrides  = "node_scoped_overrides"
	ResponsibilityObservedAggregation  = "observed_state_aggregation"
	ResponsibilityFleetCoordinationAPI = "fleet_coordination_boundary"
)

// Model describes how desired and observed state move between services.
type Model struct {
	CommunityTopology           string   `json:"communityTopology"`
	SaaSTopology                string   `json:"saasTopology"`
	ManagerMandatoryInCommunity bool     `json:"managerMandatoryInCommunity"`
	ManagerMandatoryInSaaS      bool     `json:"managerMandatoryInSaaS"`
	DesiredStateDelivery        string   `json:"desiredStateDelivery"`
	NodeConfigDelivery          string   `json:"nodeConfigDelivery"`
	DaemonHeartbeatPath         string   `json:"daemonHeartbeatPath"`
	DaemonObservedStatePath     string   `json:"daemonObservedStatePath"`
	ManagerObservedStatePath    string   `json:"managerObservedStatePath"`
	DirectControlPlaneFanout    string   `json:"directControlPlaneFanout"`
	ManagerResponsibilities     []string `json:"managerResponsibilities,omitempty"`
}

// CanonicalModel returns the supported manager-in-path topology.
func CanonicalModel() Model {
	return Model{
		CommunityTopology:           TopologyManagerMandatory,
		SaaSTopology:                TopologyManagerMandatory,
		ManagerMandatoryInCommunity: true,
		ManagerMandatoryInSaaS:      true,
		DesiredStateDelivery:        DesiredStateDeliveryControlPlaneToManager,
		NodeConfigDelivery:          NodeConfigDeliveryHeartbeatFanout,
		DaemonHeartbeatPath:         DaemonHeartbeatPathManager,
		DaemonObservedStatePath:     DaemonObservedStatePathControlPlane,
		ManagerObservedStatePath:    ManagerObservedStatePathControlPlane,
		DirectControlPlaneFanout:    DirectControlPlaneFanoutFutureOnly,
		ManagerResponsibilities: []string{
			ResponsibilityDesiredStatePolling,
			ResponsibilityHeartbeatFanout,
			ResponsibilityNodeScopedOverrides,
			ResponsibilityObservedAggregation,
			ResponsibilityFleetCoordinationAPI,
		},
	}
}

// CloneModel returns a deep copy.
func CloneModel(model Model) Model {
	cloned := model
	cloned.ManagerResponsibilities = append([]string(nil), model.ManagerResponsibilities...)
	return cloned
}
