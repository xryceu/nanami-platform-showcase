export type ConnectionStatus = "connected" | "pending" | "offline" | "degraded";

export type RuntimeFilter = "all" | "ready" | "attention";

export type RuntimeNode = {
  id: string;
  name: string;
  owner: string;
  platform: string;
  rawStatus?: string | null;
  lastSeenAt: number;
  endpoint: string | null;
  warning: string | null;
};

export function normalizeConnectionStatus(
  rawStatus?: string | null,
): ConnectionStatus {
  const status = (rawStatus ?? "").toLowerCase().trim();

  if (["connected", "online", "ready", "active"].includes(status)) {
    return "connected";
  }
  if (["pending", "provisioning"].includes(status)) {
    return "pending";
  }
  if (status === "degraded") {
    return "degraded";
  }
  return "offline";
}

export function filterRuntimeNodes(
  nodes: RuntimeNode[],
  filter: RuntimeFilter,
): RuntimeNode[] {
  if (filter === "all") {
    return nodes;
  }

  return nodes.filter((node) => {
    const status = normalizeConnectionStatus(node.rawStatus);
    return filter === "ready" ? status === "connected" : status !== "connected";
  });
}

export function getRuntimeSummary(nodes: RuntimeNode[]): {
  ready: number;
  attention: number;
} {
  return nodes.reduce(
    (summary, node) => {
      if (normalizeConnectionStatus(node.rawStatus) === "connected") {
        summary.ready += 1;
      } else {
        summary.attention += 1;
      }
      return summary;
    },
    { ready: 0, attention: 0 },
  );
}

export function formatLastSeen(lastSeenAt: number, now: number): string {
  if (!Number.isFinite(lastSeenAt) || lastSeenAt <= 0 || lastSeenAt > now) {
    return "Unknown";
  }

  const elapsedSeconds = Math.floor(now - lastSeenAt);
  if (elapsedSeconds < 60) {
    return "Just now";
  }

  const elapsedMinutes = Math.floor(elapsedSeconds / 60);
  if (elapsedMinutes < 60) {
    return `${elapsedMinutes}m ago`;
  }

  const elapsedHours = Math.floor(elapsedMinutes / 60);
  return `${elapsedHours}h ago`;
}
