import {
  createAuth0Client,
  type Auth0Client,
  type User,
} from "@auth0/auth0-spa-js";
import { browser } from "$app/environment";

let client: Auth0Client | null = $state(null);
let userData: User | null = $state(null);
let authenticated = $state(false);
let loading = $state(true);
let authError: string | null = $state(null);

export const auth = {
  get client() {
    return client;
  },
  get user() {
    return userData;
  },
  get isAuthenticated() {
    return authenticated;
  },
  get isLoading() {
    return loading;
  },
  get error() {
    return authError;
  },
};

export async function initializeAuth() {
  if (!browser) return;

  try {
    const c = await createAuth0Client({
      domain: import.meta.env.VITE_AUTH0_DOMAIN,
      clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
      authorizationParams: {
        redirect_uri: window.location.origin,
        audience: "https://dev-eg15l03z5f88rzbd.us.auth0.com/api/v2/",
      },
      useRefreshTokens: true,
      cacheLocation: "localstorage",
    });
    console.log(window.location.origin);

    client = c;

    if (window.location.search.includes("code=")) {
      await c.handleRedirectCallback();
      window.history.replaceState({}, document.title, window.location.pathname);
    }

    const isAuthed = await c.isAuthenticated();
    authenticated = isAuthed;

    if (isAuthed) {
      const u = await c.getUser();
      userData = u || null;
    }

    authError = null;
  } catch (err) {
    console.error("Auth initialization error:", err);
    authError =
      err instanceof Error
        ? err.message
        : "Authentication initialization failed";
  } finally {
    loading = false;
  }
}

export async function login() {
  if (client) {
    await client.loginWithRedirect();
  }
}

export async function logout() {
  if (client) {
    client.logout({
      logoutParams: {
        returnTo: window.location.origin,
      },
    });
  }
}

export async function getToken(): Promise<string | null> {
  if (!client) return null;

  try {
    return await client.getTokenSilently();
  } catch (err: any) {
    if (err.error === "login_required") {
      await login();
    }
    return null;
  }
}
