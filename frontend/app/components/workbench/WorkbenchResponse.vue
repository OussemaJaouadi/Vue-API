<script setup lang="ts">
import {
  PhBracketsCurly,
  PhClock,
  PhCode,
  PhDatabase,
  PhFingerprint,
  PhListDashes,
  PhPlug,
} from '@phosphor-icons/vue'
import type { BodyLanguage } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const activeView = ref<'body' | 'headers'>('body')
const responseLanguage = ref<BodyLanguage>('json')

const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

watch(
  () => workbench.responseData.value,
  (resp) => {
    if (!resp) return
    const contentType = resp.headers.find(h => h.key.toLowerCase() === 'content-type')?.value || ''
    if (contentType.includes('application/json')) {
      responseLanguage.value = 'json'
    } else if (contentType.includes('text/html')) {
      responseLanguage.value = 'html'
    } else if (contentType.includes('application/xml') || contentType.includes('text/xml')) {
      responseLanguage.value = 'xml'
    } else if (contentType.includes('application/x-yaml') || contentType.includes('text/yaml')) {
      responseLanguage.value = 'yaml'
    } else {
      responseLanguage.value = 'text'
    }
  },
  { immediate: true },
)
</script>

<template>
  <div v-if="workbench.isWebSocketRequest.value" class="flex h-full flex-col overflow-hidden select-none">
    <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-4">
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 border-2 px-2 py-0.5 font-mono text-[9px] font-black uppercase transition-all"
          :class="workbench.webSocketState.value === 'connected' ? 'border-primary/30 bg-primary/5 text-primary' : 'border-border/60 text-muted-foreground/40'"
        >
          {{ workbench.webSocketState.value }}
        </div>
        <div class="flex items-center gap-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/60">
          <PhPlug class="size-3.5" />
          Socket Stream
        </div>
      </div>
      <div class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/40">
        {{ workbench.webSocketEvents.value.length }} Frame Event{{ workbench.webSocketEvents.value.length === 1 ? '' : 's' }}
      </div>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto custom-scrollbar bg-card/10 p-3 space-y-2">
      <div
        v-for="event in workbench.webSocketEvents.value"
        :key="event.id"
        class="group grid grid-cols-[80px_80px_minmax(0,1fr)_80px] items-start border bg-background transition-all hover:border-primary/30 hover:shadow-sm"
      >
        <span class="border-r border-border/10 px-3 py-2.5 font-mono text-[10px] font-bold text-muted-foreground/40">{{ event.timestamp }}</span>
        <span
          class="border-r border-border/10 px-3 py-2.5 font-mono text-[10px] font-black uppercase tracking-tighter"
          :class="{
            'text-blue-600 dark:text-blue-400': event.direction === 'out',
            'text-primary': event.direction === 'in',
            'text-muted-foreground': event.direction === 'system',
            'text-destructive': event.direction === 'error',
          }"
        >
          {{ event.direction }}
        </span>
        <div class="min-w-0 px-3 py-2.5">
          <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight text-foreground/80 group-hover:text-foreground">{{ event.title }}</span>
          <span v-if="event.payload" class="mt-1 block truncate font-mono text-[10px] text-muted-foreground/60">{{ event.payload }}</span>
        </div>
        <span class="px-3 py-2.5 text-right font-mono text-[10px] font-black text-muted-foreground/30">
          {{ event.sizeBytes ? `${event.sizeBytes} B` : '-' }}
        </span>
      </div>

      <div v-if="workbench.webSocketEvents.value.length === 0" class="flex h-32 items-center justify-center font-mono text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground/10 italic">
        Awaiting Initial Handshake
      </div>
    </div>
  </div>

  <div v-else-if="workbench.responseData.value" class="flex h-full flex-col overflow-hidden select-none">
    <!-- Tactile Metadata Strip -->
    <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-4">
      <div class="flex items-center gap-3">
        <div 
          class="flex h-6 items-center px-2 border-2 font-mono text-[10px] font-black uppercase tracking-tight transition-all shadow-sm"
          :class="workbench.responseData.value.status < 400 ? 'border-emerald-500/30 bg-emerald-500/5 text-emerald-600 dark:text-emerald-400' : 'border-destructive/30 bg-destructive/5 text-destructive'"
        >
          {{ workbench.responseData.value.status }} {{ workbench.responseData.value.statusText }}
        </div>
        <div class="flex items-center gap-3 font-mono text-[9px] font-black text-muted-foreground/60 uppercase tracking-widest">
          <div class="flex items-center gap-1.5">
            <PhClock class="size-3.5" />
            {{ workbench.responseData.value.duration.toFixed(0) }}ms
          </div>
          <span class="opacity-20">/</span>
          <div class="flex items-center gap-1.5">
            {{ formatSize(workbench.responseData.value.size) }}
          </div>
        </div>
      </div>

      <div class="flex items-center gap-6 font-mono text-[9px] font-black uppercase tracking-[0.15em] text-muted-foreground/40">
        <div class="hidden items-center gap-2 lg:flex">
          <PhDatabase class="size-3" />
          <span class="text-foreground/40">{{ workbench.responseData.value.executionTarget }}</span>
        </div>
        <div class="flex items-center gap-2">
          <PhFingerprint class="size-3" />
          <span class="text-foreground/40">{{ workbench.responseData.value.requestId.replace('req_', '') }}</span>
        </div>
      </div>
    </div>

    <!-- View Controls -->
    <div class="flex h-9 shrink-0 items-center border-b bg-background px-2 gap-px">
      <button 
        class="group relative h-full flex items-center gap-2 px-4 font-mono text-[10px] font-black uppercase tracking-widest transition-all outline-none"
        :class="activeView === 'body' ? 'bg-primary/5 text-primary' : 'text-muted-foreground/50 hover:bg-muted/30 hover:text-foreground'"
        @click="activeView = 'body'"
      >
        <PhBracketsCurly class="size-3.5 transition-colors" :class="activeView === 'body' ? 'text-primary' : 'opacity-40 group-hover:opacity-100'" />
        Body
        <div v-if="activeView === 'body'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
      </button>
      <button 
        class="group relative h-full flex items-center gap-2 px-4 font-mono text-[10px] font-black uppercase tracking-widest transition-all outline-none"
        :class="activeView === 'headers' ? 'bg-primary/5 text-primary' : 'text-muted-foreground/50 hover:bg-muted/30 hover:text-foreground'"
        @click="activeView = 'headers'"
      >
        <PhListDashes class="size-3.5 transition-colors" :class="activeView === 'headers' ? 'text-primary' : 'opacity-40 group-hover:opacity-100'" />
        Headers
        <span class="font-mono text-[9px] opacity-70">{{ workbench.responseData.value.headers.length }}</span>
        <div v-if="activeView === 'headers'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
      </button>
    </div>

    <!-- Inspector Content -->
    <div class="min-h-0 flex-1 bg-card/10">
      <div v-if="activeView === 'body'" class="flex h-full">
        <WorkbenchCodeSurface
          :model-value="workbench.responseData.value.body"
          v-model:language="responseLanguage"
          class="min-w-0 flex-1"
          label="Response Payload"
        />
      </div>

      <div v-else class="h-full overflow-y-auto custom-scrollbar p-3">
        <div class="border bg-background shadow-sm overflow-hidden">
          <div class="grid grid-cols-[220px_1fr] border-b bg-muted/40 py-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/60">
            <span class="px-4">Key</span>
            <span class="px-4">Value</span>
          </div>
          <div class="divide-y divide-border/10">
            <div 
              v-for="header in workbench.responseData.value.headers" 
              :key="header.id"
              class="group grid grid-cols-[220px_1fr] font-mono text-[11px] transition-colors hover:bg-primary/[0.02]"
            >
              <div class="border-r border-border/5 px-4 py-2.5 font-black uppercase tracking-tight text-primary/70 group-hover:text-primary transition-colors">{{ header.key }}</div>
              <div class="min-w-0 truncate px-4 py-2.5 text-foreground/90 group-hover:text-foreground transition-colors font-bold">{{ header.value }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="flex h-full items-center justify-center p-8 text-center bg-card/5">
    <div class="max-w-xs space-y-4">
      <div class="relative mx-auto size-14">
        <PhCode class="size-full text-primary opacity-5" />
        <div class="absolute inset-0 animate-pulse border-2 border-dashed border-primary/20" />
      </div>
      <p class="font-mono text-[11px] font-black uppercase tracking-[0.25em] text-muted-foreground/30">
        Terminal Idle / Awaiting Frame
      </p>
    </div>
  </div>
</template>
