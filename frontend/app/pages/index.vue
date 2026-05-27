<script setup lang="ts">
import WorkbenchSidebar from '~/components/workbench/WorkbenchSidebar.vue'
import WorkbenchEditor from '~/components/workbench/WorkbenchEditor.vue'
import WorkbenchResponse from '~/components/workbench/WorkbenchResponse.vue'

const workbench = useWorkbench()
</script>

<template>
  <div class="flex h-[calc(100dvh-5.5rem)] border bg-card overflow-hidden">
    <WorkbenchSidebar :style="{ width: `${workbench.sidebarWidth.value}px` }" />

    <WorkbenchResizer
      v-model="workbench.sidebarWidth.value"
      orientation="horizontal"
      :min="200"
      :max="450"
    />

    <main class="flex min-w-0 flex-1 flex-col overflow-hidden">
      <WorkbenchEditor>
        <template #response>
          <WorkbenchResponse />
        </template>
      </WorkbenchEditor>
    </main>
  </div>
</template>
