# Guided code tour

This tour is designed for a focused 10-minute technical review.

## 1. Control Plane vertical slice

Start with [`internal/session`](../internal/session) and the [Go architecture walkthrough](go-architecture.md). This is a complete application slice rather than an isolated helper: domain types, use cases, reader/revoker ports, HTTP and memory adapters, and a composition root.

The key review point is dependency direction. HTTP and persistence depend inward on the application/domain contracts; the session service does not import either adapter.

## 2. Internal transport boundary

Start with [`control_plane_url.go`](../pkg/internaltransport/control_plane_url.go). It rejects ambiguous mixed configuration and normalizes one canonical internal endpoint.

Continue with [`mtls_http_client.go`](../pkg/internaltransport/mtls_http_client.go). The client validates its trust material, enforces a maximum certificate lifetime, and reloads changed certificates safely under concurrent requests.

The adjacent tests generate certificate fixtures at runtime, so no key material is committed.

## 3. Gateway endpoint safety

Open [`endpoint_dns.go`](../pkg/gatewayruntime/endpoint_dns.go). The builder turns configuration plus DNS observations into a redaction-safe runtime model with status, explanation, and remediation fields.

Then review [`endpoint_exposure.go`](../pkg/gatewayruntime/endpoint_exposure.go) and [`protocol_registry.go`](../pkg/gatewayruntime/protocol_registry.go). They keep DNS publication, transport exposure, and protocol capability as separate concepts.

The tests cover public/private IP handling, mismatched DNS records, IPv4/IPv6 behavior, direct versus proxied exposure, and protocol-specific constraints.

## 4. Desired-state reconciliation

[`reconcile.go`](../pkg/gatewayruntime/reconcile.go) models the apply boundary. Its tests focus on deterministic reconciliation and cleanup behavior, including stale runtime resources. This is where desired configuration becomes an explicit apply plan rather than an optimistic status flag.

## 5. Idempotent fan-out

[`fanout.go`](../pkg/gatewayorchestration/fanout.go) hashes the stable parts of a desired-state snapshot into a public-safe ID, then scopes its idempotency key to a daemon and worker. The tests show that identical snapshots remain stable while a revision change produces a different delivery identity.

## 6. Authenticated secret storage

[`secretbox.go`](../pkg/secretbox/secretbox.go) wraps short server-side secrets in versioned AES-256-GCM envelopes. The version and key ID are authenticated alongside the ciphertext, and malformed plaintext, wrong keys, and invalid key sizes fail closed.

## 7. Production-derived Next.js product surfaces

Start with [`frontend/PROVENANCE.md`](../frontend/PROVENANCE.md). It maps each excerpt to the exact private-monorepo source and lists the minimal adaptation made for this standalone build.

### Client App

[`client-boundary.ts`](../frontend/product/client/client-boundary.ts) is the production route-ownership and post-auth redirect boundary for SaaS and Community editions. [`transport-contract.ts`](../frontend/product/client/transport-contract.ts) derives direct, gateway-mediated, pending, and offline states from observed connectivity rather than desired configuration.

### Company Dashboard

[`incident-taxonomy.json`](../frontend/product/company/status/incident-taxonomy.json) is the production failure-mode catalog used by Company status surfaces. The adjacent [filter](../frontend/product/company/status/filters.ts) and [merge](../frontend/product/company/status/merge.ts) modules handle time windows, source failures, severity ordering, and deterministic consolidation of concurrent status collectors.

### Marketing

[`pricing.ts`](../frontend/product/marketing/pricing.ts) validates an optional remote public catalog, fails back to the repository catalog after a bounded request, filters unpublished tiers, and routes each CTA through its declared commercial action. The review component retains the production total and monthly-equivalent calculations.

### Documentation

[`docs-nav.ts`](../frontend/product/docs/docs-nav.ts) contains the real English/Russian information architecture and pure navigation contracts used by the documentation site. The review browser consumes that catalog directly, then adds interactive search without a private generated search index.

The files under `frontend/app` and `frontend/components` are intentionally a thin review harness, not a claim that the complete four applications are public. Tests verify the selected contracts but are supplementary to the production-derived implementation.

## 8. Electron privilege boundary

Start with [`electron/PROVENANCE.md`](../electron/PROVENANCE.md), then follow the call path from the typed [`DesktopBridge`](../electron/src/shared/contracts.ts) through the isolated [`preload`](../electron/src/preload/index.ts) to the main-process [`IPC registry`](../electron/src/main/ipc.ts).

[`window-policy.ts`](../electron/src/main/window-policy.ts) keeps Node integration out of the renderer, enables context isolation and sandboxing, rejects renderer navigation, and opens only validated HTTPS targets externally. [`deployment.ts`](../electron/src/main/deployment.ts) and [`deployment-link.ts`](../electron/src/main/deployment-link.ts) treat server selection as a security boundary instead of accepting an arbitrary URL from the renderer.

The renderer excerpt keeps state derivation and diagnostic sanitization pure. The Windows tray module independently derives available native actions from authenticated, observed runtime state, so a desired connection never masquerades as a connected tunnel.

## Questions this code is meant to support

- How should desired and observed state converge when runtime agents are intermittently unavailable?
- How can retries stay idempotent without placing raw desired state in an identifier?
- Where should service identity be established and narrowed?
- How should encrypted values carry enough metadata for key rotation without trusting that metadata blindly?
- What evidence is required before a gateway can safely serve multiple tenants?
- How should operational failures become safe, actionable product diagnostics?
- Where should a Next.js server/client boundary sit when only one interaction requires browser state?
- Which presentation logic belongs in pure modules rather than React components?
- How can four web products share a platform identity without erasing their different trust and user boundaries?
- Which Electron responsibilities belong in renderer, isolated preload, main process, and native adapters?
