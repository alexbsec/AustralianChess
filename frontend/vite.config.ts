import { defineConfig } from "vite";

export default defineConfig({
    server: {
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
            },
            "/api/v1/ws": {
                target: "ws://localhost:8080",
                ws: true,
                changeOrigin: true,
            },
        },
    },
});
