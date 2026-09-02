import {
  INCIDENT_TAXONOMY_CATEGORIES,
  INCIDENT_TAXONOMY_VERSION,
} from "@/product/company/status/incident-taxonomy";

export default function CompanyPage() {
  return (
    <main className="page-shell">
      <header className="page-header">
        <div>
          <p className="product-name">Nanami Company Dashboard</p>
          <h1>Operational incident taxonomy</h1>
          <p className="page-description">
            The production catalog used to map runtime reason codes to owners,
            first checks, remediation, and safe user communication.
          </p>
        </div>
        <span className="environment-label">
          Catalog {INCIDENT_TAXONOMY_VERSION}
        </span>
      </header>

      <section className="panel support-panel" aria-labelledby="taxonomy-title">
        <div className="section-heading support-heading">
          <div>
            <h2 id="taxonomy-title">Failure-mode catalog</h2>
            <p>{INCIDENT_TAXONOMY_CATEGORIES.length} operational categories</p>
          </div>
        </div>

        <div className="support-list">
          {INCIDENT_TAXONOMY_CATEGORIES.map((category) => (
            <article className="support-row" key={category.id}>
              <div className="support-case-main">
                <div className="support-case-title">
                  <strong>{category.label}</strong>
                  <span
                    className={`severity severity-${category.defaultSeverity}`}
                  >
                    {category.defaultSeverity}
                  </span>
                </div>
                <p>{category.description}</p>
                <span className="support-context">
                  Codes: {category.relatedCodes.join(", ")}
                </span>
              </div>
              <div className="support-state">
                <span className="status status-neutral">
                  {category.audience}
                </span>
                <strong>{category.primaryOwner}</strong>
                <span className="support-context">
                  {category.firstChecks.length} checks ·{" "}
                  {category.remediation.length} actions
                </span>
              </div>
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
