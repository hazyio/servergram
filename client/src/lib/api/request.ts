import createClient from "openapi-fetch";
import type { paths } from "./schema";

export function NewRequestClient() {
  return createClient<paths>();
}
