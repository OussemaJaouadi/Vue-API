<script setup lang="ts">
import CollectionsHeader from '~/components/collections/CollectionsHeader.vue'
import CollectionsRoster from '~/components/collections/CollectionsRoster.vue'
import CollectionsWorkbench from '~/components/collections/CollectionsWorkbench.vue'
import CollectionsPolicyPanel from '~/components/collections/CollectionsPolicyPanel.vue'
import CollectionsImportSheet from '~/components/collections/CollectionsImportSheet.vue'

const workbench = useWorkbench()
const {
  requestCount,
  activeCollectionName,
  expandedCollections,
  activeCollection,
  displayedCollections,
  activeRequestCount,
  environmentPolicyFor,
  selectCollection,
  toggleCollection,
  exportCollections,
  detectJsonImport,
} = useCollections()

const importInput = ref<HTMLInputElement | null>(null)
const importOpen = ref(false)
const importPreview = ref<{
  fileName: string
  format: string
  status: 'ready' | 'unsupported' | 'error'
  summary: string
  details: string[]
} | null>(null)

const previewImportFile = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const extension = file.name.split('.').pop()?.toLowerCase()
  if (extension === 'yaml' || extension === 'yml') {
    importPreview.value = {
      fileName: file.name,
      format: 'YAML spec',
      status: 'unsupported',
      summary: 'YAML upload selected',
      details: ['YAML parsing belongs in the backend parser service.', 'The UI accepts the file type so the flow is visible now.'],
    }
    importOpen.value = true
    return
  }

  try {
    const content = await file.text()
    const detected = detectJsonImport(JSON.parse(content))
    importPreview.value = {
      fileName: file.name,
      ...detected,
    }
    importOpen.value = true
  }
  catch (error) {
    importPreview.value = {
      fileName: file.name,
      format: 'Invalid JSON',
      status: 'error',
      summary: 'Could not parse file',
      details: ['Check that the file is valid JSON before importing.'],
    }
    importOpen.value = true
  }
}

const confirmImport = () => {
  // Logic to actually import would go here
  importOpen.value = false
  importPreview.value = null
}
</script>

<template>
  <div class="flex flex-col h-[calc(100dvh-5.5rem)] border bg-card overflow-hidden">
    <CollectionsHeader
      :group-count="workbench.treeItems.value.length"
      :request-count="requestCount"
      @import="importInput?.click()"
      @export="exportCollections"
      @add-collection="workbench.addFolder()"
    />

    <input
      ref="importInput"
      accept=".json,.yaml,.yml,application/json"
      class="hidden"
      type="file"
      @change="previewImportFile"
    >

    <div class="flex-1 flex min-w-0 gap-3 overflow-x-auto p-3 bg-muted/5">
      <CollectionsRoster
        :groups="workbench.treeItems.value"
        :active-collection-name="activeCollectionName"
        :total-request-count="requestCount"
        :policy-resolver="environmentPolicyFor"
        @select="selectCollection"
      />

      <CollectionsWorkbench
        :active-collection="activeCollection"
        :active-request-count="activeRequestCount"
        :displayed-collections="displayedCollections"
        :expanded-collections="expandedCollections"
        @add-request="workbench.addRequest"
        @toggle-expand="toggleCollection"
        @delete-request="id => workbench.deleteItem(id, false)"
        @delete-collection="name => workbench.deleteItem(name, true)"
      />

      <CollectionsPolicyPanel
        :active-collection-name="activeCollectionName"
        :policy="environmentPolicyFor(activeCollectionName)"
      />
    </div>

    <CollectionsImportSheet
      v-model:open="importOpen"
      :preview="importPreview"
      @import="confirmImport"
    />
  </div>
</template>
