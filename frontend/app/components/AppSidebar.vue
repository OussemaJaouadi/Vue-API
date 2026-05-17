<script setup lang="ts">
import {
  PhDatabase,
  PhFolderOpen,
  PhHouse,
  PhSignOut,
  PhSidebarSimple,
  PhX,
  PhUserCircle,
  PhActivity,
  PhDesktop,
  PhMoon,
  PhSun,
  PhCheck,
} from '@phosphor-icons/vue'
import type { Component } from 'vue'
import type { ThemePreference } from '~/composables/useThemePreference'

const props = defineProps<{
  appName: string
  collapsed?: boolean
  mobile?: boolean
  healthStatus?: 'loading' | 'ok' | 'error'
  user?: {
    username: string
    email: string
    globalRole: string
  } | null
}>()

defineEmits<{
  toggle: []
  logout: []
}>()

const theme = useThemePreference()

const healthLabel = computed(() => {
  if (props.healthStatus === 'ok') return 'Online'
  if (props.healthStatus === 'error') return 'Offline'
  return 'Check'
})

const navItems = [
  { label: 'Overview', icon: PhHouse, active: true },
  { label: 'Collections', icon: PhFolderOpen, active: false },
  { label: 'Environments', icon: PhDatabase, active: false },
]

const themeItems: Array<{
  value: ThemePreference
  label: string
  icon: Component
}> = [
  { value: 'light', label: 'Light', icon: PhSun },
  { value: 'system', label: 'System', icon: PhDesktop },
  { value: 'dark', label: 'Dark', icon: PhMoon },
]

const activeThemeItem = computed(() => themeItems.find(item => item.value === theme.preference.value) ?? themeItems[1])
</script>

