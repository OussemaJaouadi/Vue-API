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
  <div class="flex h-full min-h-0 flex-col overflow-hidden bg-card/30">
    <div class="flex h-8 shrink-0 items-center justify-between border-b bg-background px-3">
      <div class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">
        <PhPlug class="size-3 text-primary" />
        Websocket message
      </div>

      <UiBadge
        variant="outline"
        class="h-5 rounded-none px-1.5 font-mono text-[9px] uppercase"
        :class="workbench.webSocketState.value === 'connected' ? 'border-primary/30 bg-primary/10 text-primary' : 'border-muted-foreground/20 text-muted-foreground'"
      >
        {{ workbench.webSocketState.value }}
      </UiBadge>
    </div>

    <div class="min-h-0 flex-1">
      <WorkbenchCodeSurface
        v-model="workbench.webSocketMessage.value"
        v-model:language="workbench.webSocketMessageLanguage.value"
        class="min-w-0"
        editable
        label="Websocket message"
      />
    </div>

    <div class="flex h-10 shrink-0 items-center justify-between border-t bg-background px-3">
      <div class="flex items-center gap-2 font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/60">
        <PhBracketsCurly class="size-3" />
        Sent through backend execution session
      </div>

      <button
        class="flex h-7 items-center gap-2 border px-3 font-mono text-[9px] font-black uppercase tracking-widest transition-colors"
        :class="canSend ? 'border-primary bg-primary text-primary-foreground hover:bg-primary/95' : 'cursor-not-allowed border-border text-muted-foreground/40'"
        type="button"
        :disabled="!canSend"
        @click="sendMessage"
      >
        <PhPaperPlaneTilt class="size-3" />
        Send
      </button>
    </div>
  </div>
</template>
