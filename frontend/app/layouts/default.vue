<script setup lang="ts">
const { appName } = useApiConfig()
const { get } = useApiClient()
const auth = useAuthSession()

const healthStatus = ref<'loading' | 'ok' | 'error'>('loading')
const sidebarCollapsed = useState('shell:sidebar-collapsed', () => false)
const mobileSidebarOpen = ref(false)

const healthVariant = computed(() => {
  if (healthStatus.value === 'ok') {
    return 'default'
  }
  if (healthStatus.value === 'error') {
    return 'destructive'
  }
  return 'secondary'
})

const pageTitle = computed(() => {
  if (auth.user.value?.globalRole === 'manager') {
    return 'Manager console'
  }
  return 'Workspace access'
})

const pageSubtitle = computed(() => {
  if (auth.lastEvent.value?.type === 'connected') {
    return 'Connected'
  }
  return auth.lastEvent.value?.type || 'Ready'
})

onMounted(async () => {
  const authenticated = await auth.loadMe()
  if (!authenticated) {
    await navigateTo('/login')
    return
  }

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
  <div class="flex h-dvh overflow-hidden bg-background text-foreground">
    <AppSidebar
      :app-name="appName"
      class="hidden border-r md:flex"
      :collapsed="sidebarCollapsed"
      :health-status="healthStatus"
      :user="auth.user.value"
      @logout="auth.logout"
      @navigate="mobileSidebarOpen = false"
      @toggle="sidebarCollapsed = !sidebarCollapsed"
    />

    <main class="flex flex-1 flex-col overflow-hidden">
      <UiSheet v-model:open="mobileSidebarOpen">
        <AppTopbar
          project-name="Core API"
          workspace-name="Personal Workspace"
          @open-mobile="mobileSidebarOpen = true"
        />
        <UiSheetContent
          :show-close-button="false"
          accessibility-description="Primary workbench navigation and session controls."
          accessibility-title="Navigation"
          class="p-0"
          side="left"
          style="width: 18rem; max-width: 18rem;"
        >
          <AppSidebar
            :app-name="appName"
            class="flex h-full"
            :collapsed="false"
            :health-status="healthStatus"
            mobile
            :user="auth.user.value"
            @logout="auth.logout"
            @toggle="mobileSidebarOpen = false"
            @navigate="mobileSidebarOpen = false"
          />
        </UiSheetContent>
      </UiSheet>
      <UiScrollArea class="flex-1">
        <div class="p-2 md:p-3 lg:p-4">
          <slot />
        </div>
      </UiScrollArea>
    </main>
  </div>
</template>
