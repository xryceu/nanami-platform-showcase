import type { DeploymentType } from "../shared/contracts.js";
import { normalizeDeploymentBaseUrl } from "./deployment.js";

export type DeploymentLinkTarget = {
  type: DeploymentType;
  baseUrl: string;
};

export function parseDeploymentLink(
  value: string,
): DeploymentLinkTarget | null {
  let link: URL;
  try {
    link = new URL(String(value ?? "").trim());
  } catch {
    return null;
  }

  if (
    link.protocol !== "nanami:" ||
    link.hostname !== "target" ||
    (link.pathname !== "" && link.pathname !== "/") ||
    link.username ||
    link.password ||
    link.hash
  ) {
    return null;
  }

  const allowedParameters = new Set(["server", "edition"]);
  if (
    [...link.searchParams.keys()].some((key) => !allowedParameters.has(key))
  ) {
    return null;
  }

  const server = link.searchParams.get("server")?.trim() ?? "";
  const edition = link.searchParams.get("edition")?.trim().toLowerCase() ?? "";
  if (
    !server ||
    !["saas", "community", "enterprise_self_hosted"].includes(edition)
  ) {
    return null;
  }

  const type: DeploymentType = edition === "saas" ? "cloud" : "self_hosted";
  try {
    return { type, baseUrl: normalizeDeploymentBaseUrl(type, server) };
  } catch {
    return null;
  }
}

export function deploymentLinkFromArguments(
  argv: readonly string[],
): string | null {
  return argv.find((value) => String(value).startsWith("nanami://")) ?? null;
}
