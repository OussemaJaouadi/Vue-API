<script setup lang="ts">
import { PhCaretDown, PhCheck, PhPlus, PhTrash, PhWarning, PhX } from '@phosphor-icons/vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

const { workspaces, currentWorkspaceId, currentWorkspace, deleteWorkspace } = useWorkspace()

const emit = defineEmits<{
  (e: 'create'): void
}>()

const confirmOpen = ref(false)
const deleting = ref(false)
const deleteError = ref<string | null>(null)

const openDeleteConfirm = () => {
  deleteError.value = null
  confirmOpen.value = true
}

const confirmDelete = async () => {
  if (!currentWorkspace.value || deleting.value) return

  deleting.value = true
  deleteError.value = null
  try {
    await deleteWorkspace(currentWorkspace.value.id)
    confirmOpen.value = false
  } catch (err: any) {
    deleteError.value = err?.data?.error || err?.message || 'Failed to delete workspace'
  } finally {
    deleting.value = false
  }
}
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
        <DropdownMenuItem
          class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2.5 rounded-none border-l-2 border-transparent text-destructive hover:bg-destructive/10 hover:border-destructive transition-all cursor-pointer"
          @select="openDeleteConfirm"
        >
          <PhTrash class="size-4" />
          Delete current
        </DropdownMenuItem>
      </div>
    </DropdownMenuContent>
  </DropdownMenu>

  <Teleport to="body">
    <div
      v-if="confirmOpen"
      class="fixed inset-0 z-50 grid place-items-center bg-background/55 p-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
      aria-labelledby="delete-workspace-title"
    >
      <div class="w-full max-w-md rounded-none border-2 border-destructive/45 bg-background shadow-[8px_8px_0_0_rgba(239,68,68,0.18)]">
        <div class="flex items-start gap-3 border-b border-destructive/20 bg-destructive/5 p-4">
          <div class="grid size-9 shrink-0 place-items-center border border-destructive/35 bg-destructive/10 text-destructive">
            <PhWarning class="size-5" />
          </div>
          <div class="min-w-0 flex-1">
            <h2 id="delete-workspace-title" class="font-heading text-sm font-black text-foreground">
              Delete workspace
            </h2>
            <p class="mt-1 text-xs text-muted-foreground">
              This removes collections, requests, environments, grants, and memberships for this workspace.
            </p>
          </div>
          <button
            class="grid size-7 shrink-0 place-items-center text-muted-foreground transition-colors hover:text-foreground"
            type="button"
            @click="confirmOpen = false"
          >
            <PhX class="size-4" />
          </button>
        </div>

        <div class="space-y-4 p-4">
          <div class="border border-border/70 bg-muted/20 p-3">
            <span class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Workspace</span>
            <span class="mt-1 block truncate text-sm font-bold text-foreground">{{ currentWorkspace?.name }}</span>
          </div>

          <p v-if="deleteError" class="border border-destructive/30 bg-destructive/10 px-3 py-2 text-xs text-destructive">
            {{ deleteError }}
          </p>

          <div class="flex justify-end gap-2">
            <UiButton
              class="rounded-none"
              type="button"
              variant="outline"
              @click="confirmOpen = false"
            >
              Cancel
            </UiButton>
            <UiButton
              class="rounded-none bg-destructive text-destructive-foreground hover:bg-destructive/90"
              type="button"
              :disabled="deleting"
              @click="confirmDelete"
            >
              {{ deleting ? 'Deleting...' : 'Delete workspace' }}
            </UiButton>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
