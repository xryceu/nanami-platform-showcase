import type { NativeClientSnapshot } from "../../../shared/contracts.js";

export type WindowsTrayPresentation = {
  statusLabel: string;
  canConnect: boolean;
  canDisconnect: boolean;
};

export function windowsTrayPresentation(
  snapshot: NativeClientSnapshot | null,
): WindowsTrayPresentation | null {
  if (snapshot?.platform?.kind !== "windows") {
    return null;
  }

  const connected = Boolean(snapshot.connection.observedConnected);
  const disconnectable =
    connected &&
    ["connected", "degraded", "disconnecting"].includes(snapshot.state.state);
  return {
    statusLabel: windowsStatusLabel(snapshot),
    canConnect:
      snapshot.session.signedIn &&
      snapshot.platform.connectionAllowed &&
      ["approved_not_connected", "disconnected"].includes(snapshot.state.state),
    canDisconnect: snapshot.session.signedIn && disconnectable,
  };
}

export function windowsStatusLabel(snapshot: NativeClientSnapshot): string {
  const platform = snapshot.platform;
  if (!platform || platform.kind !== "windows") {
    return "Runtime unavailable";
  }
  if (platform.service.state !== "running") {
    return platform.service.state === "starting"
      ? "Service starting"
      : "Service unavailable";
  }
  if (platform.driver.state !== "ready") {
    return "Driver unavailable";
  }
  if (platform.permission.state === "elevation_required") {
    return "Administrator approval required";
  }
  if (snapshot.connection.observedConnected) {
    return snapshot.state.state === "degraded" ? "Degraded" : "Connected";
  }
  return snapshot.state.state === "connecting" ? "Connecting" : "Disconnected";
}
