<script setup lang="ts">
import {
  PhCaretDown,
  PhCircleNotch,
  PhPaperPlaneTilt,
  PhShieldCheck,
  PhDotOutline,
  PhArrowRight,
} from '@phosphor-icons/vue'
import {
  API_METHODS,
  type ApiMethod,
  METHOD_COLORS,
  METHOD_LABELS,
} from '~/composables/useWorkbench'

const workbench = useWorkbench()

const setMethod = (method: ApiMethod) => {
  workbench.setActiveRequestMethod(method)
}

const isActiveMethod = (method: ApiMethod) => workbench.activeRequest.value.method === method

const methodRowClass = (method: ApiMethod) => {
  if (isActiveMethod(method)) {
    return 'bg-primary/8'
  }

  const hoverByMethod: Record<ApiMethod, string> = {
    GET: 'hover:bg-emerald-500/10 focus:bg-emerald-500/10',
    POST: 'hover:bg-blue-500/10 focus:bg-blue-500/10',
    PUT: 'hover:bg-amber-500/10 focus:bg-amber-500/10',
    PATCH: 'hover:bg-purple-500/10 focus:bg-purple-500/10',
    DELETE: 'hover:bg-destructive/10 focus:bg-destructive/10',
    SOCKET: 'hover:bg-cyan-500/10 focus:bg-cyan-500/10',
  }

  return hoverByMethod[method]
}

// Combine for local display, but keep reactive
const fullUrl = computed({
  get: () => {
    const baseUrl = workbench.isWebSocketRequest.value ? '{{wsBaseUrl}}' : workbench.requestTarget.value.baseUrl
    return `${baseUrl}${workbench.activeRequest.value.path}`
  },
  set: (val) => {
    // Basic logic to split back for now - in production this would be smarter
    if (val.startsWith('{{wsBaseUrl}}')) {
      workbench.activeRequest.value.path = val.replace('{{wsBaseUrl}}', '')
    } else if (val.startsWith('{{apiBaseUrl}}')) {
      workbench.requestTarget.value.baseUrl = '{{apiBaseUrl}}'
      workbench.activeRequest.value.path = val.replace('{{apiBaseUrl}}', '')
    } else {
      workbench.activeRequest.value.path = val
    }
  }
})
</script>

<template>
  <div class="flex h-12 shrink-0 items-stretch border-b border-border bg-muted/15">
    <!-- Method selector -->
    <UiDropdownMenu>
      <UiDropdownMenuTrigger as-child>
        <button
          class="group relative flex w-30 items-center justify-between border-r border-border bg-background px-3 font-mono text-[11px] font-black uppercase tracking-widest transition-colors hover:bg-muted/45 active:bg-muted/70"
          type="button"
        >
          <div class="absolute inset-y-2 left-0 w-1 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.5)]" />
          <div class="flex items-center gap-1.5">
            <span :class="METHOD_COLORS[workbench.activeRequest.value.method]">
              {{ METHOD_LABELS[workbench.activeRequest.value.method] }}
            </span>
          </div>
          <PhCaretDown class="size-3 opacity-50 transition-opacity group-hover:opacity-90" />
        </button>
      </UiDropdownMenuTrigger>
      <UiDropdownMenuContent align="start" class="w-40 rounded-none border-2 border-border p-1.5 shadow-[4px_4px_0_0_rgba(0,0,0,0.08)]">
        <UiDropdownMenuItem
          v-for="method in API_METHODS"
          :key="method"
          class="group relative rounded-none p-0 transition-colors"
          :class="methodRowClass(method)"
          @click="setMethod(method)"
        >
          <div
            v-if="isActiveMethod(method)"
            class="absolute left-0 top-[20%] h-[60%] w-0.75 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.5)]"
          />
          <div class="flex w-full items-center justify-between px-2.5 py-2 font-mono text-[11px] font-black uppercase tracking-widest">
            <span :class="METHOD_COLORS[method]">{{ METHOD_LABELS[method] }}</span>
            <span v-if="isActiveMethod(method)" class="h-px w-5 bg-primary" />
          </div>
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>

    <!-- URL input -->
    <div class="group relative flex min-w-0 flex-1 items-center border-r border-border bg-background transition-colors focus-within:bg-card">
      <div class="pointer-events-none absolute inset-x-0 bottom-0 h-px bg-primary/60 opacity-0 transition-opacity group-focus-within:opacity-100" />
      <div class="flex min-w-0 flex-1 items-center gap-2 px-3 font-mono text-xs">
        <span class="hidden shrink-0 border border-primary/15 bg-primary/5 px-1.5 py-0.5 text-[8px] font-black uppercase tracking-widest text-primary/70 sm:inline">
          URL
        </span>
        <input
          v-model="fullUrl"
          aria-label="Request URL"
          class="h-full w-full min-w-0 bg-transparent font-semibold tracking-tight text-foreground outline-none placeholder:text-muted-foreground/45"
          placeholder="Enter request URL or path..."
          spellcheck="false"
        >
      </div>

      <!-- Proxy status -->
      <div class="ml-2 flex shrink-0 items-center gap-1.5 border-l border-border pl-3">
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <div class="mr-3 flex cursor-help items-center gap-1 text-primary">
              <PhDotOutline class="size-2.5" />
              <PhShieldCheck class="size-3" />
              <span class="font-mono text-[8px] font-bold uppercase tracking-widest">Proxy</span>
            </div>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">Private backend proxy enabled</UiTooltipContent>
        </UiTooltip>
      </div>
    </div>

    <!-- Execute action -->
    <button
      class="group relative flex min-w-36 items-center justify-center overflow-hidden border-l border-primary bg-primary text-primary-foreground transition-all hover:bg-primary/95 active:translate-y-px"
      type="button"
      @click="workbench.executeActiveRequest"
    >
      <div class="absolute inset-y-0 left-0 w-1 bg-primary-foreground/35" />
      <span class="relative flex h-full items-center gap-2 px-5 font-mono text-[10px] font-black uppercase tracking-[0.22em]">
        <PhCircleNotch v-if="workbench.loading.value" class="size-3.5 animate-spin" />
        <PhPaperPlaneTilt v-else class="size-3.5 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
        {{ workbench.isWebSocketRequest.value ? (workbench.webSocketState.value === 'connected' ? 'Disconnect' : 'Connect') : 'Execute' }}
        <PhArrowRight class="size-3 opacity-55 transition-transform group-hover:translate-x-1" />
      </span>
    </button>
  </div>
</template>
