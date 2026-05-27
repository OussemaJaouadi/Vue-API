<script setup lang="ts">
import EnvironmentsHeader from '~/components/environments/EnvironmentsHeader.vue'
import EnvironmentsRoster from '~/components/environments/EnvironmentsRoster.vue'
import EnvironmentsEditor from '~/components/environments/EnvironmentsEditor.vue'
import EnvironmentsPolicyPanel from '~/components/environments/EnvironmentsPolicyPanel.vue'
import EnvironmentsCreateSheet from '~/components/environments/EnvironmentsCreateSheet.vue'

const {
  environments,
  activeEnvironmentName,
  activeEnvironment,
  secretVariableCount,
  handleCreate,
  deleteEnvironment,
  addVariable,
} = useEnvironments()

const createOpen = ref(false)
</script>

<template>
  <div class="flex flex-col h-[calc(100dvh-5.5rem)] border bg-card overflow-hidden">
    <EnvironmentsHeader
      :active-environment-name="activeEnvironment.name"
      @add="createOpen = true"
    />

    <div class="flex-1 flex min-w-0 gap-3 overflow-x-auto p-3 bg-muted/5">
      <EnvironmentsRoster
        :environments="environments"
        :active-environment-name="activeEnvironmentName"
        @select="name => activeEnvironmentName = name"
      />

      <EnvironmentsEditor
        :environment-name="activeEnvironment.name"
        :variables="activeEnvironment.variables"
        @add-variable="addVariable"
      />

      <EnvironmentsPolicyPanel
        :environment="activeEnvironment"
        :secret-count="secretVariableCount"
        @delete="deleteEnvironment"
      />
    </div>

    <EnvironmentsCreateSheet
      v-model:open="createOpen"
      @create="handleCreate"
    />
  </div>
</template>
