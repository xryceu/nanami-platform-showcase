import { describe, expect, it } from "vitest";

import { windowsTrayPresentation } from "../src/main/platforms/windows/tray.js";
import type { NativeClientSnapshot } from "../src/shared/contracts.js";

describe("Windows tray presentation", () => {
  it("exposes Connect only when sign-in and platform readiness are both observed", () => {
    const snapshot = windowsSnapshot("disconnected", false);
    expect(windowsTrayPresentation(snapshot)).toEqual({
      statusLabel: "Disconnected",
      canConnect: true,
      canDisconnect: false,
    });

    snapshot.platform!.connectionAllowed = false;
    expect(windowsTrayPresentation(snapshot)?.canConnect).toBe(false);
  });

  it.each(["connected", "degraded"] as const)(
    "exposes Disconnect for an observed %s tunnel",
    (state) => {
      const snapshot = windowsSnapshot(state, true);
      expect(windowsTrayPresentation(snapshot)).toMatchObject({
        canConnect: false,
        canDisconnect: true,
      });
    },
  );
});

function windowsSnapshot(
  state: NativeClientSnapshot["state"]["state"],
  connected: boolean,
): NativeClientSnapshot {
  return {
    state: { state },
    session: { signedIn: true, sessionState: "authenticated" },
    daemon: { available: true, authenticatedIpc: true },
    connection: { desiredConnected: connected, observedConnected: connected },
    runtime: { state: connected ? "healthy" : "available", issues: [] },
    platform: {
      kind: "windows",
      service: { state: "running", installed: true },
      driver: { state: "ready", transport: "wireguard" },
      permission: { state: "standard_user" },
      connectionAllowed: true,
      reasonCode: "windows_ready",
      userMessage: "Windows runtime is ready.",
    },
  };
}
