import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import fs from "fs";
const certPath = "../certs/localhost+2.pem";
const keyPath = "../certs/localhost+2-key.pem";
const https =
  fs.existsSync(certPath) && fs.existsSync(keyPath)
    ? { cert: fs.readFileSync(certPath), key: fs.readFileSync(keyPath) }
    : undefined; // CI has no certs, and the build doesn't need them

export default defineConfig({
  server: {
    proxy: {
      "/api/": {
        target: "https://localhost:9788",
        changeOrigin: true,
        secure: false,
      },
      "/ws/": {
        target: "https://localhost:9788",
        changeOrigin: true,
        secure: false, // accept self-signed cert
        ws: true, // proxy the Upgrade request
      },
    },
    https,
  },
  plugins: [
    tailwindcss(),
    sveltekit({
      compilerOptions: {
        // Force runes mode for the project, except for libraries. Can be removed in svelte 6.
        runes: ({ filename }) =>
          filename.split(/[/\\]/).includes("node_modules") ? undefined : true,
      },

      adapter: adapter({
        pages: "../server/dist",
        assets: "../server/dist",
        fallback: "index.html",
      }),
    }),
  ],
});
