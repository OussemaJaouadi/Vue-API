<script setup lang="ts">
import {
  PhCaretRight,
  PhCaretDown,
  PhFolder,
  PhPlus,
  PhLightning,
  PhDotsThreeVertical,
  PhTrash,
  PhNotePencil,
  PhFilePlus,
  PhFolderPlus,
  PhArrowElbowDownRight,
} from '@phosphor-icons/vue'
import { type ApiMethod, METHOD_COLORS, WORKBENCH_ICONS } from '~/composables/useWorkbench'

const workbench = useWorkbench()

// Expanded state for folders
const expandedFolders = useState<Record<string, boolean>>('workbench:sidebar-expanded', () => ({
  'Authentication': true,
  'Realtime': true
}))

// Drag state for position indicators
const dragOverTarget = ref<string | null>(null)
const dragOverPosition = ref<'before' | 'after' | 'inside' | null>(null)

const toggleFolder = (name: string) => {
  expandedFolders.value[name] = !expandedFolders.value[name]
}

const methodAbreviations: Record<ApiMethod, string> = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PTCH',
  DELETE: 'DEL',
  SOCKET: 'WS',
}

const methodStyles: Record<ApiMethod, string> = {
  GET: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30',
  POST: 'bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/30',
  PUT: 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30',
  PATCH: 'bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/30',
  DELETE: 'bg-destructive/10 text-destructive border-destructive/30',
  SOCKET: 'bg-primary/20 text-primary border-primary/40',
}

const renamingId = ref<string | null>(null)
const renameValue = ref('')

const startRename = (id: string, currentName: string) => {
  renamingId.value = id
  renameValue.value = currentName
}

const vFocus = {
  mounted: (el: HTMLInputElement) => el.focus(),
}

const handleRename = (item: { name: string }, isFolder: boolean) => {
  if (renameValue.value.trim()) {
    item.name = renameValue.value.trim()
  }
  renamingId.value = null
}

const handleDragStart = (e: DragEvent, requestId: string) => {
  if (e.dataTransfer) {
    e.dataTransfer.setData('requestId', requestId)
    e.dataTransfer.effectAllowed = 'move'
  }
}

// Calculate drop position based on mouse position
const getDropPosition = (e: DragEvent, targetIndex: number, totalItems: number): 'before' | 'after' => {
  const targetEl = (e.currentTarget as HTMLElement)
  const rect = targetEl.getBoundingClientRect()
  const y = e.clientY
  const threshold = rect.top + rect.height / 2
  return y < threshold ? 'before' : 'after'
}

// Handle drag over with position detection
const handleDragOverRoot = (e: DragEvent, index: number, totalItems: number) => {
  e.preventDefault()
  const position = getDropPosition(e, index, totalItems)
  dragOverTarget.value = `root-${index}`
  dragOverPosition.value = position
}

const handleDragOverFolderHeader = (e: DragEvent, folderName: string) => {
  e.preventDefault()
  dragOverTarget.value = `folder-header-${folderName}`
  dragOverPosition.value = 'inside'
}

const handleDragOverFolderRequest = (e: DragEvent, folderName: string, index: number, totalItems: number) => {
  e.preventDefault()
  const position = getDropPosition(e, index, totalItems)
  dragOverTarget.value = `folder-${folderName}-${index}`
  dragOverPosition.value = position
}

const handleDragLeave = () => {
  dragOverTarget.value = null
  dragOverPosition.value = null
}

// Perform drop with calculated position
const handleDrop = (e: DragEvent, targetFolderName?: string, targetIndex?: number, position?: 'before' | 'after' | 'inside' | null) => {
  const requestId = e.dataTransfer?.getData('requestId')
  if (!requestId) return
  const dropPosition = position ?? dragOverPosition.value
  
  dragOverTarget.value = null
  dragOverPosition.value = null
  
  let insertIndex = targetIndex
  if (insertIndex !== undefined && dropPosition === 'after') {
    insertIndex += 1
  }
  
  workbench.moveRequest(requestId, targetFolderName, insertIndex)
}
</script>

