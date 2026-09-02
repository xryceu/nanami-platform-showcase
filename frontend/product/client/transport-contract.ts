import type { ControlCenterNode as Node } from "./types";

export type ControlCenterCurrentTransportStatus =
  "online" | "mediated" | "connecting" | "offline";

export type ControlCenterCurrentTransportConnection =
  "direct" | "mediated" | "pending" | "offline";

export function hasGatewayMediatedTransport(
  node: Pick<Node, "connectivity">,
): boolean {
  return Boolean(
    node.connectivity?.gatewayNodeId || node.connectivity?.gatewayId,
  );
}

export function resolveCurrentTransportStatus(
  node: Node,
): ControlCenterCurrentTransportStatus {
  const status = (node.status ?? "").trim().toLowerCase();
  const connectionStatus = (
    node.connectivity?.connectionStatus ??
    node.connectionStatus ??
    status
  )
    .trim()
    .toLowerCase();

  if (
    status === "offline" ||
    connectionStatus === "offline" ||
    connectionStatus === "disconnected"
  ) {
    return "offline";
  }

  if (
    status === "pending" ||
    status === "provisioning" ||
    connectionStatus === "pending" ||
    connectionStatus === "connecting" ||
    connectionStatus === "degraded"
  ) {
    return "connecting";
  }

  if (hasGatewayMediatedTransport(node)) {
    return "mediated";
  }

  return "online";
}

export function resolveCurrentTransportConnection(
  node: Node,
  statusTone: ControlCenterCurrentTransportStatus,
): ControlCenterCurrentTransportConnection {
  if (statusTone === "offline") {
    return "offline";
  }
  if (statusTone === "connecting") {
    return "pending";
  }
  if (hasGatewayMediatedTransport(node)) {
    return "mediated";
  }
  return "direct";
}

export function normalizeLegacyControlCenterStatus(
  value: string | null | undefined,
): ControlCenterCurrentTransportStatus | "all" {
  const normalized = String(value ?? "")
    .trim()
    .toLowerCase();
  switch (normalized) {
    case "online":
    case "offline":
    case "connecting":
    case "mediated":
      return normalized;
    case "relay":
      return "mediated";
    default:
      return "all";
  }
}

export function normalizeLegacyControlCenterConnection(
  value: string | null | undefined,
): "all" | "direct" | "mediated" | "logical" | "active" {
  const normalized = String(value ?? "")
    .trim()
    .toLowerCase();
  switch (normalized) {
    case "direct":
    case "mediated":
      return normalized;
    case "logical":
    case "active":
      return normalized;
    case "relay":
      return "mediated";
    default:
      return "all";
  }
}

export function emptyPeerConnectionSummary(
  connectionType: ControlCenterCurrentTransportConnection,
): string {
  return connectionType === "mediated"
    ? "Only connections through gateways are currently visible"
    : "No active peer connections are currently visible";
}

export function defaultTransitPathLabel(
  connectionType: ControlCenterCurrentTransportConnection,
): string {
  return connectionType === "mediated" ? "Through gateway" : "Direct path";
}

export function defaultTransitPathDetail(
  connectionType: ControlCenterCurrentTransportConnection,
): string {
  return connectionType === "mediated"
    ? "This connection currently goes through a gateway."
    : "No gateway-managed transit applies.";
}

export function transportReasoningDetail(
  connectionType: ControlCenterCurrentTransportConnection,
): string {
  if (connectionType === "mediated") {
    return "The current connection goes through a gateway, so check gateway health and the configured route.";
  }
  if (connectionType === "direct") {
    return "The peer currently shows direct connectivity, so route policy rather than gateway mediation is the likelier cause of missing reachability.";
  }
  return "Transport is still pending, so configured route policy may exist before live connectivity is fully visible.";
}

export function dashboardTransportDescription(): string {
  return "Current peer connections and the share using gateways.";
}

export function dashboardTransportDetail(summary: {
  directLinks: number;
  mediatedLinks: number;
}): string {
  return `${summary.directLinks} direct and ${summary.mediatedLinks} through a gateway`;
}

export function dashboardGatewayMapDescription(): string {
  return "Where gateway workers and traffic through gateways are currently active.";
}

export function dashboardMediatedPeersTitle(): string {
  return "Peers connected through a gateway";
}

export function dashboardMediatedPeersDetail(): string {
  return "Peers whose current connection still depends on a gateway.";
}
