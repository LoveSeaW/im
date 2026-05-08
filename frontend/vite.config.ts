import {fileURLToPath, URL} from 'node:url'
import {defineConfig, loadEnv} from 'vite'
import vue from '@vitejs/plugin-vue'
import type {ImportMetaEnv} from "./env";

export default defineConfig(({mode}) => {
  let env: Record<keyof ImportMetaEnv, string> = loadEnv(mode, process.cwd())

  const serverUrl =  env.VITE_SERVER_URL
  return {
    define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
    },
    plugins: [
      vue(),
    ],
    envDir: "./",
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      host: "0.0.0.0",
      port: 80,
      proxy: {
        "/api/chat/ws":{
          target: serverUrl,
          changeOrigin: true,
          ws: true,
        },
        "/api/group/ws":{
          target: serverUrl,
          changeOrigin: true,
          ws: true,
        },
        "/api": {
          target: serverUrl,
          changeOrigin: true,
        }
      }
    }
  }
})