<template>
  <aside class="hidden min-w-0 flex-col border-r bg-card/30 lg:flex select-none">
    <!-- Project Root Header -->
    <div 
      class="flex h-12 items-center justify-between border-b bg-muted/30 px-4"
      @dragover.prevent
      @dragover="handleDragOverRoot($event, workbench.rootRequests.value.length, workbench.rootRequests.value.length)"
      @dragleave="handleDragLeave"
      @drop="handleDrop($event)"
    >
      <div class="flex items-center gap-3 overflow-hidden">
        <div class="grid size-6 place-items-center bg-primary border-2 border-primary shadow-[2px_2px_0_0_rgba(16,185,129,0.2)] text-primary-foreground">
          <component :is="WORKBENCH_ICONS['PhGlobe']" class="size-3.5" />
        </div>
        <div class="flex flex-col leading-none overflow-hidden">
          <span class="truncate font-mono text-[11px] font-black uppercase tracking-tight">Core API</span>
          <span class="text-[9px] font-bold text-primary/40 uppercase tracking-widest">Explorer</span>
        </div>
      </div>
      
      <!-- Global Actions -->
      <div class="flex items-center gap-1">
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <button @click="workbench.addRequest()" class="p-1 text-muted-foreground/40 hover:text-primary transition-all hover:bg-primary/5">
              <PhFilePlus class="size-3.5" />
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">New Root Request</UiTooltipContent>
        </UiTooltip>
        <UiTooltip>
          <UiTooltipTrigger as-child>
            <button @click="workbench.addFolder()" class="p-1 text-muted-foreground/40 hover:text-primary transition-all hover:bg-primary/5">
              <PhFolderPlus class="size-3.5" />
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">New Collection</UiTooltipContent>
        </UiTooltip>
      </div>
    </div>

    <!-- The VS Code Style Explorer -->
    <nav class="flex-1 overflow-y-auto p-1 custom-scrollbar">
      <!-- Root Requests -->
      <div class="space-y-0.5 mb-2">
        <!-- Drop indicator at root level (before first item or at end) -->
        <div 
          v-if="dragOverTarget === `root-0` && dragOverPosition === 'before'"
          class="h-0.5 bg-primary"
        />
        
        <template v-for="(request, index) in workbench.rootRequests.value" :key="request.id">
          <!-- Drop indicator before this item -->
          <div 
            v-if="dragOverTarget === `root-${index}` && dragOverPosition === 'before'"
            class="h-0.5 bg-primary"
          />
          
          <button
            :data-index="index"
            draggable="true"
            @dragstart="handleDragStart($event, request.id)"
            @click="workbench.openRequest(request)"
            @dragover.prevent
            @dragover="handleDragOverRoot($event, index, workbench.rootRequests.value.length)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, undefined, index, dragOverPosition)"
            class="group relative flex h-8 w-full items-center gap-2.5 px-2 transition-all cursor-grab active:cursor-grabbing"
            :class="request.id === workbench.activeRequestId.value ? 'bg-primary/10 text-foreground shadow-sm' : 'text-muted-foreground/70 hover:bg-primary/3'"
          >
            <div class="flex h-4.5 w-10 items-center justify-center border font-mono text-[8px] font-black uppercase tracking-tighter" :class="methodStyles[request.method]">
               {{ methodAbreviations[request.method] }}
            </div>
            <input 
              v-if="renamingId === request.id"
              v-model="renameValue"
              class="min-w-0 flex-1 bg-background border border-primary/40 px-1 outline-none font-bold text-[10px]"
              @blur="handleRename(request, false)"
              @keyup.enter="handleRename(request, false)"
              v-focus
            >
            <span v-else class="truncate font-bold text-[10px] tracking-tight">{{ request.name }}</span>
            
            <!-- Inline Actions -->
            <div class="ml-auto flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
              <button @click.stop="startRename(request.id, request.name)" class="p-1 hover:text-primary"><PhNotePencil class="size-3" /></button>
              <button @click.stop="workbench.deleteItem(request.id, false)" class="p-1 hover:text-destructive"><PhTrash class="size-3" /></button>
            </div>
          </button>
        </template>
        
        <!-- Drop indicator at end of root requests -->
        <div 
          v-if="dragOverTarget === `root-${workbench.rootRequests.value.length - 1}` && dragOverPosition === 'after'"
          class="h-0.5 bg-primary"
        />
      </div>

      <!-- Folders & Nested Requests -->
      <div v-for="group in workbench.treeItems.value" :key="group.name" class="mb-0.5">
        <div 
          class="group flex w-full items-center gap-2 px-2 py-1 text-left transition-all hover:bg-primary/4"
          @click="toggleFolder(group.name)"
          @dragover.prevent
          @dragover="handleDragOverFolderHeader($event, group.name)"
          @dragleave="handleDragLeave"
          @drop="handleDrop($event, group.name, 0, 'inside')"
        >
          <component :is="expandedFolders[group.name] ? PhCaretDown : PhCaretRight" class="size-3 text-muted-foreground/40" />
          <PhFolder class="size-4 text-primary/60" />
          
          <!-- Inside indicator -->
          <div v-if="dragOverTarget === `folder-header-${group.name}`" class="ml-auto">
            <PhArrowElbowDownRight class="size-3 text-primary" />
          </div>
          
          <input 
            v-if="renamingId === group.name"
            v-model="renameValue"
            class="min-w-0 flex-1 bg-background border border-primary/40 px-1 outline-none font-black text-[10px] uppercase"
            @blur="handleRename(group, true)"
            @keyup.enter="handleRename(group, true)"
            v-focus
          >
          <span v-else class="truncate font-mono text-[10px] font-black uppercase tracking-widest text-foreground/70">{{ group.name }}</span>
          
          <div class="ml-auto flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
            <button @click.stop="workbench.addRequest(group.name)" class="p-1 hover:text-primary"><PhFilePlus class="size-3" /></button>
            <button @click.stop="startRename(group.name, group.name)" class="p-1 hover:text-primary"><PhNotePencil class="size-3" /></button>
            <button @click.stop="workbench.deleteItem(group.name, true)" class="p-1 hover:text-destructive"><PhTrash class="size-3" /></button>
          </div>
        </div>

        <!-- Nested Requests -->
        <div v-if="expandedFolders[group.name]" class="ml-4 border-l border-primary/10 pl-1 mt-0.5 space-y-0.5">
          <!-- Drop indicator at folder start -->
          <div 
            v-if="dragOverTarget === `folder-${group.name}-0` && dragOverPosition === 'before'"
            class="h-0.5 bg-primary"
          />
          
          <template v-for="(request, index) in group.requests" :key="request.id">
            <!-- Drop indicator before this item -->
            <div 
              v-if="dragOverTarget === `folder-${group.name}-${index}` && dragOverPosition === 'before'"
              class="h-0.5 bg-primary"
            />
            
            <button
              draggable="true"
              @dragstart="handleDragStart($event, request.id)"
              @click="workbench.openRequest(request)"
              @dragover.prevent
              @dragover="handleDragOverFolderRequest($event, group.name, index, group.requests.length)"
              @dragleave="handleDragLeave"
              @drop.stop="handleDrop($event, group.name, index, dragOverPosition)"
              class="group relative flex h-8 w-full items-center gap-2.5 px-2 transition-all cursor-grab active:cursor-grabbing"
              :class="request.id === workbench.activeRequestId.value ? 'bg-primary/10 text-foreground shadow-sm' : 'text-muted-foreground/70 hover:bg-primary/3'"
            >
              <div class="flex h-4.5 w-10 items-center justify-center border font-mono text-[8px] font-black uppercase tracking-tighter" :class="methodStyles[request.method]">
                 {{ methodAbreviations[request.method] }}
              </div>
              
              <input 
                v-if="renamingId === request.id"
                v-model="renameValue"
                class="min-w-0 flex-1 bg-background border border-primary/40 px-1 outline-none font-bold text-[10px]"
                @blur="handleRename(request, false)"
                @keyup.enter="handleRename(request, false)"
                v-focus
              >
              <span v-else class="truncate font-bold text-[10px] tracking-tight">{{ request.name }}</span>

              <div class="ml-auto flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click.stop="startRename(request.id, request.name)" class="p-1 hover:text-primary"><PhNotePencil class="size-3" /></button>
                <button @click.stop="workbench.deleteItem(request.id, false)" class="p-1 hover:text-destructive"><PhTrash class="size-3" /></button>
              </div>
            </button>
          </template>
          
          <!-- Drop indicator at end of folder -->
          <div 
            v-if="dragOverTarget === `folder-${group.name}-${group.requests.length - 1}` && dragOverPosition === 'after'"
            class="h-0.5 bg-primary"
          />
        </div>
      </div>
    </nav>
  </aside>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: var(--color-primary); opacity: 0.1; }
</style>
