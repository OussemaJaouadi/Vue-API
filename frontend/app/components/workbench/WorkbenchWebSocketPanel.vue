<script setup lang="ts">
import {
  PhBracketsCurly,
  PhPaperPlaneTilt,
  PhPlug,
} from '@phosphor-icons/vue'

const workbench = useWorkbench()

const canSend = computed(() => workbench.webSocketState.value === 'connected' && workbench.webSocketMessage.value.trim().length > 0)

const sendMessage = async () => {
  if (!canSend.value) return
  await workbench.sendWebSocketMessage()
}
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden bg-card/10 select-none">
    <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-4">
      <div class="flex items-center gap-3">
        <div class="grid size-6 place-items-center border-2 border-primary/20 bg-primary/5 text-primary">
          <PhPlug class="size-3.5" />
        </div>
        <div class="min-w-0">
          <span class="block font-mono text-[10px] font-black uppercase tracking-widest text-foreground/80 leading-none">Socket Uplink</span>
          <span class="mt-1 block font-mono text-[8px] font-bold uppercase tracking-widest text-muted-foreground/65 leading-none">Outbound Frame Constructor</span>
        </div>
      </div>

      <div 
        class="flex h-6 items-center px-2 border-2 font-mono text-[9px] font-black uppercase tracking-widest transition-all"
        :class="workbench.webSocketState.value === 'connected' ? 'border-primary/30 bg-primary/10 text-primary shadow-[0_0_10px_rgba(16,185,129,0.1)]' : 'border-border/60 text-muted-foreground/65'"
      >
        {{ workbench.webSocketState.value }}
      </div>
    </div>

    <div class="min-h-0 flex-1 border-b">
      <WorkbenchCodeSurface
        v-model="workbench.webSocketMessage.value"
        v-model:language="workbench.webSocketMessageLanguage.value"
        class="min-w-0"
        editable
        label="Payload"
      />
    </div>

    <div class="flex h-14 shrink-0 items-center justify-between bg-muted/20 px-4">
      <div class="flex items-center gap-2 font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/65 italic">
        <PhBracketsCurly class="size-3.5" />
        Frames are routed through backend proxy
      </div>

      <button
        class="group relative flex h-10 items-center justify-center gap-3 border-2 transition-all outline-none"
        :class="canSend 
          ? 'border-primary/30 bg-primary text-primary-foreground hover:bg-primary/90 hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.15)] active:translate-x-0.5 active:translate-y-0.5' 
          : 'border-border/40 bg-muted/10 text-muted-foreground/60 cursor-not-allowed'"
        type="button"
        :disabled="!canSend"
        @click="sendMessage"
      >
        <div class="flex items-center gap-2 px-5 font-mono text-[10px] font-black uppercase tracking-widest">
          <PhPaperPlaneTilt class="size-3.5 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
          Emit Frame
        </div>
      </button>
    </div>
  </div>
</template>
