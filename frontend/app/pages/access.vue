<script setup lang="ts">
import { PhUsersThree, PhPlus, PhWarning } from '@phosphor-icons/vue'
import type { AccessLevel } from '~/types/access'
import AccessHeader from '~/components/access/AccessHeader.vue'
import AccessUserRoster from '~/components/access/AccessUserRoster.vue'
import AccessGrantEditor from '~/components/access/AccessGrantEditor.vue'
import AccessPolicyPanel from '~/components/access/AccessPolicyPanel.vue'
import AccessInviteSheet from '~/components/access/AccessInviteSheet.vue'

const {
  users,
  usersLoading,
  usersError,
  selectedUserId,
  selectedUser,
  grantSections,
  executionRows,
  deniedTargets,
  updateRole,
  kickUser,
  updateGrant,
  resolveDenied,
  loadUsers,
} = useAccess()

const inviteOpen = ref(false)

const roleOptions = [
  { value: 'admin', label: 'Admin' },
  { value: 'developer', label: 'Developer' },
  { value: 'tester', label: 'Tester' },
]

const accessOptions: { value: AccessLevel; label: string }[] = [
  { value: 'none', label: 'Denied' },
  { value: 'read', label: 'Read' },
  { value: 'write', label: 'Write' },
  { value: 'admin', label: 'Admin' },
]
</script>

<template>
  <div class="flex flex-col h-[calc(100dvh-5.5rem)] border bg-card overflow-hidden">
    <AccessHeader @invite="inviteOpen = true" />

    <div class="flex-1 flex min-w-0 gap-3 overflow-x-auto p-3 bg-muted/5">
      <template v-if="usersLoading">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 opacity-40">
            <div class="size-12 border-4 border-primary/20 border-t-primary animate-spin" />
            <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Loading users...</span>
          </div>
        </div>
      </template>

      <template v-else-if="usersError">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-destructive/30 bg-destructive/5 text-destructive/40">
              <PhWarning class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-destructive">Failed to Load</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60 max-w-xs">{{ usersError }}</p>
            </div>
            <button class="flex h-10 items-center gap-2 border-2 border-destructive/30 bg-destructive/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-destructive transition-all hover:bg-destructive/15" type="button" @click="loadUsers()">Retry</button>
          </div>
        </div>
      </template>

      <template v-else-if="users.length === 0">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-muted-foreground/30 bg-muted/10 text-muted-foreground/40">
              <PhUsersThree class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-muted-foreground">No Members Yet</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60">Invite members to manage access and permissions.</p>
            </div>
            <button
              class="flex h-10 items-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
              type="button"
              @click="inviteOpen = true"
            >
              <PhPlus class="size-4" />
              Invite Member
            </button>
          </div>
        </div>
      </template>

      <template v-else>
      <AccessUserRoster
        :users="users"
        :selected-user-id="selectedUserId"
        @select="id => selectedUserId = id"
        class="h-full"
      />

      <AccessGrantEditor
        v-if="selectedUser"
        :user="selectedUser"
        :role-options="roleOptions"
        :sections="grantSections"
        :access-options="accessOptions"
        @update-role="updateRole"
        @update-grant="updateGrant"
        @kick-user="kickUser"
        class="h-full"
      />
      <div v-else class="flex-1 border bg-muted/5 grid place-items-center font-mono text-[10px] uppercase tracking-widest text-muted-foreground/80">
        No user selected
      </div>

      <AccessPolicyPanel
        v-if="selectedUser"
        :user="selectedUser"
        :execution-rows="executionRows"
        :denied-targets="deniedTargets"
        @resolve-denied="resolveDenied"
        class="h-full"
      />
      </template>
    </div>

    <AccessInviteSheet
      v-model:open="inviteOpen"
      :role-options="roleOptions"
      :access-options="accessOptions"
    />
  </div>
</template>
