import { describe, expect, it } from "vitest";

import {
  isApprovalWaitingSnapshot,
  isAuthenticatedSnapshot,
  sanitizeDiagnosticCopy,
} from "../src/renderer/snapshot-view.js";
import type { NativeClientSnapshot } from "../src/shared/contracts.js";

describe("renderer state boundary", () => {
  it("does not treat a signed-in approval-pending session as authenticated", () => {
    const snapshot = baseSnapshot();
    snapshot.approval = { required: true, state: "authorization_pending" };
    expect(isApprovalWaitingSnapshot(snapshot)).toBe(true);
    expect(isAuthenticatedSnapshot(snapshot)).toBe(false);
  });

  it("redacts credentials, local paths, and stack frames from product diagnostics", () => {
    const source =
      "ipc_token=private at Runtime.connect (/Users/person/app/daemon.ts:20) /var/run/nanami/daemon.sock";
    const sanitized = sanitizeDiagnosticCopy(source);
    expect(sanitized).not.toContain("private");
    expect(sanitized).not.toContain("/Users/person");
    expect(sanitized).not.toContain("daemon.sock");
    expect(sanitized).toContain("credential=[redacted]");
  });
});

function baseSnapshot(): NativeClientSnapshot {
  return {
    state: { state: "approved_not_connected" },
    session: { signedIn: true, sessionState: "authenticated" },
    daemon: { available: true, authenticatedIpc: true },
    connection: { desiredConnected: false, observedConnected: false },
    runtime: { state: "available", issues: [] },
  };
}
