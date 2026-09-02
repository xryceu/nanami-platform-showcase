import { PricingCatalog } from "@/components/pricing-catalog";
import { getPricingModel } from "@/product/marketing/pricing";
import { getMarketingDestinations } from "@/product/marketing/site-config.server";

const proofPoints = [
  "Policy-aware access to private services",
  "Hosted control plane with observed runtime",
  "Community path for self-hosted evaluation",
];

export default async function MarketingPage() {
  const pricingModel = await getPricingModel(getMarketingDestinations());

  return (
    <main className="marketing-shell">
      <section className="marketing-intro">
        <div>
          <p className="product-name">Nanami Cloud</p>
          <h1>Private networking your team can operate</h1>
          <p>
            Connect people and devices to private services with explicit policy,
            observable gateways, and one operational model.
          </p>
          <div className="marketing-actions">
            <a href="#pricing">Get started</a>
            <a className="secondary-action" href="/client">
              See the product surface
            </a>
          </div>
        </div>
        <ul className="proof-list">
          {proofPoints.map((point) => (
            <li key={point}>{point}</li>
          ))}
        </ul>
      </section>

      <div id="pricing">
        <PricingCatalog model={pricingModel} />
      </div>
    </main>
  );
}
