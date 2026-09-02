export type DocHeading = {
  title: string;
  anchor: string;
  text: string;
};

export type DocEntry = {
  href: string;
  title: string;
  section: string;
  description: string;
  headings: DocHeading[];
  body: string;
  aliases: string[];
  keywords: string[];
};

export type DocSearchResult = DocEntry & {
  resultHref: string;
  matchHeading: string | null;
  snippet: string;
  score: number;
};

export type DocSection = {
  title: string;
  entries: DocEntry[];
};

export function normalizeSearchText(value: string): string {
  return value
    .toLocaleLowerCase()
    .replace(/[’']/g, "")
    .replace(/ё/g, "е")
    .replace(/\s+/g, " ")
    .trim();
}

export function snippetAround(
  text: string,
  query: string,
  maximumLength = 180,
): string {
  const compact = text.replace(/\s+/g, " ").trim();
  if (compact.length <= maximumLength) {
    return compact;
  }

  const location = normalizeSearchText(compact).indexOf(
    normalizeSearchText(query),
  );
  if (location < 0) {
    return `${compact.slice(0, maximumLength).trim()}…`;
  }

  const contextBefore = Math.floor(maximumLength * 0.35);
  const start = Math.max(0, location - contextBefore);
  const end = Math.min(compact.length, start + maximumLength);
  return `${start > 0 ? "…" : ""}${compact.slice(start, end).trim()}${
    end < compact.length ? "…" : ""
  }`;
}

export function rankDocEntry(
  entry: DocEntry,
  rawQuery: string,
): DocSearchResult | null {
  const query = normalizeSearchText(rawQuery);
  if (!query) {
    return null;
  }

  const terms = query.split(" ").filter(Boolean);
  const title = normalizeSearchText(entry.title);
  const description = normalizeSearchText(entry.description);
  const aliases = normalizeSearchText(entry.aliases.join(" "));
  const keywords = normalizeSearchText(entry.keywords.join(" "));
  const body = normalizeSearchText(entry.body);
  const heading = entry.headings.find((candidate) =>
    terms.every((term) =>
      normalizeSearchText(`${candidate.title} ${candidate.text}`).includes(
        term,
      ),
    ),
  );
  const headingText = normalizeSearchText(
    entry.headings
      .map((candidate) => `${candidate.title} ${candidate.text}`)
      .join(" "),
  );
  const fullText = [
    title,
    description,
    aliases,
    keywords,
    headingText,
    body,
  ].join(" ");

  if (!terms.every((term) => fullText.includes(term))) {
    return null;
  }

  let score = 1;
  if (title === query) {
    score += 100;
  } else if (title.includes(query)) {
    score += 50;
  }
  if (aliases.includes(query)) {
    score += 35;
  }
  if (heading) {
    score += 22;
  }
  if (keywords.includes(query)) {
    score += 18;
  }
  if (description.includes(query)) {
    score += 12;
  }

  const source = heading?.text || entry.description || entry.body;
  return {
    ...entry,
    resultHref: heading ? `${entry.href}#${heading.anchor}` : entry.href,
    matchHeading: heading?.title ?? null,
    snippet: snippetAround(source, rawQuery),
    score,
  };
}

export function searchDocs(
  entries: DocEntry[],
  query: string,
  limit = 8,
): DocSearchResult[] {
  return entries
    .map((entry) => rankDocEntry(entry, query))
    .filter((entry): entry is DocSearchResult => entry !== null)
    .sort(
      (left, right) =>
        right.score - left.score || left.title.localeCompare(right.title),
    )
    .slice(0, limit);
}

export function groupDocs(entries: DocEntry[]): DocSection[] {
  const sections = new Map<string, DocEntry[]>();
  for (const entry of entries) {
    const sectionEntries = sections.get(entry.section) ?? [];
    sectionEntries.push(entry);
    sections.set(entry.section, sectionEntries);
  }

  return [...sections.entries()].map(([title, sectionEntries]) => ({
    title,
    entries: sectionEntries,
  }));
}

export function adjacentDocs(
  entries: DocEntry[],
  href: string,
): {
  previous: DocEntry | null;
  current: DocEntry | null;
  next: DocEntry | null;
} {
  const currentIndex = entries.findIndex((entry) => entry.href === href);
  if (currentIndex < 0) {
    return { previous: null, current: null, next: null };
  }

  return {
    previous: currentIndex > 0 ? entries[currentIndex - 1]! : null,
    current: entries[currentIndex]!,
    next: currentIndex < entries.length - 1 ? entries[currentIndex + 1]! : null,
  };
}
