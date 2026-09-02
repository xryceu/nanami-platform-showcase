type PublicDestinationName =
  | "SHOWCASE_SIGNUP_URL"
  | "SHOWCASE_CONTACT_URL"
  | "SHOWCASE_SELF_HOSTED_DOCS_URL";

function publicDestination(
  name: PublicDestinationName,
  fallback: string,
): string {
  const value = process.env[name]?.trim() || fallback;

  if (
    value.startsWith("/") &&
    !value.startsWith("//") &&
    !value.includes("\\")
  ) {
    return value;
  }

  const url = new URL(value);
  if (
    url.protocol !== "https:" ||
    url.username ||
    url.password ||
    !url.hostname
  ) {
    throw new Error(
      `${name} must be a root-relative path or an HTTPS URL without credentials`,
    );
  }

  return url.toString();
}

export function getMarketingDestinations() {
  return {
    signupUrl: publicDestination("SHOWCASE_SIGNUP_URL", "/client"),
    contactUrl: publicDestination("SHOWCASE_CONTACT_URL", "/"),
    docsSelfHostedUrl: publicDestination(
      "SHOWCASE_SELF_HOSTED_DOCS_URL",
      "/docs",
    ),
  };
}
