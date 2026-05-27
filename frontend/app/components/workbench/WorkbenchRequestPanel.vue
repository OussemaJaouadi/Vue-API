<script setup lang="ts">
import {
  PhDatabase,
  PhKey,
  PhListChecks,
  PhPlus,
  PhShieldCheck,
  PhSlidersHorizontal,
} from '@phosphor-icons/vue'

const workbench = useWorkbench()
const activeTab = ref<'params' | 'headers' | 'body'>('params')
const authOpen = ref(false)
const authMode = useState<'inherit' | 'bearer' | 'api-key' | 'basic' | 'oauth2' | 'oidc' | 'none'>('workbench:auth-mode', () => 'bearer')
const activeParamCount = computed(() => workbench.queryParams.value.filter(param => param.enabled && param.key).length)
const activeHeaderCount = computed(() => workbench.headers.value.filter(header => header.enabled && header.key).length)
const bodyLabel = computed(() => workbench.isWebSocketRequest.value ? 'Message' : 'Body')

const addActiveItem = () => {
  if (activeTab.value === 'params') {
    workbench.addQueryParam()
  }
  if (activeTab.value === 'headers') {
    workbench.addHeader()
  }
}
</script>

<template>
  <section class="flex min-h-0 flex-col overflow-hidden border-b select-none">
    <UiTabs default-value="params" class="flex h-full min-h-0 flex-col overflow-hidden gap-0">
      <div class="flex h-10 shrink-0 items-center justify-between gap-2 border-b bg-muted/30 px-2">
        <UiTabsList class="h-full gap-0 bg-transparent p-0">
          <UiTabsTrigger value="params" class="group relative h-full flex-none rounded-none border-x border-transparent px-4 font-mono text-[10px] font-black uppercase tracking-widest data-[state=active]:bg-background data-[state=active]:text-primary transition-all duration-200" @click="activeTab = 'params'">
            <div class="flex items-center gap-2">
              <PhSlidersHorizontal class="size-3.5 transition-colors" :class="activeTab === 'params' ? 'text-primary' : 'text-muted-foreground/40 group-hover:text-primary/60'" />
              Params
              <span class="font-mono text-[9px] text-muted-foreground/30">{{ activeParamCount }}</span>
            </div>
            <div v-if="activeTab === 'params'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
          </UiTabsTrigger>

          <UiTabsTrigger value="headers" class="group relative h-full flex-none rounded-none border-x border-transparent px-4 font-mono text-[10px] font-black uppercase tracking-widest data-[state=active]:bg-background data-[state=active]:text-primary transition-all duration-200" @click="activeTab = 'headers'">
            <div class="flex items-center gap-2">
              <PhListChecks class="size-3.5 transition-colors" :class="activeTab === 'headers' ? 'text-primary' : 'text-muted-foreground/40 group-hover:text-primary/60'" />
              Headers
              <span class="font-mono text-[9px] text-muted-foreground/30">{{ activeHeaderCount }}</span>
            </div>
            <div v-if="activeTab === 'headers'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
          </UiTabsTrigger>

          <UiTabsTrigger value="body" class="group relative h-full flex-none rounded-none border-x border-transparent px-4 font-mono text-[10px] font-black uppercase tracking-widest data-[state=active]:bg-background data-[state=active]:text-primary transition-all duration-200" @click="activeTab = 'body'">
            <div class="flex items-center gap-2">
              <PhKey class="size-3.5 transition-colors" :class="activeTab === 'body' ? 'text-primary' : 'text-muted-foreground/40 group-hover:text-primary/60'" />
              {{ bodyLabel }}
            </div>
            <div v-if="activeTab === 'body'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary" />
          </UiTabsTrigger>

          <div class="h-5 w-px bg-border/40 mx-1" />

          <UiSheet v-model:open="authOpen">
            <UiSheetTrigger as-child>
              <button
                class="group relative h-full inline-flex flex-none items-center justify-center gap-2 rounded-none border-x border-transparent px-4 font-mono text-[10px] font-black uppercase tracking-widest transition-all duration-200 outline-none"
                :class="authMode === 'none' ? 'text-muted-foreground hover:bg-muted/40 hover:text-foreground' : 'text-primary bg-primary/[0.03]'"
                type="button"
              >
                <PhShieldCheck class="size-3.5 transition-colors" :class="authMode !== 'none' ? 'text-primary' : 'text-muted-foreground/40 group-hover:text-primary/60'" />
                Auth
                <span class="font-mono text-[9px] opacity-30">
                  {{ authMode === 'none' ? 'Empty' : authMode }}
                </span>
                <div v-if="authMode !== 'none'" class="absolute bottom-0 left-0 h-0.75 w-full bg-primary/40" />
              </button>
            </UiSheetTrigger>
            <UiSheetContent
              accessibility-title="Request authorization"
              accessibility-description="Configure request authorization and generated auth header."
              :show-close-button="false"
              side="bottom"
              class="!inset-auto !left-1/2 !top-1/2 !h-auto max-h-[82vh] !w-[min(600px,calc(100vw-32px))] !max-w-none -translate-x-1/2 -translate-y-1/2 border-2 border-primary/20 shadow-2xl p-0"
            >
              <div class="h-14 flex items-center px-4 border-b bg-muted/20 shrink-0">
                <h3 class="font-heading text-base font-black uppercase tracking-tight">Security context</h3>
              </div>
              <div class="p-6">
                <WorkbenchAuthPanel />
              </div>
              <div class="h-16 flex items-center justify-end px-4 border-t bg-muted/10 gap-3 shrink-0">
                <UiButton
                  variant="ghost"
                  class="rounded-none border-2 font-mono text-[10px] font-black uppercase tracking-widest"
                  @click="authOpen = false"
                >
                  Close
                </UiButton>
              </div>
            </UiSheetContent>
          </UiSheet>
        </UiTabsList>

        <div class="flex items-center gap-2">
          <button
            v-if="activeTab === 'params' || activeTab === 'headers'"
            class="group flex h-6 items-center gap-1.5 border-2 border-primary/20 bg-primary/5 px-2 font-mono text-[9px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 active:translate-x-0.5 active:translate-y-0.5"
            type="button"
            @click="addActiveItem"
          >
            <PhPlus class="size-3 transition-transform group-hover:scale-110" />
            Append
          </button>

          <div class="flex items-center gap-2 px-2 py-1 border-2 border-primary/15 bg-background font-mono text-[9px] font-black uppercase tracking-widest text-primary/60">
            <PhDatabase class="size-2.5" />
            Local
          </div>
        </div>
      </div>

      <UiTabsContent value="params" class="m-0 min-h-0 flex-1 overflow-hidden p-0">
        <WorkbenchParamsTable />
      </UiTabsContent>

      <UiTabsContent value="headers" class="m-0 min-h-0 flex-1 overflow-hidden p-0">
        <WorkbenchHeadersTable />
      </UiTabsContent>

      <UiTabsContent value="body" class="m-0 min-h-0 flex-1 overflow-hidden p-0">
        <WorkbenchWebSocketPanel v-if="workbench.isWebSocketRequest.value" />
        <WorkbenchCodeSurface
          v-else
          v-model="workbench.requestBody.value"
          v-model:language="workbench.requestBodyLanguage.value"
          class="min-w-0"
          editable
          label="Request body"
        />
      </UiTabsContent>
    </UiTabs>
  </section>
</template>
