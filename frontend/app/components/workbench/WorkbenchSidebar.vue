<script setup lang="ts">
import {
  PhCaretRight,
  PhCaretDown,
  PhFolder,
  PhPlus,
  PhGlobe,
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
  GET: 'border-emerald-500/20 bg-emerald-500/5 text-emerald-600 dark:text-emerald-400',
  POST: 'border-blue-500/20 bg-blue-500/5 text-blue-600 dark:text-blue-400',
  PUT: 'border-amber-500/20 bg-amber-500/5 text-amber-600 dark:text-amber-400',
  PATCH: 'border-purple-500/20 bg-purple-500/5 text-purple-600 dark:text-purple-400',
  DELETE: 'border-destructive/20 bg-destructive/5 text-destructive',
  SOCKET: 'border-primary/30 bg-primary/10 text-primary',
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

const getDropPosition = (e: DragEvent, targetIndex: number, totalItems: number): 'before' | 'after' => {
  const targetEl = (e.currentTarget as HTMLElement)
  const rect = targetEl.getBoundingClientRect()
  const y = e.clientY
  const threshold = rect.top + rect.height / 2
  return y < threshold ? 'before' : 'after'
}

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
  <aside class="hidden w-64 flex-col border-r bg-card/30 lg:flex select-none overflow-hidden">
    <!-- Project Root Header -->
    <div 
      class="flex h-14 items-center justify-between border-b bg-muted/30 px-4"
      @dragover.prevent
      @dragover="handleDragOverRoot($event, workbench.rootRequests.value.length, workbench.rootRequests.value.length)"
      @dragleave="handleDragLeave"
      @drop="handleDrop($event)"
    >
      <div class="flex items-center gap-3 overflow-hidden">
        <div class="grid size-8 place-items-center bg-primary border-2 border-primary shadow-[2px_2px_0_0_rgba(16,185,129,0.2)] text-primary-foreground shrink-0 transition-transform hover:scale-105">
          <PhGlobe class="size-4.5" />
        </div>
        <div class="flex flex-col leading-tight overflow-hidden">
          <span class="truncate font-mono text-[11px] font-black uppercase tracking-tight">Core Registry</span>
          <span class="text-[9px] font-bold text-primary/80 uppercase tracking-widest">Workbench</span>
        </div>
      </div>
      
      <div class="flex items-center gap-1">
        <button @click="workbench.addRequest()" class="p-1.5 text-muted-foreground/80 hover:text-primary transition-all hover:bg-primary/5 active:scale-95" title="New Request">
          <PhFilePlus class="size-4" />
        </button>
        <button @click="workbench.addFolder()" class="p-1.5 text-muted-foreground/80 hover:text-primary transition-all hover:bg-primary/5 active:scale-95" title="New Collection">
          <PhFolderPlus class="size-4" />
        </button>
      </div>
    </div>

    <!-- Explorer Navigation -->
    <nav class="flex-1 overflow-y-auto p-1 custom-scrollbar">
      <!-- Root Requests -->
      <div class="space-y-0.5 mb-2">
        <div v-if="dragOverTarget === `root-0` && dragOverPosition === 'before'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
        
        <template v-for="(request, index) in workbench.rootRequests.value" :key="request.id">
          <div v-if="index > 0 && dragOverTarget === `root-${index}` && dragOverPosition === 'before'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
          
          <button
            draggable="true"
            @dragstart="handleDragStart($event, request.id)"
            @click="workbench.openRequest(request)"
            @dragover.prevent
            @dragover="handleDragOverRoot($event, index, workbench.rootRequests.value.length)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, undefined, index, dragOverPosition)"
            class="group relative flex h-10 w-full items-center gap-3 px-3 transition-all duration-200 outline-none"
            :class="request.id === workbench.activeRequestId.value ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/3 hover:text-foreground'"
          >
            <div v-if="request.id === workbench.activeRequestId.value" class="wb-active-indicator" />
            
            <div class="flex h-5 w-10 shrink-0 items-center justify-center border-2 font-mono text-[8px] font-black uppercase tracking-tighter shadow-sm" :class="methodStyles[request.method]">
               {{ methodAbreviations[request.method] }}
            </div>

            <div class="min-w-0 flex-1 text-left">
              <input 
                v-if="renamingId === request.id"
                v-model="renameValue"
                class="w-full bg-background border-2 border-primary/40 px-1 outline-none font-bold text-[11px] tracking-tight"
                @blur="handleRename(request, false)"
                @keyup.enter="handleRename(request, false)"
                v-focus
              >
              <span v-else class="block truncate font-bold text-[11px] tracking-tight uppercase">{{ request.name }}</span>
            </div>
            
            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
              <button @click.stop="startRename(request.id, request.name)" class="p-1 hover:text-primary active:scale-90"><PhNotePencil class="size-3.5" /></button>
              <button @click.stop="workbench.deleteItem(request.id, false)" class="p-1 hover:text-destructive active:scale-90"><PhTrash class="size-3.5" /></button>
            </div>
          </button>
        </template>
        
        <div v-if="dragOverTarget === `root-${workbench.rootRequests.value.length - 1}` && dragOverPosition === 'after'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
      </div>

      <!-- Folders & Nested -->
      <div v-for="group in workbench.treeItems.value" :key="group.name" class="mb-1">
        <div 
          class="group relative flex h-10 w-full items-center gap-2 px-2 transition-all duration-200 cursor-pointer outline-none"
          :class="dragOverTarget === `folder-header-${group.name}` ? 'bg-primary/5' : 'hover:bg-primary/3'"
          @click="toggleFolder(group.name)"
          @dragover.prevent
          @dragover="handleDragOverFolderHeader($event, group.name)"
          @dragleave="handleDragLeave"
          @drop="handleDrop($event, group.name, 0, 'inside')"
        >
          <div class="flex size-6 items-center justify-center">
            <component :is="expandedFolders[group.name] ? PhCaretDown : PhCaretRight" class="size-3 transition-colors" :class="expandedFolders[group.name] ? 'text-primary' : 'text-muted-foreground/30'" />
          </div>
          
          <div class="grid size-6 place-items-center border-2 transition-colors" :class="expandedFolders[group.name] ? 'border-primary/20 bg-primary/5 text-primary' : 'border-border/60 bg-muted/10 text-muted-foreground/80'">
            <PhFolder class="size-3.5" />
          </div>
          
          <div class="min-w-0 flex-1">
            <input 
              v-if="renamingId === group.name"
              v-model="renameValue"
              class="w-full bg-background border-2 border-primary/40 px-1 outline-none font-black text-[11px] uppercase tracking-widest"
              @blur="handleRename(group, true)"
              @keyup.enter="handleRename(group, true)"
              v-focus
            >
            <span v-else class="block truncate font-mono text-[11px] font-black uppercase tracking-widest" :class="expandedFolders[group.name] ? 'text-primary' : 'text-foreground/90'">{{ group.name }}</span>
          </div>

          <div v-if="dragOverTarget === `folder-header-${group.name}`" class="mr-2">
            <PhArrowElbowDownRight class="size-3.5 text-primary animate-bounce" />
          </div>
          
          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
            <button @click.stop="workbench.addRequest(group.name)" class="p-1 hover:text-primary" title="New Nested Request"><PhFilePlus class="size-3.5" /></button>
            <button @click.stop="startRename(group.name, group.name)" class="p-1 hover:text-primary"><PhNotePencil class="size-3.5" /></button>
            <button @click.stop="workbench.deleteItem(group.name, true)" class="p-1 hover:text-destructive"><PhTrash class="size-3.5" /></button>
          </div>
        </div>

        <!-- Nested Items -->
        <div v-if="expandedFolders[group.name]" class="ml-4 border-l-2 border-primary/5 pl-1 mt-0.5 space-y-0.5">
          <div v-if="dragOverTarget === `folder-${group.name}-0` && dragOverPosition === 'before'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
          
          <template v-for="(request, index) in group.requests" :key="request.id">
            <div v-if="index > 0 && dragOverTarget === `folder-${group.name}-${index}` && dragOverPosition === 'before'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
            
            <button
              draggable="true"
              @dragstart="handleDragStart($event, request.id)"
              @click="workbench.openRequest(request)"
              @dragover.prevent
              @dragover="handleDragOverFolderRequest($event, group.name, index, group.requests.length)"
              @dragleave="handleDragLeave"
              @drop.stop="handleDrop($event, group.name, index, dragOverPosition)"
              class="group relative flex h-9 w-full items-center gap-2.5 px-3 transition-all duration-200 outline-none"
              :class="request.id === workbench.activeRequestId.value ? 'bg-primary/10 text-foreground' : 'text-muted-foreground hover:bg-primary/[0.02] hover:text-foreground'"
            >
              <div v-if="request.id === workbench.activeRequestId.value" class="wb-active-indicator" />
              
              <div class="flex h-4.5 w-9 shrink-0 items-center justify-center border-2 font-mono text-[7px] font-black uppercase tracking-tighter" :class="methodStyles[request.method]">
                 {{ methodAbreviations[request.method] }}
              </div>
              
              <div class="min-w-0 flex-1 text-left">
                <input 
                  v-if="renamingId === request.id"
                  v-model="renameValue"
                  class="w-full bg-background border-2 border-primary/40 px-1 outline-none font-bold text-[10px] tracking-tight"
                  @blur="handleRename(request, false)"
                  @keyup.enter="handleRename(request, false)"
                  v-focus
                >
                <span v-else class="block truncate font-bold text-[10px] tracking-tight uppercase">{{ request.name }}</span>
              </div>

              <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
                <button @click.stop="startRename(request.id, request.name)" class="p-0.5 hover:text-primary"><PhNotePencil class="size-3" /></button>
                <button @click.stop="workbench.deleteItem(request.id, false)" class="p-0.5 hover:text-destructive"><PhTrash class="size-3" /></button>
              </div>
            </button>
          </template>
          
          <div v-if="dragOverTarget === `folder-${group.name}-${group.requests.length - 1}` && dragOverPosition === 'after'" class="h-0.75 mx-2 bg-primary/60 shadow-[0_0_8px_theme(colors.primary.DEFAULT)]" />
        </div>
      </div>
    </nav>
  </aside>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: var(--color-primary); opacity: 0.1; }
</style>
