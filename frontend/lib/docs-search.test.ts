import { describe, expect, it } from "vitest";

import {
  adjacentDocs,
  groupDocs,
  normalizeSearchText,
  rankDocEntry,
  searchDocs,
  snippetAround,
  type DocEntry,
} from "./docs-search";

const entries: DocEntry[] = [
  {
    href: "/docs/connect-device",
    title: "Connect your first device",
    section: "Get started",
    description: "Enroll a device and verify its observed connection state.",
    headings: [
      {
        title: "Verify connection",
        anchor: "verify-connection",
        text: "Wait until the device reports an observed online state.",
      },
    ],
    body: "Install the CLI, sign in, choose a network, and connect.",
    aliases: ["enroll node", "first connection"],
    keywords: ["device", "cli", "connect"],
  },
  {
    href: "/docs/no-eligible-gateway",
    title: "No eligible gateway",
    section: "Troubleshooting",
    description: "Resolve missing, stale, offline, or permission states.",
    headings: [
      {
        title: "Check observed runtime",
        anchor: "check-runtime",
        text: "Confirm that at least one permitted gateway reports online.",
      },
    ],
    body: "Do not treat desired gateway configuration as runtime success.",
    aliases: ["cannot connect", "gateway unavailable"],
    keywords: ["gateway", "offline", "reason code"],
  },
  {
    href: "/docs/private-dns",
    title: "Configure private DNS",
    section: "Configure",
    description: "Apply resolver policy to an enrolled network.",
    headings: [],
    body: "Create a DNS zone, add records, and inspect applied device state.",
    aliases: ["dns zone"],
    keywords: ["resolver", "record"],
  },
];

describe("documentation search", () => {
  it("normalizes case, apostrophes, whitespace, and Russian ё", () => {
    expect(normalizeSearchText("  Don’t   Ёлка  ")).toBe("dont елка");
  });

  it("ranks an exact title above an alias match", () => {
    const results = searchDocs(entries, "no eligible gateway");
    expect(results[0]?.href).toBe("/docs/no-eligible-gateway");
    expect(results[0]?.score).toBeGreaterThan(100);
  });

  it("links directly to a matching heading", () => {
    const result = rankDocEntry(entries[1]!, "observed runtime");
    expect(result?.resultHref).toBe("/docs/no-eligible-gateway#check-runtime");
    expect(result?.matchHeading).toBe("Check observed runtime");
  });

  it("requires every search term to be present", () => {
    expect(rankDocEntry(entries[2]!, "dns impossible-term")).toBeNull();
  });

  it("creates bounded snippets around a match", () => {
    const source = `${"prefix ".repeat(30)}gateway offline${" suffix".repeat(30)}`;
    const snippet = snippetAround(source, "gateway", 80);
    expect(snippet).toContain("gateway offline");
    expect(snippet.length).toBeLessThanOrEqual(82);
  });

  it("preserves section order while grouping entries", () => {
    expect(groupDocs(entries).map((section) => section.title)).toEqual([
      "Get started",
      "Troubleshooting",
      "Configure",
    ]);
  });

  it("returns safe adjacent navigation at both boundaries", () => {
    expect(adjacentDocs(entries, entries[0]!.href).previous).toBeNull();
    expect(adjacentDocs(entries, entries[2]!.href).next).toBeNull();
    expect(adjacentDocs(entries, "/missing").current).toBeNull();
  });
});
