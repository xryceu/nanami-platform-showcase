export type ClientState =
  | "signed_out"
  | "needs_device_approval"
  | "approved_not_connected"
  | "connecting"
  | "connected"
  | "degraded"
  | "disconnecting"
  | "disconnected"
  | "error"
  | "runtime_unavailable";

export type ThemeMode = "light" | "dark" | "system";
export type DeploymentType = "cloud" | "self_hosted";

export type DeploymentProfile = {
  id: string;
  type: DeploymentType;
  displayName: string;
  baseUrl: string;
  normalizedOrigin: string;
  deploymentId: string;
  edition: "saas" | "community" | "enterprise_self_hosted" | "unknown";
  lastValidatedAt: string;
  capabilitiesVersion: string;
};

export type DeploymentSelectionResult = {
  ok: boolean;
  state: "setup_required" | "validating" | "selected" | "unavailable";
  profile?: DeploymentProfile;
  message?: string;
  reasonCode?: string;
};

export type IssueRow = {
  severity: "info" | "warning" | "critical" | string;
  reasonCode: string;
  message: string;
  recommendedAction?: string;
};

export type WindowsPlatformStatus = {
  kind: "windows";
  service: {
    state:
      "running" | "stopped" | "missing" | "starting" | "degraded" | "unknown";
    installed: boolean;
  };
  driver: {
    state: "ready" | "missing" | "incompatible" | "unknown";
    transport: "wireguard";
  };
  permission: {
    state: "standard_user" | "administrator" | "elevation_required" | "unknown";
  };
  connectionAllowed: boolean;
  reasonCode: string;
  userMessage: string;
};

export type NativeClientSnapshot = {
  state: {
    state: ClientState;
    reasonCode?: string;
    userMessage?: string;
    recommendedAction?: { action: string; label: string };
  };
  session: {
    signedIn: boolean;
    userEmail?: string;
    workspaceId?: string;
    workspaceName?: string;
    sessionState?: string;
  };
  approval?: {
    required: boolean;
    state?: string;
  };
  daemon: {
    available: boolean;
    authenticatedIpc: boolean;
    serviceState?: string;
    version?: string;
  };
  connection: {
    desiredConnected: boolean;
    observedConnected: boolean;
    phase?: string;
    networkId?: string;
    networkName?: string;
  };
  runtime: {
    state: string;
    reasonCode?: string;
    issues?: IssueRow[];
  };
  networks?: Array<{
    id: string;
    name?: string;
    selected: boolean;
  }>;
  platform?: WindowsPlatformStatus;
};

export type ActionResult = {
  ok: boolean;
  message?: string;
  unavailable?: boolean;
  reasonCode?: string;
  snapshot: NativeClientSnapshot;
};

export type DesktopSettings = {
  theme: ThemeMode;
  startAtLogin: boolean;
  startAtLoginAvailable?: boolean;
  startAtLoginMessage?: string;
  closeToTray: boolean;
};

export type DesktopBridge = {
  status(): Promise<NativeClientSnapshot>;
  deploymentGet(): Promise<DeploymentSelectionResult>;
  deploymentSelectCloud(): Promise<DeploymentSelectionResult>;
  deploymentValidateSelfHosted(
    baseUrl: string,
  ): Promise<DeploymentSelectionResult>;
  connect(network?: string): Promise<ActionResult>;
  disconnect(): Promise<ActionResult>;
  networkSelect(network: string): Promise<ActionResult>;
  settingsGet(): Promise<DesktopSettings>;
  settingsSetTheme(theme: ThemeMode): Promise<DesktopSettings>;
  settingsSetStartAtLogin(enabled: boolean): Promise<DesktopSettings>;
  onSnapshot(callback: (snapshot: NativeClientSnapshot) => void): () => void;
  onSettingsChanged(callback: (settings: DesktopSettings) => void): () => void;
};

declare global {
  interface Window {
    nanamiDesktop: DesktopBridge;
  }
}
