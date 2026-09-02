# Security model

Nanami is designed around explicit authority, least privilege, and redaction-first diagnostics. The selected packages in this repository illustrate several of those decisions without exposing production configuration.

## Trust boundaries

### Browser to backend

Browsers communicate with a Next.js BFF, not directly with internal service endpoints. The BFF forwards only an allowlisted subset of headers and derives upstream targets from server configuration. Tenant-facing and company-internal routes fail closed when used from the wrong application mode.

### Service to service

Control-plane internal endpoints use mTLS identity. Certificate identity is narrowed by the persisted agent kind and requested operation; request headers and body fields are consistency checks, not sources of authority.

The [`internaltransport`](../pkg/internaltransport) excerpt shows:

- mandatory CA, certificate, and key configuration;
- HTTPS-only internal URL validation;
- minimum TLS version enforcement;
- bounded client-certificate lifetime;
- certificate reloading without process restart;
- tests for malformed, expired, and overlong-lived certificate material.

### Runtime secrets

WireGuard private keys stay client-local. Public keys and policy may cross the control plane, but private material is held behind native secret-store adapters. Diagnostics expose metadata and safe reason codes rather than raw tokens, keys, cookies, DSNs, or certificate bodies.

Short server-side secrets that must be persisted use versioned AES-256-GCM envelopes. The envelope version and key identifier are authenticated data, so changing metadata invalidates decryption rather than silently selecting a different trust context. The [`secretbox`](../pkg/secretbox) excerpt includes round-trip, wrong-key, plaintext, and key-size tests.

### Tenant isolation

Tenant scope is resolved server-side. Shared gateway scheduling requires explicit runtime capability evidence for route-domain isolation and source/destination ACL enforcement. Missing capability evidence is a rejection reason, not an invitation to continue in a degraded sharing mode.

## Failure model

- Mutations create desired state; only runtime observation creates applied-state truth.
- Missing or ambiguous endpoint configuration returns an actionable failure state.
- Proxying a UDP WireGuard hostname is treated as incompatible exposure rather than silently accepted.
- Invalid public origin addresses, DNS mismatches, and unsupported protocols remain visible as typed diagnostics.
- Future capabilities are labeled as unavailable instead of being inferred from adjacent groundwork.

## Public-repository hygiene

This showcase contains no environment dumps, private certificates, credentials, customer data, production database schema, or private Git history. Domain names appear only as public test-environment links or synthetic test fixtures.

Security issues concerning the live environments should be reported privately through the repository owner's GitHub contact rather than opened as a public issue.
