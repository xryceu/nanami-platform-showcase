import { describe, expect, it } from "vitest";

import {
  isCommunityOnlyClientPath,
  isCompanyOwnedPath,
  isDeletedLegacyClientPath,
  isSaaSOnlyClientPath,
  sanitizeClientPostAuthRedirect,
} from "./client-boundary";

describe("Client App boundary route contract", () => {
  it("classifies product-owned paths", () => {
    expect(isSaaSOnlyClientPath("/settings/plan")).toBe(true);
    expect(isCommunityOnlyClientPath("/settings/runtime/gateways")).toBe(true);
    expect(isCompanyOwnedPath("/problem-center/critical")).toBe(true);
    expect(isDeletedLegacyClientPath("/topology/all")).toBe(true);
  });

  it("sanitizes stale post-auth redirects by app mode", () => {
    expect(
      sanitizeClientPostAuthRedirect(
        "/settings/runtime/gateways?tab=runtime",
        "saas",
      ),
    ).toBe("/");
    expect(sanitizeClientPostAuthRedirect("/settings/plan", "community")).toBe(
      "/",
    );
    expect(
      sanitizeClientPostAuthRedirect(
        "/login/device?user_code=ABCD-2345",
        "saas",
      ),
    ).toBe("/login/device?user_code=ABCD-2345");
    expect(sanitizeClientPostAuthRedirect("/paths?tab=rules", "saas")).toBe(
      "/paths?tab=rules",
    );
  });
});
