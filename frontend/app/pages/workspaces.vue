<script setup lang="ts">
import {
  PhBuildings,
  PhCheck,
  PhPencilSimple,
  PhPlus,
  PhShieldCheck,
  PhTrash,
  PhUsersThree,
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
const selectedWorkspace = computed(() =>
  sortedWorkspaces.value.find(workspace => workspace.id === currentWorkspaceId.value) ?? sortedWorkspaces.value[0],
)
const selectedRole = computed(() => selectedWorkspace.value?.role || 'member')
const selectedMembers = computed(() => selectedWorkspace.value?.memberCount ?? 0)

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
  <div class="flex h-[calc(100dvh-5.5rem)] flex-col overflow-hidden border bg-card">
    <header class="flex h-16 shrink-0 items-center justify-between border-b bg-muted/30 px-6 select-none">
      <div class="flex items-center gap-4">
        <div class="grid size-10 place-items-center border-2 border-primary bg-primary/10 text-primary shadow-[3px_3px_0_0_rgba(16,185,129,0.15)] transition-transform hover:scale-105">
          <PhBuildings class="size-5" />
        </div>
        <div>
          <h1 class="font-heading text-lg font-black uppercase tracking-tight text-foreground">Workspaces</h1>
          <p class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
            Selected Context: {{ selectedWorkspace?.name ?? 'None' }}
          </p>
        </div>
      </div>

      <button
        class="btn-tactile-primary flex h-9 items-center gap-2 px-4 font-mono text-[10px] font-black uppercase tracking-widest outline-none"
        type="button"
        @click="createOpen = true"
      >
        <PhPlus class="size-3.5" />
        <span class="hidden sm:inline">New Workspace</span>
        <span class="sm:hidden">New</span>
      </button>
    </header>

    <div class="flex min-w-0 flex-1 gap-3 overflow-x-auto bg-muted/5 p-3">
      <aside class="flex w-64 shrink-0 select-none flex-col overflow-hidden border-r bg-card/30">
        <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-3">
          <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Registry</span>
          <span v-if="!workspacesLoading" class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/40">
            {{ workspaces.length }} nodes
          </span>
          <div v-else class="h-2 w-12 animate-pulse bg-muted-foreground/15" />
        </div>

        <div class="custom-scrollbar flex-1 space-y-0.5 overflow-y-auto p-1">
          <template v-if="workspacesLoading">
            <div v-for="idx in 3" :key="idx" class="flex h-14 w-full items-center gap-3 px-3">
              <div class="flex-1 space-y-1.5">
                <div class="h-3 w-24 animate-pulse bg-muted-foreground/15" />
                <div class="h-2 w-14 animate-pulse bg-muted-foreground/10" />
              </div>
            </div>
          </template>

          <button
            v-for="workspace in sortedWorkspaces"
            :key="workspace.id"
            class="group relative flex h-14 w-full shrink-0 items-center justify-between px-3 outline-none transition-all duration-200"
            :class="workspace.id === currentWorkspaceId ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/5 hover:text-foreground'"
            type="button"
            @click="currentWorkspaceId = workspace.id"
          >
            <div v-if="workspace.id === currentWorkspaceId" class="wb-active-indicator" />

            <div class="min-w-0 flex-1 text-left">
              <span
                class="block truncate font-mono text-[11px] font-black uppercase tracking-tight transition-colors"
                :class="workspace.id === currentWorkspaceId ? 'text-primary' : 'group-hover:text-foreground'"
              >
                {{ workspace.name }}
              </span>
              <span class="mt-1 inline-block border border-border/60 px-1 py-0.5 font-mono text-[8px] font-black uppercase tracking-tighter text-muted-foreground/80">
                {{ workspace.role || 'member' }}
              </span>
            </div>

            <div v-if="workspace.id === currentWorkspaceId" class="size-1.5 animate-pulse bg-primary/40" />
          </button>

          <div v-if="!workspacesLoading && sortedWorkspaces.length === 0" class="flex flex-col items-center gap-3 px-4 pt-10 text-center">
            <span class="font-mono text-[10px] italic text-muted-foreground/40">No workspaces yet</span>
          </div>
        </div>
      </aside>

      <template v-if="workspacesLoading">
        <div class="flex flex-1 items-center justify-center">
          <div class="flex flex-col items-center gap-4 opacity-40">
            <div class="size-12 animate-spin border-4 border-primary/20 border-t-primary" />
            <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Loading workspaces...</span>
          </div>
        </div>
      </template>

      <template v-else-if="workspacesError">
        <div class="flex flex-1 items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-destructive/30 bg-destructive/5 text-destructive/40">
              <PhWarning class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-destructive">Failed to Load</h3>
              <p class="mt-1 max-w-xs font-mono text-[10px] text-muted-foreground/60">{{ workspacesError }}</p>
            </div>
            <button
              class="flex h-10 items-center gap-2 border-2 border-destructive/30 bg-destructive/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-destructive transition-all hover:bg-destructive/15"
              type="button"
              @click="loadWorkspaces()"
            >
              Retry
            </button>
          </div>
        </div>
      </template>

      <template v-else-if="sortedWorkspaces.length === 0 || !selectedWorkspace">
        <div class="flex flex-1 items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-muted-foreground/30 bg-muted/10 text-muted-foreground/40">
              <PhBuildings class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-muted-foreground">No Workspaces Yet</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60">Create a workspace to unlock collections, environments, access grants, and execution history.</p>
            </div>
            <button
              class="flex h-10 items-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
              type="button"
              @click="createOpen = true"
            >
              <PhPlus class="size-4" />
              Create Workspace
            </button>
          </div>
        </div>
      </template>

      <template v-else>
        <section class="flex min-w-[520px] flex-1 flex-col overflow-hidden border bg-background">
          <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-3">
            <div class="flex min-w-0 items-center gap-2">
              <PhShieldCheck class="size-4 text-primary" />
              <span class="truncate font-mono text-[10px] font-black uppercase tracking-widest">
                Workspace Control
              </span>
            </div>
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/60">
              Active Boundary
            </span>
          </div>

          <div class="custom-scrollbar flex-1 overflow-y-auto p-6">
            <div class="mx-auto max-w-3xl space-y-4">
              <div class="border-2 border-primary/15 bg-primary/5 p-4 shadow-[4px_4px_0_0_rgba(16,185,129,0.08)]">
                <div class="flex flex-wrap items-start justify-between gap-4">
                  <div class="min-w-0">
                    <span class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/70">Selected workspace</span>
                    <h2 class="mt-1 truncate font-heading text-2xl font-black">{{ selectedWorkspace.name }}</h2>
                    <div class="mt-3 flex flex-wrap gap-2">
                      <span class="border border-primary/25 bg-background/60 px-2 py-1 font-mono text-[9px] font-black uppercase tracking-widest text-primary">
                        {{ selectedRole }}
                      </span>
                      <span class="border border-border/70 bg-background/60 px-2 py-1 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
                        {{ selectedMembers }} members
                      </span>
                    </div>
                  </div>

                  <button
                    class="flex h-9 items-center gap-2 border-2 border-primary/30 bg-primary/8 px-3 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
                    type="button"
                    @click="currentWorkspaceId = selectedWorkspace.id"
                  >
                    <PhCheck class="size-4" />
                    Active
                  </button>
                </div>
              </div>

              <div class="grid gap-3 md:grid-cols-2">
                <form
                  class="border bg-muted/10 p-4"
                  @submit.prevent="saveRename(selectedWorkspace)"
                >
                  <div class="mb-3 flex items-center gap-2">
                    <PhPencilSimple class="size-4 text-primary" />
                    <h3 class="font-mono text-[11px] font-black uppercase tracking-widest">Rename</h3>
                  </div>
                  <input
                    v-model="editingName"
                    class="h-11 w-full rounded-none border-2 border-primary/10 bg-background/50 px-3 font-mono text-sm outline-none transition-all placeholder:text-muted-foreground/70 hover:border-primary/40 hover:bg-background focus:border-primary"
                    :placeholder="selectedWorkspace.name"
                    type="text"
                    @focus="beginRename(selectedWorkspace)"
                  >
                  <button
                    class="mt-3 flex h-9 w-full items-center justify-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5 disabled:opacity-50"
                    :disabled="savingId === selectedWorkspace.id"
                    type="submit"
                  >
                    {{ savingId === selectedWorkspace.id ? 'Saving...' : 'Save name' }}
                  </button>
                </form>

                <div class="border border-destructive/25 bg-destructive/5 p-4">
                  <div class="mb-3 flex items-center gap-2">
                    <PhTrash class="size-4 text-destructive" />
                    <h3 class="font-mono text-[11px] font-black uppercase tracking-widest text-destructive">Danger</h3>
                  </div>
                  <p class="min-h-11 text-sm leading-5 text-muted-foreground">
                    Delete removes this workspace and all scoped collections, requests, environments, variables, memberships, and grants.
                  </p>
                  <button
                    class="mt-3 flex h-9 w-full items-center justify-center gap-2 border-2 border-destructive/30 bg-destructive/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-destructive transition-all hover:bg-destructive/15 hover:shadow-[3px_3px_0_0_rgba(239,68,68,0.12)] active:translate-x-0.5 active:translate-y-0.5"
                    type="button"
                    @click="requestDelete(selectedWorkspace)"
                  >
                    <PhTrash class="size-4" />
                    Delete workspace
                  </button>
                </div>
              </div>

              <p v-if="actionError" class="border border-destructive/30 bg-destructive/10 px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-widest text-destructive">
                {{ actionError }}
              </p>
            </div>
          </div>
        </section>

        <aside class="flex w-72 shrink-0 flex-col overflow-hidden border bg-background">
          <div class="flex h-10 shrink-0 items-center justify-between border-b bg-muted/30 px-3">
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Access Snapshot</span>
            <PhUsersThree class="size-4 text-primary" />
          </div>
          <div class="space-y-3 p-4">
            <div class="border bg-muted/15 p-3">
              <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Role</span>
              <p class="mt-1 font-mono text-lg font-black uppercase text-primary">{{ selectedRole }}</p>
            </div>
            <div class="border bg-muted/15 p-3">
              <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Members</span>
              <p class="mt-1 font-mono text-lg font-black uppercase text-primary">{{ selectedMembers }}</p>
            </div>
            <NuxtLink
              class="flex h-10 items-center justify-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
              to="/access"
            >
              <PhUsersThree class="size-4" />
              Manage access
            </NuxtLink>
          </div>
        </aside>
      </template>
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
