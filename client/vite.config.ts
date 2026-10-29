import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import fs from "fs";

export default defineConfig({
  server: {
    proxy: {
      "/api/": {
        target: "https://localhost:9788",
        changeOrigin: true,
        secure: false,
      },
    },
    https: {
      cert: fs.readFileSync("../certs/localhost+2.pem"),
      key: fs.readFileSync("../certs/localhost+2-key.pem"),
    },
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
