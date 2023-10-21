/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-21 12:41:04
 * @Description: file content
 */
import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  base: "/",
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      "nsfwjs":"nsfwjs/dist/nsfwjs.min.js",
    },
  },
  server:{
    host: "0.0.0.0",
    port: 18891,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:18890",
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  }
});
