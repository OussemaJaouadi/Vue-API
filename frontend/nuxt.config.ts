import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  css: ['~/assets/css/tailwind.css'],
  compatibilityDate: '2025-01-01',
  // Nuxt runtime config is overridden by matching runtime env vars:
  // apiBaseUrl -> NUXT_API_BASE_URL, public.appName -> NUXT_PUBLIC_APP_NAME.
  // Keep defaults literal so Docker/prod can change env without rebuilding.
  // https://nuxt.com/docs/4.x/guide/going-further/runtime-config
  runtimeConfig: {
    apiBaseUrl: 'http://localhost:8080',
    public: {
      appName: 'Vue API Workbench',
    },
  },
  vite: {
    plugins: [tailwindcss()],
  },
  modules: ['shadcn-nuxt'],
  shadcn: {
    /**
     * Prefix for all the imported component.
     * @default "Ui"
     */
    prefix: 'Ui',
    /**
     * Directory that the component lives in.
     * Will respect the Nuxt aliases.
     * @link https://nuxt.com/docs/api/nuxt-config#alias
     * @default "@/components/ui"
     */
    componentDir: '@/components/ui',
  },
})
