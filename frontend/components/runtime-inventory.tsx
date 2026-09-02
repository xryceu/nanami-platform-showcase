"use client";

import { useMemo, useState } from "react";

import {
  filterRuntimeNodes,
  formatLastSeen,
  getRuntimeSummary,
  normalizeConnectionStatus,
  type RuntimeFilter,
  type RuntimeNode,
} from "@/lib/runtime";
import { resolveCurrentTransportStatus } from "@/product/client/transport-contract";

const filters: Array<{ value: RuntimeFilter; label: string }> = [
  { value: "all", label: "All devices" },
  { value: "ready", label: "Ready" },
  { value: "attention", label: "Needs attention" },
];

export function RuntimeInventory({ nodes }: { nodes: RuntimeNode[] }) {
  const [filter, setFilter] = useState<RuntimeFilter>("all");
  const visibleNodes = useMemo(
    () => filterRuntimeNodes(nodes, filter),
    [filter, nodes],
  );
  const summary = getRuntimeSummary(nodes);

  return (
    <section className="panel" aria-labelledby="runtime-inventory-title">
      <div className="section-heading inventory-heading">
        <div>
          <h2 id="runtime-inventory-title">Managed devices</h2>
          <p>
            {summary.ready} ready · {summary.attention} need attention
          </p>
        </div>

        <div className="filter-group" aria-label="Filter managed devices">
          {filters.map((item) => (
            <button
              aria-pressed={filter === item.value}
              className={filter === item.value ? "is-active" : undefined}
              key={item.value}
              onClick={() => setFilter(item.value)}
              type="button"
            >
              {item.label}
            </button>
          ))}
        </div>
      </div>

      {visibleNodes.length === 0 ? (
        <p className="empty-state">No devices match this filter.</p>
      ) : (
        <div
          className="inventory-table"
          role="table"
          aria-label="Runtime inventory"
        >
          <div className="inventory-row inventory-header" role="row">
            <span role="columnheader">Device</span>
            <span role="columnheader">Owner</span>
            <span role="columnheader">Endpoint</span>
            <span role="columnheader">Observed state</span>
            <span role="columnheader">Last seen</span>
          </div>

          {visibleNodes.map((node) => {
            const status = normalizeConnectionStatus(node.rawStatus);
            const transportStatus = resolveCurrentTransportStatus({
              status: node.rawStatus,
              connectionStatus: node.rawStatus,
            });

            return (
              <article className="inventory-row" key={node.id} role="row">
                <div className="device-cell" role="cell">
                  <span
                    aria-hidden="true"
                    className={`status-dot status-dot-${status}`}
                  />
                  <span>
                    <strong>{node.name}</strong>
                    <small>{node.platform}</small>
                  </span>
                </div>
                <div className="mobile-field" role="cell">
                  <span className="mobile-label">Owner</span>
                  <span>{node.owner}</span>
                </div>
                <div className="mobile-field mono" role="cell">
                  <span className="mobile-label">Endpoint</span>
                  <span>{node.endpoint ?? "Not reported"}</span>
                </div>
                <div className="mobile-field" role="cell">
                  <span className="mobile-label">Observed state</span>
                  <span className={`status status-${status}`}>
                    {transportStatus}
                  </span>
                  {node.warning ? (
                    <small className="warning-copy">{node.warning}</small>
                  ) : null}
                </div>
                <div className="mobile-field" role="cell">
                  <span className="mobile-label">Last seen</span>
                  <time
                    dateTime={new Date(node.lastSeenAt * 1000).toISOString()}
                  >
                    {formatLastSeen(node.lastSeenAt, 1_788_171_600)}
                  </time>
                </div>
              </article>
            );
          })}
        </div>
      )}
    </section>
  );
}
