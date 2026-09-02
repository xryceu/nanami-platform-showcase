import catalog from "./incident-taxonomy.json";
import type { StatusAudience, StatusSeverity, StatusSource } from "./types";

export type IncidentTaxonomyAudience = StatusAudience | "both";

export type IncidentTaxonomyCategory = {
  id: string;
  label: string;
  description: string;
  defaultSeverity: StatusSeverity;
  primaryOwner: string;
  audience: IncidentTaxonomyAudience;
  sources: StatusSource[];
  relatedCodes: string[];
  symptoms: string[];
  firstChecks: string[];
  remediation: string[];
  evidence: string[];
  userMessageGuidance: string;
};

type IncidentTaxonomyCatalog = {
  version: string;
  categories: IncidentTaxonomyCategory[];
};

const typedCatalog = catalog as IncidentTaxonomyCatalog;

export const INCIDENT_TAXONOMY_VERSION = typedCatalog.version;

export const INCIDENT_TAXONOMY_CATEGORIES = typedCatalog.categories;

export const INCIDENT_TAXONOMY_BY_ID = new Map(
  INCIDENT_TAXONOMY_CATEGORIES.map((category) => [category.id, category]),
);

export function incidentTaxonomyForAudience(audience: StatusAudience) {
  return INCIDENT_TAXONOMY_CATEGORIES.filter(
    (category) =>
      category.audience === "both" || category.audience === audience,
  );
}

export function incidentTaxonomyForCode(
  code: string,
): IncidentTaxonomyCategory | null {
  const normalized = code.trim().toLowerCase();
  if (!normalized) {
    return null;
  }

  return (
    INCIDENT_TAXONOMY_CATEGORIES.find((category) =>
      category.relatedCodes.some((candidate) => {
        const mapped = candidate.toLowerCase();
        return normalized === mapped || normalized.startsWith(mapped);
      }),
    ) ?? null
  );
}
