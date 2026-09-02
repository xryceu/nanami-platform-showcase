import { describe, expect, it } from "vitest";

import { createPricingModel } from "./pricing";

describe("pricing destinations", () => {
  it("keeps standalone destinations root-relative", () => {
    const model = createPricingModel({
      signupUrl: "/client",
      contactUrl: "/",
      docsSelfHostedUrl: "/docs",
    });

    expect(model.tiers.find((tier) => tier.key === "free")?.cta.href).toBe(
      "/client",
    );
    expect(model.tiers.find((tier) => tier.key === "team")?.cta.href).toBe(
      "/client?redirect=%2Fsettings%2Fplan%3FrequestedPlan%3Dteam",
    );
    expect(
      model.tiers.find((tier) => tier.key === "open_source")?.cta.href,
    ).toBe("/docs");
  });

  it("preserves configured HTTPS origins", () => {
    const model = createPricingModel({
      signupUrl: "https://app.example.com/login",
      contactUrl: "https://example.com/contact",
      docsSelfHostedUrl: "https://docs.example.com/self-hosted",
    });

    expect(model.tiers.find((tier) => tier.key === "team")?.cta.href).toBe(
      "https://app.example.com/login?redirect=%2Fsettings%2Fplan%3FrequestedPlan%3Dteam",
    );
  });
});
