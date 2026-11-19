import { self } from "$app/service-worker";
import { assets, immutable, prerendered } from "$app/manifest";

const CACHE_PREFIX = "servergram-";

/**
 * Derive the cache name from the build's immutable assets so that a new deploy
 * installs into a fresh cache and the `activate` handler can safely drop the
 * previous one. A fixed name would keep serving assets from older builds.
 */
function buildVersion() {
  const fingerprint = immutable
    .map(({ path }) => path)
    .sort()
    .join("|");

  let hash = 2166136261;
  for (let index = 0; index < fingerprint.length; index += 1) {
    hash ^= fingerprint.charCodeAt(index);
    hash = Math.imul(hash, 16777619);
  }

  return (hash >>> 0).toString(36);
}

const CACHE_NAME = `${CACHE_PREFIX}${buildVersion()}`;

// The service worker lives at `<base>/service-worker.js`, so this resolves the
// base path without importing `$app/paths`, which is not available here.
const BASE = self.location.pathname.replace(/\/service-worker\.js$/, "");
const OFFLINE_URL = `${BASE}/`;
const OFFLINE_FALLBACK_URL = `${BASE}/index.html`;

function url(path: string) {
  return `${BASE}/${path.replace(/^\//, "")}`;
}

const precache = [
  ...immutable.map(({ path }) => url(path)),
  ...assets.map(({ path }) => url(path)),
  ...prerendered.map(({ path }) => url(path)),
  OFFLINE_URL,
  OFFLINE_FALLBACK_URL,
];

const precacheUrls = [...new Set(precache)];

async function cacheUrls() {
  const cache = await caches.open(CACHE_NAME);

  await Promise.allSettled(
    precacheUrls.map(async (url) => {
      try {
        // `reload` bypasses the browser's HTTP cache so a new worker always
        // installs the assets that belong to its current build.
        const response = await fetch(new Request(url, { cache: "reload" }));
        if (response.ok) {
          await cache.put(url, response);
        }
      } catch {
        // A single missing asset should not prevent the worker installing.
      }
    }),
  );
}

self.addEventListener("install", (event) => {
  event.waitUntil(cacheUrls());
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      const keys = await caches.keys();
      await Promise.all(
        keys
          .filter((key) => key.startsWith(CACHE_PREFIX) && key !== CACHE_NAME)
          .map((key) => caches.delete(key)),
      );

      await self.clients.claim();
    })(),
  );
});

async function handleNavigation(request: Request) {
  try {
    return await fetch(request);
  } catch {
    const cache = await caches.open(CACHE_NAME);
    return (
      (await cache.match(OFFLINE_URL)) ??
      (await cache.match(OFFLINE_FALLBACK_URL)) ??
      Response.error()
    );
  }
}

async function handleAsset(request: Request) {
  const cache = await caches.open(CACHE_NAME);
  const cached = await cache.match(request);

  if (cached) {
    return cached;
  }

  try {
    const response = await fetch(request);
    if (response.ok) {
      await cache.put(request, response.clone());
    }
    return response;
  } catch {
    return Response.error();
  }
}

self.addEventListener("fetch", (event) => {
  const { request } = event;

  if (request.method !== "GET") return;

  const requestUrl = new URL(request.url);
  if (requestUrl.origin !== self.location.origin) return;

  // API calls and websocket upgrades must always reach the server. When the
  // network is unavailable, callers surface the error through ErrorSection.
  if (
    requestUrl.pathname.startsWith(`${BASE}/api/`) ||
    requestUrl.pathname.startsWith(`${BASE}/ws/`)
  ) {
    return;
  }

  // Never cache the worker script itself; the browser must always be able to
  // fetch the newest build's worker.
  if (requestUrl.pathname === `${BASE}/service-worker.js`) return;

  if (request.mode === "navigate") {
    event.respondWith(handleNavigation(request));
    return;
  }

  // Data payloads are not safe to serve stale, so leave them to the network.
  if (requestUrl.pathname.endsWith("/__data.json")) return;
  if (
    requestUrl.pathname.startsWith("/@") ||
    requestUrl.pathname.startsWith("/node_modules/") ||
    requestUrl.pathname.startsWith("/src/") ||
    requestUrl.searchParams.has("v") ||
    requestUrl.searchParams.has("t")
  )
    return;
  event.respondWith(handleAsset(request));
});
