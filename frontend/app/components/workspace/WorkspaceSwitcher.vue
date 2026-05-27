<script setup lang="ts">
import { PhCaretDown, PhCheck, PhPlus } from '@phosphor-icons/vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

const { workspaces, currentWorkspaceId, currentWorkspace } = useWorkspace()

const emit = defineEmits<{
  (e: 'create'): void
}>()
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <button class="flex h-8 w-full items-center gap-2 border border-border/40 bg-muted/20 px-3 font-mono text-[11px] font-black uppercase tracking-widest text-foreground hover:border-primary/30 hover:bg-muted/40 transition-all outline-none">
        <span class="truncate flex-1 text-left">{{ currentWorkspace?.name ?? 'Select workspace' }}</span>
        <PhCaretDown class="size-3 shrink-0 text-muted-foreground/60" />
      </button>
    </DropdownMenuTrigger>
    <DropdownMenuContent class="w-64 rounded-none border-2 border-primary/20 bg-background p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="start">
      <div class="px-3 py-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80 border-b border-border/30 mb-1">
        Workspaces
      </div>
      <DropdownMenuItem
        v-for="ws in workspaces"
        :key="ws.id"
        class="flex items-center justify-between font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2.5 rounded-none border-l-2 transition-all mb-0.5 last:mb-0 cursor-pointer"
        :class="ws.id === currentWorkspaceId ? 'bg-primary/10 text-primary border-primary' : 'border-transparent text-muted-foreground hover:bg-primary/5 hover:text-foreground hover:border-primary/30'"
        @select="currentWorkspaceId = ws.id"
      >
        <span class="truncate">{{ ws.name }}</span>
        <PhCheck v-if="ws.id === currentWorkspaceId" class="size-3 shrink-0 text-primary" />
      </DropdownMenuItem>
      <div class="border-t border-border/30 mt-1 pt-1">
        <DropdownMenuItem
          class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2.5 rounded-none border-l-2 border-transparent text-primary hover:bg-primary/10 hover:border-primary transition-all cursor-pointer"
          @select="emit('create')"
        >
          <PhPlus class="size-4" />
          Create workspace
        </DropdownMenuItem>
      </div>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
