"use client";

import { useMemo, useState } from "react";

import {
  type BillingPeriod,
  type BillingPeriodKey,
  type PricingModel,
  type PricingTier,
} from "@/product/marketing/pricing";

function formatCurrency(cents: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(cents / 100);
}

function periodTotalCents(tier: PricingTier, period: BillingPeriod): number {
  if (tier.monthlyPriceCents === null) return 0;
  const gross = tier.monthlyPriceCents * period.months;
  return Math.round(gross * (1 - period.discountPercent / 100));
}

function metricValue(tier: PricingTier, label: RegExp): string {
  return tier.metrics.find((metric) => label.test(metric.label))?.value ?? "—";
}

export function PricingCatalog({ model }: { model: PricingModel }) {
  const [periodKey, setPeriodKey] = useState<BillingPeriodKey>(
    model.defaultPeriod,
  );
  const period = useMemo(
    () =>
      model.periods.find((candidate) => candidate.key === periodKey) ??
      model.periods[0],
    [model.periods, periodKey],
  );

  if (!period) return null;

  return (
    <section className="pricing-section" aria-labelledby="pricing-title">
      <div className="pricing-heading">
        <div>
          <h2 id="pricing-title">Nanami Cloud plans</h2>
          <p>Published catalog values with the exact selected-period total.</p>
        </div>
        <div className="period-switch" aria-label="Billing period">
          {model.periods.map((candidate) => (
            <button
              aria-pressed={periodKey === candidate.key}
              className={periodKey === candidate.key ? "is-active" : undefined}
              key={candidate.key}
              onClick={() => setPeriodKey(candidate.key)}
              type="button"
            >
              {candidate.label}
              {candidate.discountPercent > 0 ? (
                <span>Save {candidate.discountPercent}%</span>
              ) : null}
            </button>
          ))}
        </div>
      </div>

      <div className="pricing-grid">
        {model.tiers.map((tier) => {
          const total = periodTotalCents(tier, period);
          const monthlyEquivalent = Math.round(total / period.months);

          return (
            <article
              className={
                tier.recommended
                  ? "price-column is-recommended"
                  : "price-column"
              }
              key={tier.key}
            >
              <div className="price-column-header">
                <div>
                  <h3>{tier.name}</h3>
                  {tier.recommended ? <span>Recommended</span> : null}
                </div>
                <p>{tier.metrics.map((metric) => metric.value).join(" · ")}</p>
              </div>

              <div className="price-block">
                <strong>
                  {tier.monthlyPriceCents === null
                    ? "Self-hosted"
                    : formatCurrency(monthlyEquivalent)}
                </strong>
                <span>per month equivalent</span>
                <p>
                  {tier.monthlyPriceCents === null
                    ? "Operated by your team"
                    : `${formatCurrency(total)} billed for ${period.months} month${period.months === 1 ? "" : "s"}`}
                </p>
              </div>

              <dl className="plan-limits">
                <div>
                  <dt>Users</dt>
                  <dd>{metricValue(tier, /users/i)}</dd>
                </div>
                <div>
                  <dt>Devices</dt>
                  <dd>{metricValue(tier, /devices/i)}</dd>
                </div>
                <div>
                  <dt>Catalog source</dt>
                  <dd>{model.catalogSource}</dd>
                </div>
              </dl>

              <ul className="feature-list">
                {tier.metrics.map((metric) => (
                  <li key={metric.label}>
                    {metric.label}: {metric.value}
                  </li>
                ))}
              </ul>

              <a className="plan-action" href={tier.cta.href}>
                {tier.cta.label}
              </a>
            </article>
          );
        })}
      </div>

      <p className="activation-note">
        Paid plans use assisted activation in this excerpt. No payment or
        checkout completion is implied.
      </p>
    </section>
  );
}
