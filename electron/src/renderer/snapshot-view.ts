import type {
  ClientState,
  IssueRow,
  NativeClientSnapshot,
} from "../shared/contracts.js";

export function selectedNetwork(snapshot: NativeClientSnapshot | null): string {
  const selected = snapshot?.networks?.find((network) => network.selected);
  return (
    selected?.id ||
    selected?.name ||
    snapshot?.connection.networkId ||
    snapshot?.connection.networkName ||
    ""
  );
}

export function isAuthenticatedSnapshot(
  snapshot: NativeClientSnapshot | null,
): boolean {
  if (!snapshot?.session.signedIn) {
    return false;
  }
  return !isSignedOutSnapshot(snapshot) && !isApprovalWaitingSnapshot(snapshot);
}

export function isSignedOutSnapshot(
  snapshot: NativeClientSnapshot | null,
): boolean {
  return (
    snapshot?.state.state === "signed_out" ||
    snapshot?.session.sessionState === "loggedout"
  );
}

export function isApprovalWaitingSnapshot(
  snapshot: NativeClientSnapshot | null,
): boolean {
  const approvalState = String(
    snapshot?.approval?.state ?? snapshot?.session.sessionState ?? "",
  ).toLowerCase();
  return (
    snapshot?.state.state === "needs_device_approval" ||
    approvalState === "authorization_pending" ||
    approvalState === "slow_down" ||
    approvalState === "authenticating"
  );
}

export function stateLabel(state: ClientState): string {
  switch (state) {
    case "signed_out":
      return "Signed out";
    case "needs_device_approval":
      return "Approval";
    case "connected":
      return "Connected";
    case "connecting":
      return "Connecting";
    case "degraded":
      return "Degraded";
    case "runtime_unavailable":
      return "Runtime unavailable";
    default:
      return "Disconnected";
  }
}

export function primaryStateMessage(
  snapshot: NativeClientSnapshot | null,
): string {
  if (
    snapshot?.platform?.kind === "windows" &&
    snapshot.platform.reasonCode !== "windows_ready"
  ) {
    return snapshot.platform.userMessage;
  }
  const state = snapshot?.state.state;
  if (state === "runtime_unavailable") {
    return "Nanami runtime is not running. Start Nanami to connect this device.";
  }
  if (state === "needs_device_approval" || state === "signed_out") {
    return "Sign in to connect this device.";
  }
  return (
    primaryActionMessage(snapshot?.state.userMessage) ||
    "Local runtime status is loading."
  );
}

export function primaryRuntimeIssueMessage(issue: IssueRow): string {
  const mapped = primaryReasonMessages[issue.reasonCode];
  if (mapped) {
    return mapped;
  }
  if (hasTechnicalRuntimeDetail(issue.message)) {
    return primaryReasonMessages.runtime_unavailable;
  }
  return "Nanami could not complete this connection. Open Diagnostics for a safe reason code.";
}

export function sanitizeDiagnosticCopy(value: string): string {
  return sanitizePrimaryCopy(value)
    .replace(
      /\b(?:access|refresh|session|ipc)[_-]?token\b\s*[:=]\s*\S+/gi,
      "credential=[redacted]",
    )
    .replace(/(?:[A-Za-z]:\\|\/)(?:[^\s]+[\\/])+[^\s]+/g, "[local path]")
    .replace(/\bat\s+[\w$.<>]+\s*\([^)]*\)/g, "")
    .replace(/\s+/g, " ")
    .trim();
}

const primaryReasonMessages: Record<string, string> = {
  signed_out: "Sign in to connect this device.",
  device_approval_required:
    "Complete sign-in approval before connecting this device.",
  runtime_unavailable:
    "Nanami is not running on this device. Start Nanami and try again.",
  runtime_degraded:
    "The local Nanami service needs attention before this connection can continue.",
  runtime_truth_unknown:
    "Nanami is waiting for this device to report its connection state.",
  protocol_runtime_unverified:
    "WireGuard readiness has not been confirmed on this device.",
  protocol_not_supported_by_endpoint:
    "This destination does not support the selected connection method.",
  gateway_quota_exceeded: "This workspace has reached its gateway limit.",
  shared_gateway_capacity_unavailable:
    "No gateway capacity is available right now. Try again later.",
};

function primaryActionMessage(value: string | undefined): string | null {
  if (!value) {
    return null;
  }
  if (hasTechnicalRuntimeDetail(value)) {
    return "Nanami runtime is not running. Start Nanami to connect this device.";
  }
  const sanitized = sanitizePrimaryCopy(value);
  return sanitized || null;
}

function hasTechnicalRuntimeDetail(value: string): boolean {
  return /daemon\.sock|ECONNREFUSED|connect\s+[^ ]+\.sock|socket path/i.test(
    value,
  );
}

function sanitizePrimaryCopy(value: string): string {
  return value
    .replace(/\/[^\s"]*daemon\.sock/gi, "local runtime")
    .replace(/connect\s+ECONNREFUSED\s+\S+/gi, "Nanami runtime is not running")
    .replace(/\s+/g, " ")
    .trim();
}
