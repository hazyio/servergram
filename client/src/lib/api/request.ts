import createClient, { type Middleware } from "openapi-fetch";
import { goto } from "$app/navigation";
import type { paths } from "./schema";

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
