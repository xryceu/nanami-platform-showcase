import { afterEach, describe, expect, it } from "vitest";

import {
  deploymentContextName,
  normalizeDeploymentBaseUrl,
} from "../src/main/deployment.js";
import { parseDeploymentLink } from "../src/main/deployment-link.js";

afterEach(() => {
  delete process.env.NANAMI_DESKTOP_ALLOW_INSECURE_LOOPBACK;
});

describe("deployment boundary", () => {
  it("normalizes HTTPS targets and creates stable deployment-scoped identities", () => {
    const baseUrl = normalizeDeploymentBaseUrl(
      "self_hosted",
      "https://EDGE.Example.test//api/",
    );
    expect(baseUrl).toBe("https://edge.example.test/api");
    expect(deploymentContextName(baseUrl)).toMatch(/^desktop-[a-f0-9]{20}$/);
    expect(deploymentContextName(baseUrl)).toBe(deploymentContextName(baseUrl));
  });

  it.each([
    "http://edge.example.test",
    "https://user:password@edge.example.test",
    "https://edge.example.test?token=private",
    "https://edge.example.test/#fragment",
  ])("rejects unsafe target %s", (target) => {
    expect(() => normalizeDeploymentBaseUrl("self_hosted", target)).toThrow();
  });

  it("allows HTTP only for an explicitly enabled loopback development target", () => {
    process.env.NANAMI_DESKTOP_ALLOW_INSECURE_LOOPBACK = "1";
    expect(
      normalizeDeploymentBaseUrl("self_hosted", "http://127.0.0.1:8080"),
    ).toBe("http://127.0.0.1:8080");
  });

  it("accepts an allowlisted deep-link shape and rejects extra parameters", () => {
    expect(
      parseDeploymentLink(
        "nanami://target?server=https%3A%2F%2Fedge.example.test&edition=community",
      ),
    ).toEqual({ type: "self_hosted", baseUrl: "https://edge.example.test" });
    expect(
      parseDeploymentLink(
        "nanami://target?server=https%3A%2F%2Fedge.example.test&edition=community&next=unsafe",
      ),
    ).toBeNull();
  });
});
