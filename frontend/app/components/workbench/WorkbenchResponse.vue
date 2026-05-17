<script setup lang="ts">
import {
  PhBracketsCurly,
  PhClock,
  PhCode,
  PhDatabase,
  PhFingerprint,
  PhShieldCheck,
  PhListDashes,
  PhPlug,
} from '@phosphor-icons/vue'
import type { BodyLanguage } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const activeView = ref<'body' | 'headers'>('body')
const responseLanguage = ref<BodyLanguage>('json')

watch(
  () => workbench.responseData.value?.decoded,
  (decoded) => {
    responseLanguage.value = decoded === 'json' ? 'json' : 'text'
  },
  { immediate: true },
)
</script>

<template>
  <div v-if="workbench.isWebSocketRequest.value" class="flex h-full flex-col overflow-hidden">
    <div class="flex h-9 shrink-0 items-center justify-between border-b bg-muted/30 px-3">
      <div class="flex items-center gap-2">
        <UiBadge
          variant="outline"
          class="h-5 rounded-none font-mono text-[9px] font-black uppercase"
          :class="workbench.webSocketState.value === 'connected' ? 'border-primary/30 bg-primary/10 text-primary' : 'border-muted-foreground/20 text-muted-foreground'"
        >
          {{ workbench.webSocketState.value }}
        </UiBadge>
        <div class="flex items-center gap-1.5 font-mono text-[9px] font-bold uppercase text-muted-foreground/60">
          <PhPlug class="size-3" />
          Backend-owned socket stream
        </div>
      </div>
      <div class="font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/40">
        {{ workbench.webSocketEvents.value.length }} events
      </div>
    </div>

    <div class="min-h-0 flex-1 overflow-auto bg-card/30 p-3">
      <div class="grid gap-2">
        <button
          v-for="event in workbench.webSocketEvents.value"
          :key="event.id"
          class="grid grid-cols-[64px_72px_minmax(0,1fr)_64px] items-start border bg-background text-left font-mono text-[10px] transition-colors hover:border-primary/30"
          type="button"
        >
          <span class="border-r px-2 py-2 font-black text-muted-foreground/45">{{ event.timestamp }}</span>
          <span
            class="border-r px-2 py-2 font-black uppercase"
            :class="{
              'text-blue-600 dark:text-blue-400': event.direction === 'out',
              'text-primary': event.direction === 'in',
              'text-muted-foreground': event.direction === 'system',
              'text-destructive': event.direction === 'error',
            }"
          >
            {{ event.direction }}
          </span>
          <span class="min-w-0 px-2 py-2">
            <span class="block truncate font-black uppercase tracking-widest text-foreground/80">{{ event.title }}</span>
            <span v-if="event.payload" class="mt-1 block truncate text-muted-foreground">{{ event.payload }}</span>
          </span>
          <span class="px-2 py-2 text-right font-black text-muted-foreground/45">
            {{ event.sizeBytes ? `${event.sizeBytes} B` : '-' }}
          </span>
        </button>
      </div>
    </div>
  </div>

  <div v-else-if="workbench.responseData.value" class="flex h-full flex-col overflow-hidden">
    <!-- Elite Metadata Strip -->
    <div class="flex h-9 shrink-0 items-center justify-between border-b bg-muted/30 px-3">
      <!-- Left: Primary Status -->
      <div class="flex items-center gap-2">
        <UiBadge variant="outline" class="h-5 rounded-none border-emerald-500/30 bg-emerald-500/10 font-mono text-[9px] font-black uppercase text-emerald-600 dark:text-emerald-400">
          {{ workbench.responseData.value.status }} {{ workbench.responseData.value.statusText }}
        </UiBadge>
        <div class="flex items-center gap-1.5 font-mono text-[9px] font-bold text-muted-foreground/60 uppercase">
          <PhClock class="size-3" />
          {{ workbench.responseData.value.duration }}ms
          <span class="px-1 opacity-20">/</span>
          {{ workbench.responseData.value.size }}
        </div>
      </div>

      <!-- Right: Technical Metrics (Dense) -->
      <div class="flex items-center gap-4 font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/40">
        <div class="flex items-center gap-1.5">
          <PhDatabase class="size-2.5" />
          Target: <span class="text-foreground/60">{{ workbench.responseData.value.executionTarget }}</span>
        </div>
        <div class="flex items-center gap-1.5">
          <PhClock class="size-2.5" />
          TTFB: <span class="text-foreground/60">{{ workbench.responseData.value.ttfb }}ms</span>
        </div>
        <div class="hidden items-center gap-1.5 sm:flex">
          <PhFingerprint class="size-2.5" />
          ID: <span class="text-foreground/60">{{ workbench.responseData.value.requestId.split('_')[1] }}</span>
        </div>
      </div>
    </div>

    <!-- View Controls -->
    <div class="flex h-7 shrink-0 items-center border-b bg-background px-2">
      <button 
        class="flex h-full items-center gap-1.5 px-2 text-[9px] font-black uppercase tracking-widest transition-colors"
        :class="activeView === 'body' ? 'text-primary shadow-[inset_0_-1px_0_0_rgba(16,185,129,1)]' : 'text-muted-foreground/60 hover:text-foreground'"
        @click="activeView = 'body'"
      >
        <PhBracketsCurly class="size-3" />
        Body
      </button>
      <button 
        class="flex h-full items-center gap-1.5 px-2 text-[9px] font-black uppercase tracking-widest transition-colors"
        :class="activeView === 'headers' ? 'text-primary shadow-[inset_0_-1px_0_0_rgba(16,185,129,1)]' : 'text-muted-foreground/60 hover:text-foreground'"
        @click="activeView = 'headers'"
      >
        <PhListDashes class="size-3" />
        Headers
        <span class="ml-1 rounded-full bg-muted-foreground/10 px-1 text-[8px] opacity-60">{{ workbench.responseData.value.headers.length }}</span>
      </button>
    </div>

    <!-- Inspector Content -->
    <div class="min-h-0 flex-1 bg-card/30">
      <!-- Body Inspector -->
      <div v-if="activeView === 'body'" class="flex h-full">
        <WorkbenchCodeSurface
          :model-value="workbench.responseData.value.body"
          v-model:language="responseLanguage"
          class="min-w-0 flex-1"
          label="Response body"
        />
      </div>

      <!-- Headers Inspector -->
      <div v-else class="h-full overflow-auto p-4">
        <div class="grid gap-px border border-primary/10 bg-primary/10">
          <div 
            v-for="header in workbench.responseData.value.headers" 
            :key="header.id"
            class="grid grid-cols-[200px_1fr] bg-background p-2 font-mono text-[10px]"
          >
            <div class="font-black uppercase tracking-tighter text-primary/70">{{ header.key }}</div>
            <div class="truncate text-muted-foreground">{{ header.value }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="flex h-full items-center justify-center p-8 text-center text-muted-foreground">
    <div class="max-w-xs space-y-3">
      <div class="relative mx-auto size-12">
        <PhCode class="size-full opacity-10" />
        <div class="absolute inset-0 animate-pulse border-2 border-primary/10" />
      </div>
      <p class="font-mono text-[10px] font-bold uppercase tracking-widest opacity-40">Awaiting Execution...</p>
    </div>
  </div>
</template>
