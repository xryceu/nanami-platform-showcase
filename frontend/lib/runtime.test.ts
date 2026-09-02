import { describe, expect, it } from "vitest";

import {
  filterRuntimeNodes,
  formatLastSeen,
  getRuntimeSummary,
  normalizeConnectionStatus,
  type RuntimeNode,
} from "./runtime";

const nodes: RuntimeNode[] = [
  {
    id: "ready",
    name: "Ready gateway",
    owner: "Operations",
    platform: "Linux",
    rawStatus: "READY",
    lastSeenAt: 1_000,
    endpoint: "198.51.100.10:51820",
    warning: null,
  },
  {
    id: "pending",
    name: "Pending gateway",
    owner: "Operations",
    platform: "Linux",
    rawStatus: "provisioning",
    lastSeenAt: 900,
    endpoint: null,
    warning: "Waiting for runtime confirmation.",
  },
  {
    id: "unknown",
    name: "Unknown device",
    owner: "Engineering",
    platform: "macOS",
    rawStatus: "unexpected-value",
    lastSeenAt: 800,
    endpoint: null,
    warning: "No status reported.",
  },
];

describe("runtime presentation", () => {
  it("maps backend-compatible success values to connected", () => {
    expect(normalizeConnectionStatus(" online ")).toBe("connected");
    expect(normalizeConnectionStatus("READY")).toBe("connected");
  });

  it("fails unknown values closed as offline", () => {
    expect(normalizeConnectionStatus("unexpected-value")).toBe("offline");
    expect(normalizeConnectionStatus(null)).toBe("offline");
  });

  it("filters by runtime posture without changing source order", () => {
    expect(filterRuntimeNodes(nodes, "ready").map((node) => node.id)).toEqual([
      "ready",
    ]);
    expect(
      filterRuntimeNodes(nodes, "attention").map((node) => node.id),
    ).toEqual(["pending", "unknown"]);
  });

  it("summarizes ready and attention states", () => {
    expect(getRuntimeSummary(nodes)).toEqual({ ready: 1, attention: 2 });
  });

  it("formats bounded relative time and rejects invalid timestamps", () => {
    expect(formatLastSeen(970, 1_000)).toBe("Just now");
    expect(formatLastSeen(700, 1_000)).toBe("5m ago");
    expect(formatLastSeen(0, 1_000)).toBe("Unknown");
    expect(formatLastSeen(1_001, 1_000)).toBe("Unknown");
  });
});
