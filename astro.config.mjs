// @ts-check
import { defineConfig } from "astro/config";
import viteBasicSslPlugin from "@vitejs/plugin-basic-ssl";

import react from "@astrojs/react";

// https://astro.build/config
export default defineConfig({
  vite: {
    plugins: [viteBasicSslPlugin()],
    server: {
      https: true,
    },
  },

  integrations: [react()],
});
