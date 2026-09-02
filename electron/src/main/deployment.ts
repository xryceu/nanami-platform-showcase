import crypto from "node:crypto";

import type { DeploymentProfile, DeploymentType } from "../shared/contracts.js";

export function normalizeDeploymentBaseUrl(
  type: DeploymentType,
  value?: string,
): string {
  const raw = String(value ?? "").trim();
  if (!raw) {
    throw new Error(
      type === "cloud"
        ? "The cloud endpoint is not configured."
        : "Enter the Nanami server address.",
    );
  }

  let parsed: URL;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error("Enter a valid Nanami server address.");
  }

  if (parsed.username || parsed.password || parsed.hash || parsed.search) {
    throw new Error(
      "The Nanami server address cannot contain credentials, query parameters, or fragments.",
    );
  }

  const loopback = ["localhost", "127.0.0.1", "::1"].includes(
    parsed.hostname.toLowerCase(),
  );
  const allowLoopbackHttp =
    process.env.NANAMI_DESKTOP_ALLOW_INSECURE_LOOPBACK === "1";
  if (
    parsed.protocol !== "https:" &&
    !(parsed.protocol === "http:" && loopback && allowLoopbackHttp)
  ) {
    throw new Error("Use an HTTPS Nanami server address.");
  }
  if (!parsed.hostname) {
    throw new Error("Enter a valid Nanami server address.");
  }

  parsed.hostname = parsed.hostname.toLowerCase();
  parsed.pathname = parsed.pathname.replace(/\/{2,}/g, "/").replace(/\/$/, "");
  return parsed.toString().replace(/\/$/, "");
}

export function deploymentContextName(baseUrl: string): string {
  return `desktop-${crypto.createHash("sha256").update(baseUrl).digest("hex").slice(0, 20)}`;
}

export function safeDeploymentProfile(
  value: unknown,
): DeploymentProfile | null {
  if (!value || typeof value !== "object") {
    return null;
  }

  const candidate = value as Partial<DeploymentProfile>;
  if (candidate.type !== "cloud" && candidate.type !== "self_hosted") {
    return null;
  }

  try {
    const baseUrl = normalizeDeploymentBaseUrl(
      candidate.type,
      candidate.baseUrl,
    );
    if (!candidate.deploymentId || !candidate.capabilitiesVersion) {
      return null;
    }
    return {
      id: deploymentContextName(baseUrl),
      type: candidate.type,
      displayName:
        String(candidate.displayName ?? "").trim() || new URL(baseUrl).hostname,
      baseUrl,
      normalizedOrigin: new URL(baseUrl).origin,
      deploymentId: String(candidate.deploymentId).trim(),
      edition: ["saas", "community", "enterprise_self_hosted"].includes(
        String(candidate.edition),
      )
        ? (candidate.edition as DeploymentProfile["edition"])
        : "unknown",
      lastValidatedAt: String(candidate.lastValidatedAt ?? ""),
      capabilitiesVersion: String(candidate.capabilitiesVersion).trim(),
    };
  } catch {
    return null;
  }
}
