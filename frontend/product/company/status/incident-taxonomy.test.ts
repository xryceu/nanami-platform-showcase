import { describe, expect, it } from "vitest";

import {
  INCIDENT_TAXONOMY_CATEGORIES,
  INCIDENT_TAXONOMY_VERSION,
  incidentTaxonomyForAudience,
  incidentTaxonomyForCode,
} from "./incident-taxonomy";

const requiredCategoryIds = [
  "auth-session",
  "gateway-runtime",
  "route-path-apply",
  "dns-resolution",
  "access-ownership",
  "audit-log-sink",
  "connector-integration",
  "control-plane-availability",
  "deployment-config-secrets",
  "public-distribution-install",
];

describe("incident taxonomy", () => {
  it("covers the current operational failure modes with playbook metadata", () => {
    expect(INCIDENT_TAXONOMY_VERSION).toMatch(/^\d{4}-\d{2}-\d{2}$/);
    expect(INCIDENT_TAXONOMY_CATEGORIES.map((category) => category.id)).toEqual(
      requiredCategoryIds,
    );

    for (const category of INCIDENT_TAXONOMY_CATEGORIES) {
      expect(category.label).toBeTruthy();
      expect(category.description).toBeTruthy();
      expect(["critical", "warning", "info"]).toContain(
        category.defaultSeverity,
      );
      expect(category.primaryOwner).toBeTruthy();
      expect(["platform", "system", "both"]).toContain(category.audience);
      expect(category.sources.length).toBeGreaterThan(0);
      expect(category.relatedCodes.length).toBeGreaterThan(0);
      expect(category.symptoms.length).toBeGreaterThanOrEqual(2);
      expect(category.firstChecks.length).toBeGreaterThanOrEqual(2);
      expect(category.remediation.length).toBeGreaterThanOrEqual(2);
      expect(category.evidence.length).toBeGreaterThanOrEqual(2);
      expect(category.userMessageGuidance).toBeTruthy();
    }
  });

  it("keeps platform-only playbooks out of system audience guidance", () => {
    const systemCategories = incidentTaxonomyForAudience("system");
    const platformOnly = systemCategories.filter(
      (category) => category.audience === "platform",
    );

    expect(platformOnly).toEqual([]);
    expect(systemCategories.map((category) => category.id)).toContain(
      "gateway-runtime",
    );
    expect(systemCategories.map((category) => category.id)).toContain(
      "public-distribution-install",
    );
  });

  it("maps current incident codes back to taxonomy categories", () => {
    expect(incidentTaxonomyForCode("gateway_down")?.id).toBe("gateway-runtime");
    expect(incidentTaxonomyForCode("ROUTE_APPLY_FAILED")?.id).toBe(
      "route-path-apply",
    );
    expect(incidentTaxonomyForCode("audit_log_sink_delivery_failed")?.id).toBe(
      "audit-log-sink",
    );
    expect(incidentTaxonomyForCode("routing_apply_mode_off")?.id).toBe(
      "route-path-apply",
    );
    expect(incidentTaxonomyForCode("audit_error_spike:auth.login")?.id).toBe(
      "auth-session",
    );
  });
});
