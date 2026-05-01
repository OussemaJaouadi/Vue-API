<script setup lang="ts">
const { appName } = useApiConfig()
const { get } = useApiClient()
const route = useRoute()
const auth = useAuthSession()

const healthStatus = ref<'loading' | 'ok' | 'error'>('loading')
const isAuthRoute = computed(() => route.path === '/login' || route.path === '/register')

onMounted(async () => {
  try {
    const res = await get<{ status: string }>('/healthz')
    if (res.status === 'ok') {
      healthStatus.value = 'ok'
    } else {
      healthStatus.value = 'error'
    }
  } catch (e) {
    healthStatus.value = 'error'
  }
})
</script>

<template>
  <div v-if="isAuthRoute" class="min-h-screen bg-background text-foreground">
    <slot />
  </div>

  <div v-else class="flex h-screen overflow-hidden bg-background">
    <aside class="flex w-64 flex-col border-r bg-sidebar">
      <div class="p-4 border-b">
        <h1 class="font-heading font-bold text-xl">{{ appName }}</h1>
      </div>
      
      <nav class="flex-1 p-4 space-y-2">
        <div class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
          Workbench
        </div>
        <div class="p-2 rounded bg-accent text-accent-foreground cursor-pointer">
          Collections
        </div>
        <div class="p-2 rounded hover:bg-accent hover:text-accent-foreground cursor-pointer text-muted-foreground transition-colors">
          Environments
        </div>
      </nav>

      <div class="space-y-3 border-t p-4 text-xs">
        <div v-if="auth.user.value" class="space-y-1">
          <div class="font-medium text-sidebar-foreground">{{ auth.user.value.username }}</div>
          <div class="text-muted-foreground">{{ auth.user.value.email }}</div>
        </div>

        <div class="flex items-center gap-2">
          <div
            class="h-2 w-2 rounded-full"
            :class="{
              'bg-yellow-400 animate-pulse': healthStatus === 'loading',
              'bg-green-500': healthStatus === 'ok',
              'bg-red-500': healthStatus === 'error'
            }"
          />
          <span class="text-muted-foreground">Backend: {{ healthStatus }}</span>
        </div>

        <button
          v-if="auth.user.value"
          class="w-full rounded-md border px-3 py-2 text-left font-medium text-sidebar-foreground transition-colors hover:bg-sidebar-accent"
          type="button"
          @click="auth.logout"
        >
          Sign out
        </button>
      </div>
    </aside>

    <main class="flex flex-1 flex-col overflow-hidden">
      <header class="h-12 border-b flex items-center px-4 bg-card">
        <div class="text-sm font-medium">
          <span v-if="auth.user.value?.globalRole === 'manager'">Manager console</span>
          <span v-else>Workspace access</span>
        </div>
      </header>

      <div class="flex-1 overflow-auto p-6">
        <slot />
      </div>
    </main>
  </div>
</template>
