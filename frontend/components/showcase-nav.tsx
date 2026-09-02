import Link from "next/link";

const links = [
  { href: "/client", label: "Client" },
  { href: "/company", label: "Company" },
  { href: "/marketing", label: "Marketing" },
  { href: "/docs", label: "Docs" },
];

export function ShowcaseNav() {
  return (
    <header className="showcase-nav">
      <nav aria-label="Frontend showcase">
        <Link className="showcase-brand" href="/">
          Nanami
        </Link>
        <div className="showcase-links">
          {links.map((link) => (
            <Link href={link.href} key={link.href}>
              {link.label}
            </Link>
          ))}
        </div>
        <a
          className="source-link"
          href="https://github.com/xryceu/nanami-platform-showcase"
        >
          Source
        </a>
      </nav>
    </header>
  );
}
