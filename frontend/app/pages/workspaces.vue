<script setup lang="ts">
import {
  PhBuildings,
  PhCheck,
  PhPencilSimple,
  PhPlus,
  PhTrash,
  PhWarning,
  PhX,
} from '@phosphor-icons/vue'
import type { Workspace } from '~/types/workspace'

const {
  workspaces,
  currentWorkspaceId,
  currentWorkspace,
  workspacesLoading,
  workspacesError,
  loadWorkspaces,
  renameWorkspace,
  deleteWorkspace,
} = useWorkspace()

const createOpen = ref(false)
const editingId = ref('')
const editingName = ref('')
const savingId = ref('')
const deletingId = ref('')
const deleteTarget = ref<Workspace | null>(null)
const actionError = ref<string | null>(null)

const sortedWorkspaces = computed(() => workspaces.value)

const beginRename = (workspace: Workspace) => {
  actionError.value = null
  editingId.value = workspace.id
  editingName.value = workspace.name
}

const cancelRename = () => {
  editingId.value = ''
  editingName.value = ''
}

const saveRename = async (workspace: Workspace) => {
  const nextName = editingName.value.trim()
  if (!nextName || nextName === workspace.name || savingId.value) {
    cancelRename()
    return
  }

  savingId.value = workspace.id
  actionError.value = null
  try {
    await renameWorkspace(workspace.id, nextName)
    cancelRename()
  } catch (err: any) {
    actionError.value = err?.data?.error || err?.message || 'Failed to rename workspace'
  } finally {
    savingId.value = ''
  }
}

const requestDelete = (workspace: Workspace) => {
  actionError.value = null
  deleteTarget.value = workspace
}

const confirmDelete = async () => {
  if (!deleteTarget.value || deletingId.value) return

  deletingId.value = deleteTarget.value.id
  actionError.value = null
  try {
    await deleteWorkspace(deleteTarget.value.id)
    deleteTarget.value = null
  } catch (err: any) {
    actionError.value = err?.data?.error || err?.message || 'Failed to delete workspace'
  } finally {
    deletingId.value = ''
  }
}

onMounted(() => {
  if (workspaces.value.length === 0 && !workspacesLoading.value) {
    loadWorkspaces()
  }
})
</script>

