import { getAccessToken } from "./auth.js";
import { config } from "../env.js";

export async function fetchFlags(key) {
  const token = await getAccessToken();

  let url = `${config.API_BASE_URL}/api/${key}`;

  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: "application/json",
    },
  });

  if (res.status === 304) {
    return { notModified: true, status: 304 };
  }

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API request failed: ${res.status} ${text}`);
  }

  return { notModified: false, status: res.status, data: await res.json() };
}
