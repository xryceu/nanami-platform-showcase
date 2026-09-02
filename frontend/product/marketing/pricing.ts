import catalog from "./public-pricing-catalog.json";

export type BillingPeriodKey =
  "monthly" | "quarterly" | "six_months" | "yearly";

export type PricingTierKey =
  "open_source" | "free" | "starter" | "team" | "scale" | "enterprise";

export type BillingPeriod = {
  key: BillingPeriodKey;
  label: string;
  shortLabel: string;
  months: number;
  discountPercent: number;
  helper: string;
};

export type PricingTier = {
  key: PricingTierKey;
  sortOrder: number;
  name: string;
  monthlyPriceCents: number | null;
  cta: { label: string; href: string; external?: boolean };
  recommended?: boolean;
  metrics: Array<{ label: string; value: string }>;
};

export type PricingModel = {
  periods: BillingPeriod[];
  defaultPeriod: BillingPeriodKey;
  tiers: PricingTier[];
  catalogSource: string;
};

type PublicPricingTier = Omit<PricingTier, "cta"> & {
  shortName: string;
  eyebrow: string;
  description: string;
  priceNote: string;
  billingNote: string;
  cta: {
    kind: "signup" | "contact" | "docs_self_hosted";
    label: string;
    external?: boolean;
  };
  visibility: string;
  lifecycleState: string;
};

type PublicPricingCatalog = {
  source: string;
  lifecycleState: string;
  defaultBillingPeriod: BillingPeriodKey;
  billingPeriods: BillingPeriod[];
  tiers: PublicPricingTier[];
  comparisonGroups: unknown[];
  booster: Record<string, unknown>;
};

const repoPricingCatalog = catalog.publicPricingCatalog as PublicPricingCatalog;

export const pricingCatalogUrlEnv = "MARKETING_PRICING_CATALOG_URL";
export const publicPricingCatalogPath = "/api/v1/public/pricing-catalog";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isPublicPricingCatalog(value: unknown): value is PublicPricingCatalog {
  return Boolean(
    isRecord(value) &&
    value.lifecycleState === "published" &&
    typeof value.source === "string" &&
    Array.isArray(value.billingPeriods) &&
    Array.isArray(value.tiers),
  );
}

function paidPlanEntryUrl(signupUrl: string, tierKey: PricingTierKey): string {
  if (tierKey === "free") {
    return signupUrl;
  }

  const destination = `/settings/plan?requestedPlan=${tierKey}`;
  const isRootRelative =
    signupUrl.startsWith("/") && !signupUrl.startsWith("//");
  const url = new URL(signupUrl, "https://showcase.local");
  url.searchParams.set("redirect", destination);
  return isRootRelative
    ? `${url.pathname}${url.search}${url.hash}`
    : url.toString();
}

function ctaHref(
  tier: PublicPricingTier,
  urls: { signupUrl: string; contactUrl: string; docsSelfHostedUrl: string },
): string {
  switch (tier.cta.kind) {
    case "signup":
      return paidPlanEntryUrl(urls.signupUrl, tier.key);
    case "contact":
      return urls.contactUrl;
    case "docs_self_hosted":
      return urls.docsSelfHostedUrl;
  }
}

async function loadConfiguredPricingCatalog(): Promise<PublicPricingCatalog | null> {
  const catalogUrl = process.env[pricingCatalogUrlEnv]?.trim();
  if (!catalogUrl) {
    return null;
  }

  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 1500);

  try {
    const response = await fetch(catalogUrl, {
      headers: { accept: "application/json" },
      signal: controller.signal,
      next: { revalidate: 300 },
    });
    if (!response.ok) {
      return null;
    }
    const body: unknown = await response.json();
    return isPublicPricingCatalog(body) ? body : null;
  } catch {
    return null;
  } finally {
    clearTimeout(timeout);
  }
}

export async function getPricingModel(options: {
  signupUrl: string;
  contactUrl: string;
  docsSelfHostedUrl: string;
}): Promise<PricingModel> {
  const configuredCatalog = await loadConfiguredPricingCatalog();
  return createPricingModel(options, configuredCatalog ?? repoPricingCatalog);
}

export function createPricingModel(
  urls: { signupUrl: string; contactUrl: string; docsSelfHostedUrl: string },
  sourceCatalog: PublicPricingCatalog = repoPricingCatalog,
): PricingModel {
  const tiers = sourceCatalog.tiers
    .filter(
      (tier) =>
        tier.visibility === "public" && tier.lifecycleState === "published",
    )
    .sort((left, right) => left.sortOrder - right.sortOrder)
    .map((tier) => ({
      key: tier.key,
      sortOrder: tier.sortOrder,
      name: tier.name,
      monthlyPriceCents: tier.monthlyPriceCents,
      recommended: tier.recommended,
      metrics: tier.metrics.map(({ label, value }) => ({ label, value })),
      cta: {
        label: tier.cta.label,
        href: ctaHref(tier, urls),
        external: tier.cta.external,
      },
    }));

  return {
    periods: sourceCatalog.billingPeriods.filter(
      (period) => period.key === "monthly" || period.key === "yearly",
    ),
    defaultPeriod: sourceCatalog.defaultBillingPeriod,
    tiers,
    catalogSource: sourceCatalog.source,
  };
}
