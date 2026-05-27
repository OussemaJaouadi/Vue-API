<script setup lang="ts">
import {
  PhDotsSixVertical,
  PhX,
} from '@phosphor-icons/vue'
import type { HeaderItem } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const draggedHeaderId = ref<string | null>(null)
const dropTarget = ref<{ id: string, position: 'before' | 'after' } | null>(null)

const toggleHeader = (header: HeaderItem) => {
  header.enabled = !header.enabled
}

const isGhostHeader = (header: HeaderItem, index: number) => {
  const isLast = index === workbench.headers.value.length - 1
  return isLast && !header.key && !header.value
}

const startDrag = (header: HeaderItem, index: number, event: DragEvent) => {
  if (isGhostHeader(header, index)) return

  draggedHeaderId.value = header.id
  event.dataTransfer?.setData('text/plain', header.id)
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

const updateDropTarget = (header: HeaderItem, event: DragEvent) => {
  if (!draggedHeaderId.value || draggedHeaderId.value === header.id) return

  const row = event.currentTarget as HTMLElement
  const rect = row.getBoundingClientRect()
  const position = event.clientY < rect.top + rect.height / 2 ? 'before' : 'after'
  dropTarget.value = { id: header.id, position }
}

const dropHeader = (header: HeaderItem, event: DragEvent) => {
  const sourceId = draggedHeaderId.value || event.dataTransfer?.getData('text/plain')
  if (!sourceId || sourceId === header.id) return

  const targetIndex = workbench.headers.value.findIndex(item => item.id === header.id)
  if (targetIndex === -1) return

  const insertIndex = dropTarget.value?.position === 'after' ? targetIndex + 1 : targetIndex
  workbench.moveHeader(sourceId, insertIndex)
  draggedHeaderId.value = null
  dropTarget.value = null
}

const endDrag = () => {
  draggedHeaderId.value = null
  dropTarget.value = null
}

// Watch for changes on rows to auto-add ghost row
watch(workbench.headers.value, (newHeaders) => {
  const lastHeader = newHeaders[newHeaders.length - 1]
  if (lastHeader && (lastHeader.key !== '' || lastHeader.value !== '')) {
    workbench.addHeader()
  }
}, { deep: true })

onMounted(() => {
  // Ensure we start with at least one empty row if none exist
  if (workbench.headers.value.length === 0) {
    workbench.addHeader()
  }
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden border-0 bg-card select-none">
    <div class="grid shrink-0 grid-cols-[36px_42px_minmax(120px,1fr)_minmax(120px,1.3fr)_36px] border-b bg-muted/40 py-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
      <div class="flex items-center justify-center border-r border-border/10">#</div>
      <div class="flex items-center justify-center border-r border-border/10">Use</div>
      <span class="flex items-center border-r border-border/10 px-4">Key</span>
      <span class="flex items-center border-r border-border/10 px-4">Value</span>
      <span />
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto custom-scrollbar">
      <div
        v-for="(header, index) in workbench.headers.value"
        :key="header.id"
        class="group relative grid min-h-10 grid-cols-[36px_42px_minmax(120px,1fr)_minmax(120px,1.3fr)_36px] items-stretch border-b transition-colors"
        :class="[
          header.enabled ? 'bg-background hover:bg-primary/[0.03]' : 'bg-muted/10 text-muted-foreground/65 hover:bg-muted/20',
          draggedHeaderId === header.id && 'opacity-35',
        ]"
        :draggable="!isGhostHeader(header, index)"
        @dragstart="startDrag(header, index, $event)"
        @dragend="endDrag"
        @dragover.prevent="updateDropTarget(header, $event)"
        @dragleave="dropTarget?.id === header.id && (dropTarget = null)"
        @drop.prevent="dropHeader(header, $event)"
      >
        <div
          v-if="dropTarget?.id === header.id"
          class="pointer-events-none absolute inset-x-0 z-10 h-0.5 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.45)]"
          :class="dropTarget.position === 'before' ? 'top-0' : 'bottom-0'"
        />

        <div
          class="flex cursor-grab items-center justify-center border-r border-border/5 text-[9px] font-black opacity-70 group-hover:opacity-100 active:cursor-grabbing"
          :class="isGhostHeader(header, index) && 'invisible'"
        >
          <PhDotsSixVertical class="size-3.5" />
        </div>

        <button
          class="flex items-center justify-center border-r border-border/5 outline-none transition-colors"
          type="button"
          @click="toggleHeader(header)"
        >
          <div
            class="flex size-4 items-center justify-center border-2 transition-all"
            :class="header.enabled ? 'border-primary bg-primary text-background' : 'border-muted-foreground/20 hover:border-primary/40'"
          >
            <div v-if="header.enabled" class="size-2 bg-background" />
          </div>
        </button>

        <div class="flex items-center border-r border-border/5">
          <input
            v-model="header.key"
            class="h-full w-full bg-transparent px-4 font-mono text-[11px] font-black uppercase tracking-tight outline-none placeholder:text-muted-foreground/50"
            placeholder="Header-Key"
          >
        </div>

        <div class="flex items-center border-r border-border/5">
          <input
            v-model="header.value"
            class="h-full w-full bg-transparent px-4 font-mono text-[11px] font-bold outline-none placeholder:text-muted-foreground/50"
            placeholder="Value"
          >
        </div>

        <button
          class="flex items-center justify-center text-muted-foreground/20 outline-none transition-all hover:bg-destructive/10 hover:text-destructive group-hover:opacity-100"
          :class="isGhostHeader(header, index) ? 'invisible' : 'opacity-0'"
          type="button"
          @click="workbench.removeHeader(header.id)"
        >
          <PhX class="size-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
