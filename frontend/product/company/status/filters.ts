import type {
  StatusIncident,
  StatusRange,
  StatusSeverity,
  StatusSeverityFilter,
} from "./types";

const rangeToDurationMs: Record<StatusRange, number> = {
  "15m": 15 * 60 * 1000,
  "1h": 60 * 60 * 1000,
  "24h": 24 * 60 * 60 * 1000,
  "7d": 7 * 24 * 60 * 60 * 1000,
};

const severityWeight: Record<StatusSeverity, number> = {
  critical: 0,
  warning: 1,
  info: 2,
};

export function buildIncidentId(
  source: string,
  code: string,
  scopeId: string,
): string {
  return `${source}:${code}:${scopeId}`;
}

export function parseStatusRange(value?: string | null): StatusRange {
  if (value === "15m" || value === "1h" || value === "24h" || value === "7d") {
    return value;
  }
  return "1h";
}

export function parseStatusSeverityFilter(
  value?: string | null,
): StatusSeverityFilter {
  if (value === "critical" || value === "warning" || value === "info") {
    return value;
  }
  return "all";
}

export function getStatusWindow(
  range: StatusRange,
  nowMs = Date.now(),
): {
  fromIso: string;
  toIso: string;
} {
  const durationMs = rangeToDurationMs[range];
  return {
    fromIso: new Date(nowMs - durationMs).toISOString(),
    toIso: new Date(nowMs).toISOString(),
  };
}

function parseIsoTimestamp(value?: string): number {
  if (!value) {
    return 0;
  }
  const parsed = Date.parse(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

export function sortStatusIncidents(
  incidents: StatusIncident[],
): StatusIncident[] {
  return [...incidents].sort((a, b) => {
    const severityDelta =
      severityWeight[a.severity] - severityWeight[b.severity];
    if (severityDelta !== 0) {
      return severityDelta;
    }
    const bTs = parseIsoTimestamp(b.lastSeenAt);
    const aTs = parseIsoTimestamp(a.lastSeenAt);
    if (bTs !== aTs) {
      return bTs - aTs;
    }
    return a.id.localeCompare(b.id);
  });
}

export function applySeverityFilter(
  incidents: StatusIncident[],
  severity: StatusSeverityFilter,
): StatusIncident[] {
  if (severity === "all") {
    return incidents;
  }
  return incidents.filter((incident) => incident.severity === severity);
}

export function buildStatusSourceFailedIncident(
  sourceName: string,
  reason: string,
): StatusIncident {
  return {
    id: buildIncidentId("control-plane", "status_source_failed", sourceName),
    severity: "warning",
    source: "control-plane",
    code: "status_source_failed",
    title: `Data source failed: ${sourceName}`,
    message: `Failed to load ${sourceName}: ${reason}`,
    lastSeenAt: new Date().toISOString(),
    actions: [
      {
        label: "Retry",
        onClick: "retry",
      },
    ],
  };
}

export function rangeLabel(range: StatusRange): string {
  switch (range) {
    case "15m":
      return "15m";
    case "1h":
      return "1h";
    case "24h":
      return "24h";
    case "7d":
      return "7d";
    default:
      return "1h";
  }
}

function criticalThreshold(range: StatusRange): number {
  switch (range) {
    case "15m":
      return 5;
    case "1h":
      return 10;
    case "24h":
      return 20;
    case "7d":
      return 40;
    default:
      return 10;
  }
}

function warningThreshold(range: StatusRange): number {
  switch (range) {
    case "15m":
      return 2;
    case "1h":
      return 4;
    case "24h":
      return 8;
    case "7d":
      return 15;
    default:
      return 4;
  }
}

export function severityFromCount(
  count: number,
  range: StatusRange,
): StatusSeverity {
  if (count >= criticalThreshold(range)) {
    return "critical";
  }
  if (count >= warningThreshold(range)) {
    return "warning";
  }
  return "info";
}
