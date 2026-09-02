import { ipcMain, type BrowserWindow } from "electron";

import type {
  ActionResult,
  DeploymentSelectionResult,
  DesktopSettings,
  NativeClientSnapshot,
  ThemeMode,
} from "../shared/contracts.js";

export interface DesktopRuntime {
  snapshot(): Promise<NativeClientSnapshot>;
  selectCloudDeployment(): Promise<DeploymentSelectionResult>;
  validateSelfHostedDeployment(
    baseUrl: string,
  ): Promise<DeploymentSelectionResult>;
  connect(network?: string): Promise<ActionResult>;
  disconnect(): Promise<ActionResult>;
  selectNetwork(network: string): Promise<ActionResult>;
}

export interface SettingsPort {
  get(): DesktopSettings;
  setTheme(theme: ThemeMode): DesktopSettings;
  setStartAtLogin(enabled: boolean): DesktopSettings;
}

export type IpcDependencies = {
  runtime: DesktopRuntime;
  settings: SettingsPort;
  selectedDeployment: () => DeploymentSelectionResult;
  window: () => BrowserWindow | null;
};

const invokeChannels = [
  "nanami:status",
  "nanami:deployment_get",
  "nanami:deployment_select_cloud",
  "nanami:deployment_validate_self_hosted",
  "nanami:connect",
  "nanami:disconnect",
  "nanami:network_select",
  "nanami:settings_get",
  "nanami:settings_set_theme",
  "nanami:settings_set_start_at_login",
] as const;

export function registerDesktopIpc(dependencies: IpcDependencies): () => void {
  const { runtime, settings, selectedDeployment, window } = dependencies;

  ipcMain.handle("nanami:status", () => runtime.snapshot());
  ipcMain.handle("nanami:deployment_get", () => selectedDeployment());
  ipcMain.handle("nanami:deployment_select_cloud", () =>
    runtime.selectCloudDeployment(),
  );
  ipcMain.handle(
    "nanami:deployment_validate_self_hosted",
    (_event, baseUrl: unknown) =>
      runtime.validateSelfHostedDeployment(
        requireBoundedString(baseUrl, "baseUrl", 2048),
      ),
  );
  ipcMain.handle("nanami:connect", (_event, network: unknown) =>
    runtime.connect(optionalBoundedString(network, "network", 160)),
  );
  ipcMain.handle("nanami:disconnect", () => runtime.disconnect());
  ipcMain.handle("nanami:network_select", (_event, network: unknown) =>
    runtime.selectNetwork(requireBoundedString(network, "network", 160)),
  );
  ipcMain.handle("nanami:settings_get", () => settings.get());
  ipcMain.handle("nanami:settings_set_theme", (_event, value: unknown) => {
    const next = settings.setTheme(requireTheme(value));
    window()?.webContents.send("nanami:settings", next);
    return next;
  });
  ipcMain.handle(
    "nanami:settings_set_start_at_login",
    (_event, enabled: unknown) => {
      if (typeof enabled !== "boolean") {
        throw new TypeError("enabled must be a boolean");
      }
      const next = settings.setStartAtLogin(enabled);
      window()?.webContents.send("nanami:settings", next);
      return next;
    },
  );

  return () => {
    for (const channel of invokeChannels) {
      ipcMain.removeHandler(channel);
    }
  };
}

function requireTheme(value: unknown): ThemeMode {
  if (value === "light" || value === "dark" || value === "system") {
    return value;
  }
  throw new TypeError("theme must be light, dark, or system");
}

function requireBoundedString(
  value: unknown,
  name: string,
  maxLength: number,
): string {
  if (typeof value !== "string") {
    throw new TypeError(`${name} must be a string`);
  }
  const normalized = value.trim();
  if (
    !normalized ||
    normalized.length > maxLength ||
    normalized.includes("\0")
  ) {
    throw new TypeError(`${name} is invalid`);
  }
  return normalized;
}

function optionalBoundedString(
  value: unknown,
  name: string,
  maxLength: number,
): string | undefined {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  return requireBoundedString(value, name, maxLength);
}
