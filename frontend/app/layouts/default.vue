<script setup lang="ts">
const { appName } = useApiConfig()
const { get } = useApiClient()

const healthStatus = ref<'loading' | 'ok' | 'error'>('loading')

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
  <div class="flex h-screen overflow-hidden bg-background">
    <!-- Sidebar -->
    <aside class="w-64 border-r bg-sidebar flex flex-col">
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

      <div class="p-4 border-t text-xs flex items-center gap-2">
        <div 
          class="w-2 h-2 rounded-full" 
          :class="{
            'bg-yellow-400 animate-pulse': healthStatus === 'loading',
            'bg-green-500': healthStatus === 'ok',
            'bg-red-500': healthStatus === 'error'
          }"
        />
        <span class="text-muted-foreground">
          Backend: {{ healthStatus }}
        </span>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- Header / Tabs Placeholder -->
      <header class="h-12 border-b flex items-center px-4 bg-card">
        <div class="text-sm font-medium">New Request</div>
      </header>

      <!-- Content Area -->
      <div class="flex-1 overflow-auto p-6">
        <slot />
      </div>
    </main>
  </div>
</template>
