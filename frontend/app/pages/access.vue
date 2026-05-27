<script setup lang="ts">
import type { AccessLevel } from '~/types/access'
import AccessHeader from '~/components/access/AccessHeader.vue'
import AccessUserRoster from '~/components/access/AccessUserRoster.vue'
import AccessGrantEditor from '~/components/access/AccessGrantEditor.vue'
import AccessPolicyPanel from '~/components/access/AccessPolicyPanel.vue'
import AccessInviteSheet from '~/components/access/AccessInviteSheet.vue'

const {
  users,
  usersLoading,
  selectedUserId,
  selectedUser,
  grantSections,
  executionRows,
  deniedTargets,
  updateRole,
  kickUser,
  updateGrant,
  resolveDenied,
} = useAccess()

const inviteOpen = ref(false)

const roleOptions = [
  { value: 'manager', label: 'Manager' },
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
      <div v-if="usersLoading" class="flex-1 grid place-items-center">
        <span class="font-mono text-[10px] uppercase tracking-widest text-muted-foreground/60">Loading users...</span>
      </div>

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
