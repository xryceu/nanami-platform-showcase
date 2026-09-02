import { describe, expect, it } from "vitest";

import {
  dashboardTransportDetail,
  normalizeLegacyControlCenterConnection,
  normalizeLegacyControlCenterStatus,
  resolveCurrentTransportConnection,
  resolveCurrentTransportStatus,
} from "./transport-contract";

describe("control center transport contract", () => {
  it("normalizes legacy relay filters into the current mediated vocabulary", () => {
    expect(normalizeLegacyControlCenterStatus("relay")).toBe("mediated");
    expect(normalizeLegacyControlCenterConnection("relay")).toBe("mediated");
  });

  it("derives mediated transport only from current gateway-mediated connectivity", () => {
    const node = {
      status: "online",
      connectionStatus: "online",
      connectivity: {
        gatewayId: "gateway-1",
        gatewayNodeId: "gateway-1",
        connectionStatus: "online",
      },
    } as never;

    const statusTone = resolveCurrentTransportStatus(node);
    expect(statusTone).toBe("mediated");
    expect(resolveCurrentTransportConnection(node, statusTone)).toBe(
      "mediated",
    );
  });

  it("describes current board summaries in user-facing language", () => {
    expect(
      dashboardTransportDetail({
        directLinks: 4,
        mediatedLinks: 2,
      }),
    ).toContain("through a gateway");
  });
});
