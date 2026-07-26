import { config } from "../env.js";

let cachedToken = null;
let tokenExpiresAt = 0;

export async function getAccessToken() {
  if (cachedToken && Date.now() < tokenExpiresAt - 60000) {
    return cachedToken;
  }

  const domain = config.AUTH0_DOMAIN;
  const body = new URLSearchParams({
    grant_type: "client_credentials",
    client_id: config.AUTH0_CLIENT_ID,
    client_secret: config.AUTH0_CLIENT_SECRET,
    audience: config.AUTH0_AUDIENCE,
    organization: "org_lf4I9FRW65XClil1",
    // scope: config.AUTH0_SCOPE,
  });

  const res = await fetch(`https://${domain}/oauth/token`, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Auth0 token request failed: ${res.status} ${text}`);
  }

  const data = await res.json();
  cachedToken = data.access_token;
  tokenExpiresAt = Date.now() + data.expires_in * 1000;

  return cachedToken;
}

export async function clearTokenCache() {
  cachedToken = null;
  tokenExpiresAt = 0;
}
