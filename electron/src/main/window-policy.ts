import path from "node:path";

import { BrowserWindow, nativeImage, nativeTheme, shell } from "electron";

export type WindowPaths = {
  rendererIndex: string;
  preload: string;
  icon: string;
};

export function createMainWindow(paths: WindowPaths): BrowserWindow {
  const appIcon = nativeImage.createFromPath(paths.icon);
  const window = new BrowserWindow({
    width: 1040,
    height: 720,
    minWidth: 680,
    minHeight: 560,
    title: "Nanami",
    backgroundColor: nativeTheme.shouldUseDarkColors ? "#0b0b0c" : "#fafafa",
    icon: appIcon.isEmpty() ? undefined : appIcon,
    show: false,
    webPreferences: {
      preload: path.resolve(paths.preload),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: true,
      webSecurity: true,
    },
  });

  window.removeMenu();
  void window.loadFile(paths.rendererIndex);
  window.once("ready-to-show", () => window.show());
  window.webContents.setWindowOpenHandler(({ url }) => {
    if (isSafeExternalUrl(url)) {
      void shell.openExternal(url);
    }
    return { action: "deny" };
  });
  window.webContents.on("will-navigate", (event) => event.preventDefault());
  return window;
}

export function isSafeExternalUrl(value: string): boolean {
  try {
    const parsed = new URL(value);
    return parsed.protocol === "https:" && !parsed.username && !parsed.password;
  } catch {
    return false;
  }
}
