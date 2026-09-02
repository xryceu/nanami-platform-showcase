import { sortStatusIncidents } from "./filters";
import type {
  StatusIncident,
  StatusRange,
  StatusResponse,
  StatusSeverity,
} from "./types";

const severityRank: Record<StatusSeverity, number> = {
  critical: 3,
  warning: 2,
  info: 1,
};

function incidentScopeKey(incident: StatusIncident): string {
  const scope = incident.scope;
  if (!scope) {
    return "global";
  }
  return [
    scope.gatewayId ?? "",
    scope.nodeId ?? "",
    scope.regionId ?? "",
    scope.tenantId ?? "",
    scope.tenantSlug ?? "",
  ].join(":");
}

function incidentMergeKey(incident: StatusIncident): string {
  const id = incident.id.trim();
  if (id) {
    return `id:${id}`;
  }
  return `shape:${incident.source}:${incident.code}:${incidentScopeKey(incident)}`;
}

function parseTimestamp(value?: string): number {
  if (!value) {
    return 0;
  }
  const parsed = Date.parse(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function shouldReplaceIncident(
  current: StatusIncident,
  next: StatusIncident,
): boolean {
  const currentRank = severityRank[current.severity] ?? 0;
  const nextRank = severityRank[next.severity] ?? 0;
  if (nextRank !== currentRank) {
    return nextRank > currentRank;
  }

  const currentLastSeen = parseTimestamp(current.lastSeenAt);
  const nextLastSeen = parseTimestamp(next.lastSeenAt);
  return nextLastSeen > currentLastSeen;
}

export function mergeStatusResponses(
  responses: StatusResponse[],
  range: StatusRange,
): StatusResponse {
  const merged = new Map<string, StatusIncident>();
  let generatedAt = 0;

  for (const response of responses) {
    generatedAt = Math.max(generatedAt, parseTimestamp(response.generatedAt));
    for (const incident of response.incidents) {
      const key = incidentMergeKey(incident);
      const current = merged.get(key);
      if (!current || shouldReplaceIncident(current, incident)) {
        merged.set(key, incident);
      }
    }
  }

  return {
    incidents: sortStatusIncidents(Array.from(merged.values())),
    generatedAt:
      generatedAt > 0
        ? new Date(generatedAt).toISOString()
        : new Date().toISOString(),
    range,
  };
}
