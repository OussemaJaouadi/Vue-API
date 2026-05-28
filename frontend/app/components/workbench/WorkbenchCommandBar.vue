<script setup lang="ts">
import {
  PhCaretDown,
  PhCircleNotch,
  PhPaperPlaneTilt,
  PhShieldCheck,
  PhFloppyDisk,
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
    return 'bg-primary/10 text-primary border-primary'
  }

  const colorByMethod: Record<ApiMethod, string> = {
    GET: 'hover:bg-emerald-500/10 hover:text-emerald-600 dark:hover:text-emerald-400 focus:bg-emerald-500/10',
    POST: 'hover:bg-blue-500/10 hover:text-blue-600 dark:hover:text-blue-400 focus:bg-blue-500/10',
    PUT: 'hover:bg-amber-500/10 hover:text-amber-600 dark:hover:text-amber-400 focus:bg-amber-500/10',
    PATCH: 'hover:bg-purple-500/10 hover:text-purple-600 dark:hover:text-purple-400 focus:bg-purple-500/10',
    DELETE: 'hover:bg-destructive/10 hover:text-destructive focus:bg-destructive/10',
    SOCKET: 'hover:bg-cyan-500/10 hover:text-cyan-600 dark:hover:text-cyan-400 focus:bg-cyan-500/10',
  }

  return colorByMethod[method] + ' border-transparent'
}

const fullUrl = computed({
  get: () => {
    const baseUrl = workbench.isWebSocketRequest.value ? '{{wsBaseUrl}}' : workbench.requestTarget.value.baseUrl
    return `${baseUrl}${workbench.activeRequest.value.path}`
  },
  set: (val) => {
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
  <div class="flex h-14 shrink-0 items-stretch border-b bg-card/30 select-none overflow-hidden">
    <!-- Method selector -->
    <UiDropdownMenu>
      <UiDropdownMenuTrigger as-child>
        <button
          class="group relative flex w-32 items-center justify-between border-r bg-background px-4 font-mono text-[11px] font-black uppercase tracking-widest transition-all hover:bg-muted/40 outline-none"
          type="button"
        >
          <div class="wb-active-indicator" />
          <span :class="METHOD_COLORS[workbench.activeRequest.value.method]">
            {{ METHOD_LABELS[workbench.activeRequest.value.method] }}
          </span>
          <PhCaretDown class="size-3 opacity-70 group-hover:opacity-100" />
        </button>
      </UiDropdownMenuTrigger>
      <UiDropdownMenuContent align="start" class="w-48 rounded-none border-2 border-primary/20 bg-background p-1 shadow-[8px_8px_0_0_rgba(16,185,129,0.12)]">
        <UiDropdownMenuItem
          v-for="method in API_METHODS"
          :key="method"
          class="font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2.5 rounded-none border-l-2 transition-all mb-0.5 last:mb-0"
          :class="methodRowClass(method)"
          @click="setMethod(method)"
        >
          <div class="flex w-full items-center justify-between">
            <span>{{ METHOD_LABELS[method] }}</span>
            <div v-if="isActiveMethod(method)" class="size-1.5 rounded-full bg-primary shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
          </div>
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>

    <!-- URL Input Context -->
    <div class="group relative flex min-w-0 flex-1 items-center bg-background focus-within:bg-card transition-all">
      <div class="flex min-w-0 flex-1 items-center gap-4 px-4 h-full">
        <div class="flex h-6 items-center gap-2 border-2 border-primary/15 bg-primary/5 px-2 font-mono text-[9px] font-black uppercase tracking-widest text-primary/70 shrink-0">
          <PhShieldCheck class="size-3" />
          Gateway
        </div>
        <input
          v-model="fullUrl"
          class="h-full w-full min-w-0 bg-transparent font-mono text-[11px] font-bold tracking-tight text-foreground outline-none placeholder:text-muted-foreground/50"
          placeholder="Path or qualified endpoint..."
          spellcheck="false"
        >
      </div>

      <div class="h-6 w-px bg-border/20" />

      <!-- Action Context -->
      <div class="flex items-center px-4 gap-3">
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <button
              class="group flex size-8 items-center justify-center border-2 border-primary/10 bg-background text-muted-foreground transition-all hover:border-primary/40 hover:bg-primary/5 hover:text-primary active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
              type="button"
              @click="workbench.saveActiveRequestState"
            >
              <PhFloppyDisk class="size-4" />
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">Persist State</UiTooltipContent>
        </UiTooltip>
      </div>
    </div>

    <!-- Terminal Execute -->
    <button
      class="group relative flex min-w-[160px] items-center justify-center border-l-2 border-primary/20 bg-primary/5 text-primary transition-all hover:bg-primary/10 hover:shadow-[inset_4px_0_0_0_theme(colors.primary.DEFAULT)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none outline-none overflow-hidden"
      type="button"
      @click="workbench.executeActiveRequest"
    >
      <div class="flex items-center gap-3 font-mono text-[11px] font-black uppercase tracking-widest">
        <PhCircleNotch v-if="workbench.loading.value" class="size-4 animate-spin" />
        <PhPaperPlaneTilt v-else class="size-4 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
        {{ workbench.isWebSocketRequest.value ? (workbench.webSocketState.value === 'connected' ? 'Kill' : 'Init') : 'Send' }}
      </div>
      
      <!-- Tactile Decor -->
      <div class="absolute right-2 top-1 size-1 bg-primary/20" />
      <div class="absolute right-1 bottom-2 size-1 bg-primary/10" />
    </button>
  </div>
</template>
