import createClient from "openapi-fetch";
import type { paths } from "./schema";

export const requestClient = createClient<paths>();
