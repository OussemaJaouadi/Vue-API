<script setup lang="ts">
import {
  PhActivity as Activity,
  PhCaretDown as CaretDown,
  PhDatabase as Database,
  PhFolderOpen as FolderOpen,
  PhHouse as House,
  PhSignOut as SignOut,
  PhUserCircle as UserCircle,
} from '@phosphor-icons/vue'

const { appName } = useApiConfig()
const { get } = useApiClient()
const route = useRoute()
const auth = useAuthSession()

const healthStatus = ref<'loading' | 'ok' | 'error'>('loading')
const isAuthRoute = computed(() => route.path === '/login' || route.path === '/register')
const navItems = [
  { label: 'Overview', icon: House, active: true },
  { label: 'Collections', icon: FolderOpen, active: false },
  { label: 'Environments', icon: Database, active: false },
]

const healthVariant = computed(() => {
  if (healthStatus.value === 'ok') {
    return 'default'
  }
  if (healthStatus.value === 'error') {
    return 'destructive'
  }
  return 'secondary'
})

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

  <div v-else class="flex h-dvh overflow-hidden bg-background text-foreground">
    <aside class="hidden w-64 shrink-0 flex-col border-r bg-sidebar md:flex">
      <div class="flex h-14 items-center gap-2 border-b px-4">
        <AppLogo class="size-9 shrink-0" />
        <div class="min-w-0">
          <h1 class="truncate font-heading text-sm font-semibold">{{ appName }}</h1>
          <p class="text-xs text-muted-foreground">API workbench</p>
        </div>
      </div>

      <nav class="flex-1 space-y-1 p-3">
        <div class="px-2 pb-2 text-xs font-medium uppercase text-muted-foreground">
          Workbench
        </div>
        <UiButton
          v-for="item in navItems"
          :key="item.label"
          class="w-full justify-start"
          :variant="item.active ? 'secondary' : 'ghost'"
        >
          <component :is="item.icon" class="size-4" />
          {{ item.label }}
        </UiButton>
      </nav>

      <div class="space-y-3 border-t p-4 text-xs">
        <div class="flex items-center justify-between">
          <span class="text-muted-foreground">Backend</span>
          <UiBadge :variant="healthVariant">
            <Activity class="size-3" />
            {{ healthStatus }}
          </UiBadge>
        </div>

        <UiDropdownMenu v-if="auth.user.value">
          <UiDropdownMenuTrigger as-child>
            <UiButton class="h-auto w-full justify-between px-0 py-0" variant="ghost">
              <span class="flex min-w-0 items-center gap-2 px-2 py-1.5">
                <UiAvatar class="size-7">
                  <UiAvatarFallback>
                    <UserCircle class="size-4" />
                  </UiAvatarFallback>
                </UiAvatar>
                <span class="min-w-0 text-left">
                  <span class="block truncate font-medium">{{ auth.user.value.username }}</span>
                  <span class="block truncate text-muted-foreground">{{ auth.user.value.email }}</span>
                </span>
              </span>
              <CaretDown class="mr-2 size-3 text-muted-foreground" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end" class="w-56">
            <UiDropdownMenuLabel>{{ auth.user.value.globalRole }}</UiDropdownMenuLabel>
            <UiDropdownMenuSeparator />
            <UiDropdownMenuItem @click="auth.logout">
              <SignOut class="size-4" />
              Sign out
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>
    </aside>

    <main class="flex flex-1 flex-col overflow-hidden">
      <header class="flex h-14 items-center justify-between border-b bg-card px-4">
        <div class="min-w-0">
          <div class="truncate text-sm font-medium">
            <span v-if="auth.user.value?.globalRole === 'manager'">Manager console</span>
            <span v-else>Workspace access</span>
          </div>
          <div class="text-xs text-muted-foreground">
            {{ auth.lastEvent.value?.type || 'Realtime idle' }}
          </div>
        </div>
        <UiBadge variant="outline" class="md:hidden">
          {{ healthStatus }}
        </UiBadge>
      </header>

      <UiScrollArea class="flex-1">
        <div class="p-4 md:p-6">
          <slot />
        </div>
      </UiScrollArea>
    </main>
  </div>
</template>
