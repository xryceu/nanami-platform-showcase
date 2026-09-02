import { AccessDecision } from "@/components/access-decision";
import { RuntimeInventory } from "@/components/runtime-inventory";
import type { RuntimeNode } from "@/lib/runtime";

const runtimeNodes: RuntimeNode[] = [
  {
    id: "node-berlin",
    name: "Berlin gateway",
    owner: "Platform operations",
    platform: "Linux",
    rawStatus: "ready",
    lastSeenAt: 1_788_171_540,
    endpoint: "198.51.100.24:51820",
    warning: null,
  },
  {
    id: "node-dubai",
    name: "Dubai gateway",
    owner: "Platform operations",
    platform: "Linux",
    rawStatus: "degraded",
    lastSeenAt: 1_788_171_420,
    endpoint: "203.0.113.18:51820",
    warning: "Latest desired state is waiting for runtime confirmation.",
  },
  {
    id: "node-build",
    name: "build-server-01",
    owner: "Engineering",
    platform: "Ubuntu",
    rawStatus: "offline",
    lastSeenAt: 1_788_168_000,
    endpoint: null,
    warning: "No recent handshake was reported.",
  },
];

export default function ClientPage() {
  return (
    <main className="page-shell">
      <header className="page-header">
        <div>
          <p className="product-name">Nanami Client App</p>
          <h1>Runtime overview</h1>
          <p className="page-description">
            Desired configuration and observed device state stay visibly
            separate.
          </p>
        </div>
        <span className="environment-label">Public-safe runtime sample</span>
      </header>

      <AccessDecision
        subject="engineering-team"
        service="admin-console"
        policy="Engineering private services"
        gateway="Berlin gateway"
      />

      <RuntimeInventory nodes={runtimeNodes} />
    </main>
  );
}