<template>
  <div class="min-h-[calc(100dvh-5.5rem)] border bg-card">
    <header class="flex h-14 items-center justify-between border-b bg-muted/20 px-4">
      <div class="flex min-w-0 items-center gap-3">
        <div class="grid size-8 place-items-center border-2 border-primary/25 bg-primary/10 text-primary">
          <PhBuildings class="size-4" />
        </div>
        <div class="min-w-0">
          <h1 class="truncate font-heading text-base font-bold">Workspaces</h1>
          <p class="font-mono text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            CRUD / membership boundary
          </p>
        </div>
      </div>

      <UiButton
        class="rounded-none border-2 border-primary/25 bg-primary/5 font-mono text-[10px] font-black uppercase tracking-widest text-primary hover:bg-primary/10"
        type="button"
        variant="ghost"
        @click="createOpen = true"
      >
        <PhPlus class="mr-2 size-4" />
        Workspace
      </UiButton>
    </header>

    <div class="grid gap-3 p-3 lg:grid-cols-[minmax(0,1fr)_320px]">
      <section class="border bg-background">
        <div class="flex items-center justify-between border-b bg-muted/15 px-4 py-3">
          <div>
            <h2 class="font-mono text-[11px] font-black uppercase tracking-widest">Workspace registry</h2>
            <p class="mt-1 text-sm text-muted-foreground">Switch, rename, or remove isolated project spaces.</p>
          </div>
          <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">
            {{ workspaces.length }} total
          </span>
        </div>

        <div v-if="workspacesLoading" class="grid gap-2 p-4">
          <div v-for="idx in 4" :key="idx" class="h-20 animate-pulse border bg-muted/20" />
        </div>

        <div v-else-if="workspacesError" class="m-4 border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
          {{ workspacesError }}
        </div>

        <div v-else-if="sortedWorkspaces.length === 0" class="grid min-h-80 place-items-center p-6 text-center">
          <div class="max-w-sm">
            <div class="mx-auto grid size-12 place-items-center border-2 border-primary/25 bg-primary/10 text-primary">
              <PhBuildings class="size-6" />
            </div>
            <h3 class="mt-4 font-heading text-lg font-bold">No workspaces</h3>
            <p class="mt-2 text-sm text-muted-foreground">Create one to unlock collections, environments, access grants, and execution history.</p>
            <UiButton
              class="mt-4 rounded-none"
              type="button"
              @click="createOpen = true"
            >
              <PhPlus class="mr-2 size-4" />
              Create workspace
            </UiButton>
          </div>
        </div>

        <div v-else class="divide-y">
          <article
            v-for="workspace in sortedWorkspaces"
            :key="workspace.id"
            class="grid gap-3 p-4 transition-colors md:grid-cols-[minmax(0,1fr)_auto]"
            :class="workspace.id === currentWorkspaceId ? 'bg-primary/6' : 'hover:bg-muted/20'"
          >
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <button
                  class="grid size-8 place-items-center border transition-colors"
                  :class="workspace.id === currentWorkspaceId ? 'border-primary bg-primary/10 text-primary' : 'border-border text-muted-foreground hover:border-primary/40 hover:text-primary'"
                  type="button"
                  @click="currentWorkspaceId = workspace.id"
                >
                  <PhCheck v-if="workspace.id === currentWorkspaceId" class="size-4" />
                  <PhBuildings v-else class="size-4" />
                </button>

                <form
                  v-if="editingId === workspace.id"
                  class="flex min-w-0 flex-1 items-center gap-2"
                  @submit.prevent="saveRename(workspace)"
                >
                  <input
                    v-model="editingName"
                    class="h-9 min-w-0 flex-1 rounded-none border-2 border-primary/20 bg-background px-3 font-mono text-sm font-bold outline-none focus:border-primary"
                    type="text"
                  >
                  <UiButton class="rounded-none" size="sm" type="submit" :disabled="savingId === workspace.id">
                    {{ savingId === workspace.id ? 'Saving' : 'Save' }}
                  </UiButton>
                  <UiButton class="rounded-none" size="icon-sm" type="button" variant="ghost" @click="cancelRename">
                    <PhX class="size-4" />
                  </UiButton>
                </form>

                <template v-else>
                  <button
                    class="min-w-0 text-left"
                    type="button"
                    @click="currentWorkspaceId = workspace.id"
                  >
                    <span class="block truncate font-heading text-sm font-bold">{{ workspace.name }}</span>
                    <span class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
                      {{ workspace.role || 'member' }} / {{ workspace.memberCount ?? 0 }} members
                    </span>
                  </button>
                  <span
                    v-if="workspace.id === currentWorkspaceId"
                    class="border border-primary/25 bg-primary/10 px-2 py-1 font-mono text-[9px] font-black uppercase tracking-widest text-primary"
                  >
                    Active
                  </span>
                </template>
              </div>
            </div>

            <div class="flex items-center gap-2 md:justify-end">
              <UiButton
                class="rounded-none"
                size="icon-sm"
                type="button"
                variant="outline"
                @click="beginRename(workspace)"
              >
                <PhPencilSimple class="size-4" />
              </UiButton>
              <UiButton
                class="rounded-none border-destructive/30 text-destructive hover:bg-destructive/10"
                size="icon-sm"
                type="button"
                variant="outline"
                @click="requestDelete(workspace)"
              >
                <PhTrash class="size-4" />
              </UiButton>
            </div>
          </article>
        </div>
      </section>

      <aside class="border bg-background">
        <div class="border-b bg-muted/15 px-4 py-3">
          <h2 class="font-mono text-[11px] font-black uppercase tracking-widest">Current boundary</h2>
          <p class="mt-1 text-sm text-muted-foreground">Everything below this workspace is scoped.</p>
        </div>
        <div class="space-y-3 p-4">
          <div class="border bg-muted/15 p-3">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Active workspace</span>
            <p class="mt-1 truncate font-heading text-base font-bold">{{ currentWorkspace?.name ?? 'None' }}</p>
          </div>
          <div class="border bg-muted/15 p-3">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Delete behavior</span>
            <p class="mt-1 text-sm leading-5 text-muted-foreground">
              Removing a workspace also removes its collections, requests, environments, variables, memberships, and access grants.
            </p>
          </div>
          <p v-if="actionError" class="border border-destructive/30 bg-destructive/10 px-3 py-2 text-xs text-destructive">
            {{ actionError }}
          </p>
        </div>
      </aside>
    </div>

    <Teleport to="body">
      <div
        v-if="deleteTarget"
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
              <h2 id="delete-workspace-title" class="font-heading text-sm font-black text-foreground">Delete workspace</h2>
              <p class="mt-1 text-xs text-muted-foreground">This action removes all workspace-scoped data.</p>
            </div>
            <button
              class="grid size-7 shrink-0 place-items-center text-muted-foreground transition-colors hover:text-foreground"
              type="button"
              @click="deleteTarget = null"
            >
              <PhX class="size-4" />
            </button>
          </div>

          <div class="space-y-4 p-4">
            <div class="border border-border/70 bg-muted/20 p-3">
              <span class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Workspace</span>
              <span class="mt-1 block truncate text-sm font-bold text-foreground">{{ deleteTarget.name }}</span>
            </div>

            <div class="flex justify-end gap-2">
              <UiButton class="rounded-none" type="button" variant="outline" @click="deleteTarget = null">
                Cancel
              </UiButton>
              <UiButton
                class="rounded-none bg-destructive text-destructive-foreground hover:bg-destructive/90"
                type="button"
                :disabled="deletingId === deleteTarget.id"
                @click="confirmDelete"
              >
                {{ deletingId === deleteTarget.id ? 'Deleting...' : 'Delete workspace' }}
              </UiButton>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <WorkspaceCreateModal v-model:open="createOpen" />
  </div>
</template>
