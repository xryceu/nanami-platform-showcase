# Repository scope

The source project is a private monorepo. This public repository was created as an independent export with no shared Git history.

## Selection criteria

Code was included when it met all of these conditions:

1. It represents a meaningful production engineering concern.
2. It has useful tests and can be reviewed independently.
3. It does not disclose credentials, customer data, production infrastructure, or commercially sensitive implementation.
4. It can compile without private services or generated artifacts.
5. It demonstrates decisions that are easier to evaluate in code than in screenshots.

## Why these packages

The selected packages cover a complete Control Plane session-management vertical slice, configuration validation, mTLS, live certificate reload, DNS and endpoint modeling, protocol capability checks, runtime reconciliation, idempotent fan-out, authenticated secret envelopes, production-derived Next.js domain logic from all four web products, and an Electron main/preload/renderer security boundary. Together they form a coherent full-stack engineering sample without reconstructing the full platform.

## Excluded material

- Full control-plane handlers, services, database models, and migrations
- Production Gin/GORM wiring for the exported session slice; the public version uses standard-library HTTP and memory adapters against the same inward-facing contracts
- Complete RBAC and tenant policy implementation
- Gateway scheduling and production dataplane application
- Complete web and native application source trees
- Desktop daemon lifecycle, authenticated local transport, full renderer, platform packaging, signing, notarization, and installer implementation
- Deployment manifests, infrastructure identifiers, and operational scripts
- Entitlement rules, checkout, billing-provider integrations, and non-public commercial policy
- Internal audit, support, and incident-response surfaces
- Private repository commits, branches, issues, and release evidence

The Next.js excerpt keeps Client, Company, Marketing, and Docs code in separate routes and production-derived domain modules. It includes the public pricing payload and the real operator incident taxonomy, but excludes authentication, BFF routes, API clients, private environment configuration, generated backend types, the non-public commercial catalog, and the complete product UI. Fictional records exist only in the review harness where a live backend would normally supply data.

## Provenance

All code and product assets in this repository originate from the repository owner's Nanami project. Exact source mappings and adaptations are recorded in [`frontend/PROVENANCE.md`](../frontend/PROVENANCE.md) and [`electron/PROVENANCE.md`](../electron/PROVENANCE.md).
