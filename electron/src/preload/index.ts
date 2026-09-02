import { contextBridge, ipcRenderer } from "electron";

import type {
  DesktopBridge,
  DesktopSettings,
  NativeClientSnapshot,
  ThemeMode,
} from "../shared/contracts.js";

const bridge: DesktopBridge = {
  status: () => ipcRenderer.invoke("nanami:status"),
  deploymentGet: () => ipcRenderer.invoke("nanami:deployment_get"),
  deploymentSelectCloud: () =>
    ipcRenderer.invoke("nanami:deployment_select_cloud"),
  deploymentValidateSelfHosted: (baseUrl: string) =>
    ipcRenderer.invoke("nanami:deployment_validate_self_hosted", baseUrl),
  connect: (network?: string) => ipcRenderer.invoke("nanami:connect", network),
  disconnect: () => ipcRenderer.invoke("nanami:disconnect"),
  networkSelect: (network: string) =>
    ipcRenderer.invoke("nanami:network_select", network),
  settingsGet: () => ipcRenderer.invoke("nanami:settings_get"),
  settingsSetTheme: (theme: ThemeMode) =>
    ipcRenderer.invoke("nanami:settings_set_theme", theme),
  settingsSetStartAtLogin: (enabled: boolean) =>
    ipcRenderer.invoke("nanami:settings_set_start_at_login", enabled),
  onSnapshot: (callback: (snapshot: NativeClientSnapshot) => void) => {
    const listener = (
      _event: Electron.IpcRendererEvent,
      snapshot: NativeClientSnapshot,
    ) => callback(snapshot);
    ipcRenderer.on("nanami:snapshot", listener);
    return () => ipcRenderer.removeListener("nanami:snapshot", listener);
  },
  onSettingsChanged: (callback: (settings: DesktopSettings) => void) => {
    const listener = (
      _event: Electron.IpcRendererEvent,
      settings: DesktopSettings,
    ) => callback(settings);
    ipcRenderer.on("nanami:settings", listener);
    return () => ipcRenderer.removeListener("nanami:settings", listener);
  },
};

contextBridge.exposeInMainWorld("nanamiDesktop", bridge);
