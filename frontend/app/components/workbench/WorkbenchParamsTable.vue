<script setup lang="ts">
import {
  PhDotsSixVertical,
  PhX,
} from '@phosphor-icons/vue'
import type { QueryParamItem } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const draggedParamId = ref<string | null>(null)
const dropTarget = ref<{ id: string, position: 'before' | 'after' } | null>(null)

const toggleParam = (param: QueryParamItem) => {
  param.enabled = !param.enabled
}

const isGhostParam = (param: QueryParamItem, index: number) => {
  const isLast = index === workbench.queryParams.value.length - 1
  return isLast && !param.key && !param.value
}

const startDrag = (param: QueryParamItem, index: number, event: DragEvent) => {
  if (isGhostParam(param, index)) return

  draggedParamId.value = param.id
  event.dataTransfer?.setData('text/plain', param.id)
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

const updateDropTarget = (param: QueryParamItem, event: DragEvent) => {
  if (!draggedParamId.value || draggedParamId.value === param.id) return

  const row = event.currentTarget as HTMLElement
  const rect = row.getBoundingClientRect()
  const position = event.clientY < rect.top + rect.height / 2 ? 'before' : 'after'
  dropTarget.value = { id: param.id, position }
}

const dropParam = (param: QueryParamItem, event: DragEvent) => {
  const sourceId = draggedParamId.value || event.dataTransfer?.getData('text/plain')
  if (!sourceId || sourceId === param.id) return

  const targetIndex = workbench.queryParams.value.findIndex(item => item.id === param.id)
  if (targetIndex === -1) return

  const insertIndex = dropTarget.value?.position === 'after' ? targetIndex + 1 : targetIndex
  workbench.moveQueryParam(sourceId, insertIndex)
  draggedParamId.value = null
  dropTarget.value = null
}

const endDrag = () => {
  draggedParamId.value = null
  dropTarget.value = null
}

// Watch for changes on rows to auto-add ghost row
watch(workbench.queryParams.value, (newParams) => {
  const lastParam = newParams[newParams.length - 1]
  if (lastParam && (lastParam.key !== '' || lastParam.value !== '')) {
    workbench.addQueryParam()
  }
}, { deep: true })

onMounted(() => {
  // Ensure we start with at least one empty row if none exist
  if (workbench.queryParams.value.length === 0) {
    workbench.addQueryParam()
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
        v-for="(param, index) in workbench.queryParams.value"
        :key="param.id"
        class="group relative grid min-h-10 grid-cols-[36px_42px_minmax(120px,1fr)_minmax(120px,1.3fr)_36px] items-stretch border-b transition-colors"
        :class="[
          param.enabled ? 'bg-background hover:bg-primary/[0.03]' : 'bg-muted/10 text-muted-foreground/40 hover:bg-muted/20',
          draggedParamId === param.id && 'opacity-35',
        ]"
        :draggable="!isGhostParam(param, index)"
        @dragstart="startDrag(param, index, $event)"
        @dragend="endDrag"
        @dragover.prevent="updateDropTarget(param, $event)"
        @dragleave="dropTarget?.id === param.id && (dropTarget = null)"
        @drop.prevent="dropParam(param, $event)"
      >
        <div
          v-if="dropTarget?.id === param.id"
          class="pointer-events-none absolute inset-x-0 z-10 h-0.5 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.45)]"
          :class="dropTarget.position === 'before' ? 'top-0' : 'bottom-0'"
        />

        <div
          class="flex cursor-grab items-center justify-center border-r border-border/5 text-[9px] font-black opacity-70 group-hover:opacity-100 active:cursor-grabbing"
          :class="isGhostParam(param, index) && 'invisible'"
        >
          <PhDotsSixVertical class="size-3.5" />
        </div>

        <button
          class="flex items-center justify-center border-r border-border/5 outline-none transition-colors"
          type="button"
          @click="toggleParam(param)"
        >
          <div
            class="flex size-4 items-center justify-center border-2 transition-all"
            :class="param.enabled ? 'border-primary bg-primary text-background' : 'border-muted-foreground/20 hover:border-primary/40'"
          >
            <div v-if="param.enabled" class="size-2 bg-background" />
          </div>
        </button>

        <div class="flex items-center border-r border-border/5">
          <input
            v-model="param.key"
            class="h-full w-full bg-transparent px-4 font-mono text-[11px] font-black uppercase tracking-tight outline-none placeholder:text-muted-foreground/20"
            placeholder="Parameter Key"
          >
        </div>

        <div class="flex items-center border-r border-border/5">
          <input
            v-model="param.value"
            class="h-full w-full bg-transparent px-4 font-mono text-[11px] font-bold outline-none placeholder:text-muted-foreground/20"
            placeholder="Value"
          >
        </div>

        <button
          class="flex items-center justify-center text-muted-foreground/20 outline-none transition-all hover:bg-destructive/10 hover:text-destructive group-hover:opacity-100"
          :class="isGhostParam(param, index) ? 'invisible' : 'opacity-0'"
          type="button"
          @click="workbench.removeQueryParam(param.id)"
        >
          <PhX class="size-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