<template>
   <aside
     class="flex min-h-0 shrink-0 flex-col bg-muted/10 border-r text-foreground transition-[width] duration-150 ease-out overflow-hidden"
     :class="collapsed ? 'w-16' : 'w-72 md:w-64'"
   >
    <!-- Header -->
    <div class="flex h-14 items-center gap-3 border-b px-4" :class="collapsed && 'justify-center px-0'">
      <UiTooltip v-if="collapsed">
        <UiTooltipTrigger as-child>
          <UiButton class="hidden md:inline-flex rounded-none" size="icon-sm" variant="ghost" @click="$emit('toggle')">
            <PhSidebarSimple class="size-4" />
          </UiButton>
        </UiTooltipTrigger>
        <UiTooltipContent side="right">Expand sidebar</UiTooltipContent>
      </UiTooltip>
      <template v-else>
        <AppLogo class="size-8 shrink-0 text-primary" />
        <div class="min-w-0 flex flex-col justify-center">
          <h1 class="truncate font-heading text-sm font-bold leading-tight">{{ appName }}</h1>
          <p class="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">API workbench</p>
        </div>
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <UiButton
              class="ml-auto inline-flex rounded-none text-muted-foreground hover:text-foreground"
              size="icon-xs"
              variant="ghost"
              @click="$emit('toggle')"
            >
              <PhX v-if="mobile" class="size-4" />
              <PhSidebarSimple v-else class="hidden size-4 md:block" />
            </UiButton>
          </UiTooltipTrigger>
          <UiTooltipContent side="right">{{ mobile ? 'Close navigation' : 'Collapse sidebar' }}</UiTooltipContent>
        </UiTooltip>
      </template>
    </div>

    <!-- Navigation -->
    <nav class="min-h-0 flex-1 space-y-1 p-3">
      <div class="mb-4 flex items-center gap-2 px-2" :class="collapsed && 'sr-only'">
        <div class="h-px flex-1 bg-primary/10" />
        <span class="font-mono text-[9px] font-black uppercase tracking-[0.3em] text-primary/40">Workbench</span>
        <div class="h-px w-2 bg-primary/10" />
      </div>

      <template v-for="(item, index) in navItems" :key="item.label">
        <UiTooltip :disabled="!collapsed">
          <UiTooltipTrigger as-child>
            <UiButton
              class="group relative h-10 w-full rounded-none border border-transparent transition-all duration-200"
              :class="[
                collapsed ? 'justify-center px-0' : 'justify-start px-3',
                item.active 
                  ? 'bg-primary/7 text-primary font-black border-primary/10 shadow-xs' 
                  : 'text-muted-foreground hover:bg-primary/3 hover:text-foreground hover:border-primary/5'
              ]"
              variant="ghost"
            >
              <component 
                :is="item.icon" 
                class="size-4.5 shrink-0 transition-transform group-hover:scale-110" 
                :class="[!collapsed && 'mr-3', item.active ? 'text-primary' : 'text-muted-foreground/40']" 
              />
              <span v-if="!collapsed" class="truncate text-[10px] font-bold uppercase tracking-widest">{{ item.label }}</span>

              <!-- Tactical Index -->
              <span v-if="!collapsed" class="ml-auto font-mono text-[8px] font-black opacity-20 group-hover:opacity-60 transition-opacity">0{{ index + 1 }}</span>

              <!-- Active State Glow Bar -->
              <div v-if="item.active" class="absolute left-0 top-[20%] h-[60%] w-0.75 bg-primary shadow-[0_0_12px_rgba(16,185,129,0.55)]" />
            </UiButton>
          </UiTooltipTrigger>
          <UiTooltipContent v-if="collapsed" side="right">{{ item.label }}</UiTooltipContent>
        </UiTooltip>
      </template>
    </nav>

    <!-- Modular Footer -->
    <div class="mt-auto flex flex-col border-t-2 border-primary/10 bg-primary/2 p-2 gap-2">
      <!-- Module 0: Theme Preference -->
      <UiDropdownMenu>
        <UiDropdownMenuTrigger as-child>
          <button
            class="flex h-9 w-full items-center rounded-none border border-primary/10 bg-background/50 px-3 text-muted-foreground shadow-inner transition-colors hover:border-primary/30 hover:text-foreground"
            :class="collapsed ? 'justify-center px-0' : 'justify-start'"
            type="button"
            :title="`Theme: ${activeThemeItem.label}`"
          >
            <component :is="activeThemeItem.icon" class="size-4 shrink-0 text-primary" :class="!collapsed && 'mr-3'" />
            <span v-if="!collapsed" class="truncate font-mono text-[9px] font-black uppercase tracking-[0.2em]">
              Theme
            </span>
            <span v-if="!collapsed" class="ml-auto font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/55">
              {{ activeThemeItem.label }}
            </span>
          </button>
        </UiDropdownMenuTrigger>

        <UiDropdownMenuContent
          align="start"
          side="right"
          class="w-44 rounded-none border-2 p-1"
        >
          <UiDropdownMenuItem
            v-for="item in themeItems"
            :key="item.value"
            class="grid grid-cols-[18px_minmax(0,1fr)_14px] gap-2 rounded-none px-2 py-2"
            :class="theme.preference.value === item.value ? 'bg-primary/8 text-foreground' : 'text-muted-foreground'"
            @click="theme.setTheme(item.value)"
          >
            <component
              :is="item.icon"
              class="size-3.5 self-center"
              :class="theme.preference.value === item.value ? 'text-primary' : 'text-muted-foreground'"
            />
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">{{ item.label }}</span>
            <PhCheck v-if="theme.preference.value === item.value" class="size-3 self-center text-primary" />
          </UiDropdownMenuItem>
        </UiDropdownMenuContent>
      </UiDropdownMenu>
      
      <!-- Module 1: System Status -->
      <UiTooltip>
        <UiTooltipTrigger as-child>
          <div class="rounded-none border border-primary/10 bg-background/50 p-2 shadow-inner transition-colors hover:border-primary/30 cursor-help">
            <div class="flex items-center gap-3" :class="collapsed ? 'justify-center' : ''">
              <div class="relative flex size-2 shrink-0">
                <span class="absolute inline-flex h-full w-full animate-ping rounded-full opacity-75" :class="healthStatus === 'ok' ? 'bg-emerald-400' : 'bg-destructive'" />
                <span class="relative inline-flex size-2 rounded-full" :class="healthStatus === 'ok' ? 'bg-emerald-500' : 'bg-destructive'" />
              </div>
              <div v-if="!collapsed" class="flex flex-1 items-center justify-between min-w-0">
                <span class="font-mono text-[9px] font-bold uppercase tracking-tighter text-muted-foreground">Engine Status</span>
                <span class="font-mono text-[9px] font-black uppercase text-emerald-500" :class="healthStatus !== 'ok' && 'text-destructive'">{{ healthLabel }}</span>
              </div>
            </div>
          </div>
        </UiTooltipTrigger>
        <UiTooltipContent side="right">System: {{ healthLabel }}</UiTooltipContent>
      </UiTooltip>

      <!-- Module 2: Identity Card -->
      <div v-if="user" class="rounded-none border border-primary/10 bg-background/50 p-2 shadow-inner">
        <div class="flex items-center gap-3" :class="collapsed ? 'justify-center' : ''">
          <UiTooltip :disabled="!collapsed">
            <UiTooltipTrigger as-child>
              <UiAvatar class="size-8 rounded-none border-2 border-primary/20 shrink-0 transition-all hover:border-primary/50">
                <UiAvatarFallback class="rounded-none bg-primary/5 text-primary">
                  <PhUserCircle class="size-5" />
                </UiAvatarFallback>
              </UiAvatar>
            </UiTooltipTrigger>
            <UiTooltipContent side="right" class="p-0 border-0 shadow-none bg-transparent">
              <div class="w-64 rounded-none border-2 border-primary/40 bg-background/95 backdrop-blur-xl shadow-[6px_6px_0_0_rgba(16,185,129,0.2)]">
                <div class="flex flex-col gap-1.5 p-4 bg-primary/5">
                  <div class="flex items-center justify-between">
                    <span class="text-[9px] font-black uppercase tracking-[0.2em] text-primary">Identity Card</span>
                    <div class="flex items-center gap-1.5">
                      <div class="size-1.5 rounded-full" :class="healthStatus === 'ok' ? 'bg-emerald-500' : 'bg-destructive'" />
                      <span class="font-mono text-[8px] font-bold text-muted-foreground uppercase">{{ healthLabel }}</span>
                    </div>
                  </div>
                  <div class="mt-2 space-y-0.5">
                    <span class="block truncate text-sm font-bold text-foreground">{{ user.username }}</span>
                    <span class="block truncate font-mono text-[10px] text-muted-foreground">{{ user.email }}</span>
                    <div class="pt-1">
                      <span class="text-[8px] font-black px-1.5 py-0.5 bg-primary/10 text-primary uppercase rounded-[2px]">{{ user.globalRole }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </UiTooltipContent>
          </UiTooltip>
          
          <div v-if="!collapsed" class="min-w-0 flex-1 flex flex-col justify-center">
            <div class="flex items-center justify-between">
              <span class="truncate text-[10px] font-black uppercase tracking-tighter text-foreground/90">{{ user.username }}</span>
              <span class="text-[8px] font-bold px-1 bg-primary/10 text-primary uppercase rounded-[2px]">{{ user.globalRole }}</span>
            </div>
            <span class="truncate font-mono text-[8px] uppercase tracking-widest text-muted-foreground/50">{{ user.email }}</span>
          </div>
        </div>
      </div>

      <!-- Module 3: Session Control -->
      <UiTooltip>
        <UiTooltipTrigger as-child>
          <UiButton
            class="w-full rounded-none border border-destructive/20 bg-destructive/5 shadow-inner h-9 text-destructive transition-all hover:text-white hover:bg-destructive hover:border-destructive hover:shadow-[0_0_10px_rgba(220,38,38,0.3)]"
            :class="collapsed ? 'justify-center p-0' : 'justify-start px-3'"
            variant="ghost"
            @click="$emit('logout')"
          >
            <PhSignOut class="size-4 shrink-0" :class="!collapsed && 'mr-3'" />
            <span v-if="!collapsed" class="truncate text-[9px] font-black uppercase tracking-[0.2em]">Terminate Session</span>
          </UiButton>
        </UiTooltipTrigger>
        <UiTooltipContent side="right">Terminate Session</UiTooltipContent>
      </UiTooltip>

    </div>
  </aside>
</template>
