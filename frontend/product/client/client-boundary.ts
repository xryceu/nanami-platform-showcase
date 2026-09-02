import type { Edition } from "./types";

type PathRule = {
  readonly root: string;
  readonly exact?: boolean;
};

const normalizeRoot = (root: string) =>
  root.length > 1 && root.endsWith("/") ? root.slice(0, -1) : root;

const matchesRule = (pathname: string, rule: PathRule): boolean => {
  const root = normalizeRoot(rule.root);
  if (rule.exact) {
    return pathname === root;
  }
  return pathname === root || pathname.startsWith(`${root}/`);
};

function normalizeLocalPath(value?: string | null): string {
  if (!value) {
    return "/";
  }

  const trimmed = value.trim();
  if (
    !trimmed ||
    !trimmed.startsWith("/") ||
    trimmed.startsWith("//") ||
    trimmed.includes("\\") ||
    /[\u0000-\u001f\u007f]/.test(trimmed)
  ) {
    return "/";
  }

  try {
    const parsed = new URL(trimmed, "https://nanami.local");
    if (parsed.origin !== "https://nanami.local") {
      return "/";
    }
    return `${parsed.pathname}${parsed.search}${parsed.hash}`;
  } catch {
    return "/";
  }
}

export function clientBoundaryPathname(value?: string | null): string {
  const local = normalizeLocalPath(value);
  try {
    const parsed = new URL(local, "https://nanami.local");
    const pathname = parsed.pathname || "/";
    return pathname.length > 1 && pathname.endsWith("/")
      ? pathname.slice(0, -1)
      : pathname;
  } catch {
    return "/";
  }
}

export const deletedLegacyClientPathRules: readonly PathRule[] = [
  { root: "/dashboard" },
  { root: "/routes" },
  { root: "/groups" },
  { root: "/users" },
  { root: "/audit" },
  { root: "/topology" },
  { root: "/gateways" },
  { root: "/agents" },
  { root: "/regions" },
  { root: "/tunnels" },
];

export const companyOwnedClientPathRules: readonly PathRule[] = [
  { root: "/company" },
  { root: "/platform" },
  { root: "/admin" },
  { root: "/tenants" },
  { root: "/client-users" },
  { root: "/problem-center" },
  { root: "/support/inbox" },
  { root: "/team/groups" },
];

export const saasOnlyClientPathRules: readonly PathRule[] = [
  { root: "/settings/plan" },
  { root: "/support" },
];

export const communityOnlyClientPathRules: readonly PathRule[] = [
  { root: "/settings/platform-security" },
  { root: "/settings/integrations" },
  { root: "/settings/runtime" },
];

export const debugClientPathRules: readonly PathRule[] = [{ root: "/debug" }];

export function isDeletedLegacyClientPath(value?: string | null): boolean {
  const pathname = clientBoundaryPathname(value);
  return deletedLegacyClientPathRules.some((rule) =>
    matchesRule(pathname, rule),
  );
}

export function isCompanyOwnedPath(value?: string | null): boolean {
  const pathname = clientBoundaryPathname(value);
  return companyOwnedClientPathRules.some((rule) =>
    matchesRule(pathname, rule),
  );
}

export function isSaaSOnlyClientPath(value?: string | null): boolean {
  const pathname = clientBoundaryPathname(value);
  return saasOnlyClientPathRules.some((rule) => matchesRule(pathname, rule));
}

export function isCommunityOnlyClientPath(value?: string | null): boolean {
  const pathname = clientBoundaryPathname(value);
  return communityOnlyClientPathRules.some((rule) =>
    matchesRule(pathname, rule),
  );
}

export function isDebugClientPath(value?: string | null): boolean {
  const pathname = clientBoundaryPathname(value);
  return debugClientPathRules.some((rule) => matchesRule(pathname, rule));
}

export function isModeForbiddenClientPath(
  value: string | null | undefined,
  mode: Edition,
): boolean {
  if (mode === "saas") {
    return isCommunityOnlyClientPath(value);
  }
  return isSaaSOnlyClientPath(value);
}

export function sanitizeClientPostAuthRedirect(
  requestedPath: string | null | undefined,
  mode: Edition,
): string {
  const localPath = normalizeLocalPath(requestedPath);
  const pathname = clientBoundaryPathname(localPath);

  if (pathname === "/login/device") {
    return localPath;
  }
  if (
    isCompanyOwnedPath(pathname) ||
    isDeletedLegacyClientPath(pathname) ||
    isDebugClientPath(pathname) ||
    isModeForbiddenClientPath(pathname, mode)
  ) {
    return "/";
  }

  return localPath;
}
