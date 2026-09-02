import Link from "next/link";

const surfaces = [
  {
    href: "/client",
    name: "Client App",
    role: "Tenant product",
    description:
      "Production route-boundary and direct/gateway transport contracts, wired into a small review harness.",
    code: "product/client",
  },
  {
    href: "/company",
    name: "Company Dashboard",
    role: "Internal operations",
    description:
      "Production incident taxonomy, severity ordering, response merging, ownership, and remediation guidance.",
    code: "product/company/status",
  },
  {
    href: "/marketing",
    name: "Marketing",
    role: "Public product surface",
    description:
      "Production catalog validation, remote fallback behavior, CTA routing, and exact billing totals.",
    code: "product/marketing",
  },
  {
    href: "/docs",
    name: "Documentation",
    role: "Public operator guide",
    description:
      "Production bilingual information architecture and navigation, paired with the searchable review UI.",
    code: "product/docs",
  },
];

export default function HomePage() {
  return (
    <main className="page-shell">
      <header className="page-header">
        <div>
          <p className="product-name">Production-derived frontend source</p>
          <h1>Four product surfaces, one platform</h1>
          <p className="page-description">
            Real modules exported from Nanami&apos;s Client, Company, Marketing,
            and Documentation applications with private integrations removed.
          </p>
        </div>
        <span className="environment-label">Next.js 16 · React 19</span>
      </header>

      <section
        className="surface-directory"
        aria-labelledby="surface-directory-title"
      >
        <div className="directory-heading">
          <h2 id="surface-directory-title">Included surfaces</h2>
          <p>Each route keeps its own UI and presentation logic.</p>
        </div>
        <div className="directory-list">
          {surfaces.map((surface, index) => (
            <Link
              className="directory-row"
              href={surface.href}
              key={surface.href}
            >
              <span className="directory-index" aria-hidden="true">
                {String(index + 1).padStart(2, "0")}
              </span>
              <span className="directory-main">
                <strong>{surface.name}</strong>
                <span>{surface.description}</span>
              </span>
              <span className="directory-meta">
                <span>{surface.role}</span>
                <code>{surface.code}</code>
              </span>
              <span className="directory-arrow" aria-hidden="true">
                →
              </span>
            </Link>
          ))}
        </div>
      </section>
    </main>
  );
}
