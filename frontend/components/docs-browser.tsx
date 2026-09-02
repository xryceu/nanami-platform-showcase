"use client";

import { useMemo, useState } from "react";

import { groupDocs, searchDocs, type DocEntry } from "@/lib/docs-search";

export function DocsBrowser({ entries }: { entries: DocEntry[] }) {
  const [query, setQuery] = useState("");
  const trimmedQuery = query.trim();
  const results = useMemo(
    () => (trimmedQuery ? searchDocs(entries, trimmedQuery) : []),
    [entries, trimmedQuery],
  );
  const sections = useMemo(() => groupDocs(entries), [entries]);

  return (
    <div className="docs-layout">
      <aside className="docs-sidebar">
        <strong>Documentation</strong>
        <nav aria-label="Documentation sections">
          {sections.map((section) => (
            <div key={section.title}>
              <span>{section.title}</span>
              {section.entries.map((entry) => (
                <a href={entry.href} key={entry.href}>
                  {entry.title}
                </a>
              ))}
            </div>
          ))}
        </nav>
      </aside>

      <section className="docs-content" aria-labelledby="docs-heading">
        <div className="docs-intro">
          <h1 id="docs-heading">What do you want to do?</h1>
          <p>
            Search by task, symptom, command, or the reason shown by the
            product.
          </p>
          <label className="docs-search-field">
            <span className="visually-hidden">Search documentation</span>
            <input
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Try “gateway unavailable” or “private DNS”"
              type="search"
              value={query}
            />
            <kbd>/</kbd>
          </label>
        </div>

        {trimmedQuery ? (
          <div className="docs-results" aria-live="polite">
            <div className="docs-results-heading">
              <h2>Search results</h2>
              <span>{results.length} matches</span>
            </div>
            {results.length === 0 ? (
              <div className="empty-state">
                <strong>No documentation matched this search</strong>
                <p>Try a shorter symptom, feature name, or reason code.</p>
              </div>
            ) : (
              results.map((result) => (
                <a
                  className="docs-result"
                  href={result.resultHref}
                  key={result.resultHref}
                >
                  <span>{result.section}</span>
                  <strong>{result.title}</strong>
                  {result.matchHeading ? <em>{result.matchHeading}</em> : null}
                  <p>{result.snippet}</p>
                </a>
              ))
            )}
          </div>
        ) : (
          <div className="docs-task-sections">
            {sections.map((section) => (
              <section key={section.title}>
                <h2>{section.title}</h2>
                <div>
                  {section.entries.map((entry) => (
                    <a href={entry.href} key={entry.href}>
                      <strong>{entry.title}</strong>
                      <span>{entry.description}</span>
                    </a>
                  ))}
                </div>
              </section>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
