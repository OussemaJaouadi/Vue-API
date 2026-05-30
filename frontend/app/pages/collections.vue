<script setup lang="ts">
import { PhFolderOpen, PhPlus, PhWarning } from '@phosphor-icons/vue'
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
  importCollections,
  detectJsonImport,
} = useCollections()

const importInput = ref<HTMLInputElement | null>(null)
const importOpen = ref(false)
const importPayload = ref<unknown | null>(null)
const importing = ref(false)
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
    importPayload.value = null
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
    const parsed = JSON.parse(content)
    const detected = detectJsonImport(parsed)
    importPayload.value = detected.status === 'ready' ? parsed : null
    importPreview.value = {
      fileName: file.name,
      ...detected,
    }
    importOpen.value = true
  }
  catch (error) {
    importPayload.value = null
    importPreview.value = {
      fileName: file.name,
      format: 'Invalid JSON',
      status: 'error',
      summary: 'Could not parse file',
      details: ['Check that the file is valid JSON before importing.'],
    }
    importOpen.value = true
  }

  if (importInput.value) {
    importInput.value.value = ''
  }
}

const confirmImport = async () => {
  if (!importPayload.value || !importPreview.value || importing.value) return

  importing.value = true
  try {
    await importCollections(importPayload.value)
    await workbench.loadCollections()
    importOpen.value = false
    importPayload.value = null
    importPreview.value = null
  }
  catch (error: any) {
    importPreview.value = {
      ...importPreview.value,
      status: 'error',
      summary: error?.data?.error || error?.message || 'Import failed',
      details: ['The backend rejected this import payload.', 'Nothing was written to the selected workspace.'],
    }
  }
  finally {
    importing.value = false
  }
}
</script>

<template>
  <div class="flex flex-col h-[calc(100dvh-5.5rem)] border bg-card overflow-hidden">
    <CollectionsHeader
      :group-count="workbench.treeItems.value.length"
      :request-count="requestCount"
      :collection-names="workbench.treeItems.value.map(group => group.name)"
      :active-collection-name="activeCollectionName"
      @import="importInput?.click()"
      @export="exportCollections"
      @add-collection="workbench.addFolder()"
      @select-collection="selectCollection"
    />

    <input
      ref="importInput"
      accept=".json,.yaml,.yml,application/json"
      class="hidden"
      type="file"
      @change="previewImportFile"
    >

    <div class="flex-1 flex min-w-0 gap-3 overflow-x-auto p-3 bg-muted/5">
      <template v-if="workbench.collectionsLoading.value">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 opacity-40">
            <div class="size-12 border-4 border-primary/20 border-t-primary animate-spin" />
            <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Loading collections...</span>
          </div>
        </div>
      </template>

      <template v-else-if="workbench.collectionsError.value">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-destructive/30 bg-destructive/5 text-destructive/40">
              <PhWarning class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-destructive">Failed to Load</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60 max-w-xs">{{ workbench.collectionsError.value }}</p>
            </div>
            <button class="flex h-10 items-center gap-2 border-2 border-destructive/30 bg-destructive/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-destructive transition-all hover:bg-destructive/15" type="button" @click="workbench.loadCollections()">Retry</button>
          </div>
        </div>
      </template>

      <template v-else-if="workbench.treeItems.value.length === 0 && workbench.rootRequests.value.length === 0">
        <div class="flex-1 flex items-center justify-center">
          <div class="flex flex-col items-center gap-4 px-8 text-center">
            <div class="grid size-16 place-items-center border-2 border-dashed border-muted-foreground/30 bg-muted/10 text-muted-foreground/40">
              <PhFolderOpen class="size-8" />
            </div>
            <div>
              <h3 class="font-mono text-[13px] font-black uppercase tracking-tight text-muted-foreground">No Collections Yet</h3>
              <p class="mt-1 font-mono text-[10px] text-muted-foreground/60">Create a collection to organize your API requests.</p>
            </div>
            <button
              class="flex h-10 items-center gap-2 border-2 border-primary/30 bg-primary/8 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:bg-primary/15 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.12)] active:translate-x-0.5 active:translate-y-0.5"
              type="button"
              @click="workbench.addFolder()"
            >
              <PhPlus class="size-4" />
              Create Collection
            </button>
          </div>
        </div>
      </template>

      <template v-else>
      <CollectionsRoster
        :loading="false"
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
      </template>
    </div>

    <CollectionsImportSheet
      v-model:open="importOpen"
      :preview="importPreview"
      :importing="importing"
      @import="confirmImport"
    />
  </div>
</template>
