export type Edition = "saas" | "community";

export type ControlCenterNode = {
  status?: string | null;
  connectionStatus?: string | null;
  connectivity?: {
    gatewayId?: string | null;
    gatewayNodeId?: string | null;
    connectionStatus?: string | null;
  } | null;
};
