import type { paths } from "#lib/types/request-schema.js";
import { goto } from "$app/navigation";
import createClient, { type Middleware } from "openapi-fetch";

const authMiddleware: Middleware = {
  async onResponse({ response }) {
    console.log(response);
    // A wrong password also returns 401, so don't redirect from the login page.
    if (
      response.status === 401 &&
      !location.pathname.startsWith("/account/login")
    ) {
      const current = location.pathname + location.search;
      await goto(`/account/login?redirect=${encodeURIComponent(current)}`);
    }
    return response; // always return it (or undefined) so the chain continues
  },
};

export function NewRequestClient() {
  const client = createClient<paths>();
  client.use(authMiddleware);
  return client;
}

export function unWrapErrorMessage(
  error: {
    detail?: string | undefined;
    errors?:
      | {
          location?: string | undefined;
          message?: string | undefined;
          value?: unknown;
        }[]
      | null
      | undefined;
    instance?: string | undefined;
    status?: number | undefined;
    title?: string | undefined;
    type: string;
  },
  or: string,
) {
  return error.detail ?? error.title ?? or;
}
