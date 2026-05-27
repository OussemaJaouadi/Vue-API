<script setup lang="ts">
import { PhDatabase, PhPlus } from '@phosphor-icons/vue'
import EnvironmentsHeader from '~/components/environments/EnvironmentsHeader.vue'
import EnvironmentsRoster from '~/components/environments/EnvironmentsRoster.vue'
import EnvironmentsEditor from '~/components/environments/EnvironmentsEditor.vue'
import EnvironmentsPolicyPanel from '~/components/environments/EnvironmentsPolicyPanel.vue'
import EnvironmentsCreateSheet from '~/components/environments/EnvironmentsCreateSheet.vue'

const {
  environments,
  envsLoading,
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
      :active-environment-name="activeEnvironment?.name ?? 'None'"
      @add="createOpen = true"
    />

    <div class="flex-1 flex min-w-0 gap-3 overflow-x-auto p-3 bg-muted/5">
      <EnvironmentsRoster
        :loading="envsLoading"
        :environments="environments"
        :active-environment-name="activeEnvironmentName"
        @select="name => activeEnvironmentName = name"
      />

      <template v-if="envsLoading">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 opacity-40">
            <div class="size-12 border-4 border-primary/20 border-t-primary animate-spin" />
            <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Loading environments...</span>
          </div>
        </div>
      </template>

      <template v-else-if="environments.length === 0 || !activeEnvironment">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-muted-foreground/30 bg-muted/10 text-muted-foreground/40">
              <PhDatabase class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-muted-foreground">No Environments Yet</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60">Create an environment to manage variables and secrets.</p>
            </div>
            <button
              class="flex h-10 items-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
              type="button"
              @click="createOpen = true"
            >
              <PhPlus class="size-4" />
              Create Environment
            </button>
          </div>
        </div>
      </template>

      <template v-else>
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
      </template>
    </div>

    <EnvironmentsCreateSheet
      v-model:open="createOpen"
      @create="handleCreate"
    />
  </div>
</template>
