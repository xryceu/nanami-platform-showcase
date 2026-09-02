export type StatusSeverity = "critical" | "warning" | "info";
export type StatusSeverityFilter = "all" | StatusSeverity;
export type StatusSource =
  | "gateway"
  | "audit"
  | "config"
  | "control-plane"
  | "dashboard"
  | "gateway-manager"
  | "unknown";
export type StatusAudience = "platform" | "system";
export type StatusRange = "15m" | "1h" | "24h" | "7d";
export type StatusIncidentAction = {
  label: string;
  href?: string;
  onClick?: "retry" | "openLogs" | "openDetails";
};
export type StatusIncident = {
  id: string;
  severity: StatusSeverity;
  source: StatusSource;
  code: string;
  title: string;
  message: string;
  scope?: {
    tenantId?: string;
    tenantSlug?: string;
    nodeId?: string;
    gatewayId?: string;
    regionId?: string;
  };
  firstSeenAt?: string;
  lastSeenAt?: string;
  count?: number;
  details?: Record<string, unknown>;
  playbook?: string[];
  actions?: StatusIncidentAction[];
};
export type StatusResponse = {
  incidents: StatusIncident[];
  generatedAt: string;
  range: StatusRange;
};
