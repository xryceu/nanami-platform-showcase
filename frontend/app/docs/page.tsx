import { DocsBrowser } from "@/components/docs-browser";
import type { DocEntry } from "@/lib/docs-search";
import { getLocalizedDocs } from "@/product/docs/docs-nav";

const docs: DocEntry[] = getLocalizedDocs("en").map((entry) => ({
  href: `/docs?article=${entry.slug.join("/")}`,
  title: entry.title,
  section: entry.sectionTitle,
  description: entry.summary,
  headings: [],
  body: entry.summary,
  aliases: [],
  keywords: entry.keywords,
}));

export default function DocsPage() {
  return (
    <main className="docs-shell">
      <DocsBrowser entries={docs} />
    </main>
  );
}
