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
  <section class="flex min-h-0 flex-col overflow-hidden border-b">
    <UiTabs default-value="params" class="flex h-full min-h-0 flex-col overflow-hidden gap-0">
      <div class="flex h-8 shrink-0 items-center justify-between gap-2 border-b bg-background px-2">
        <UiTabsList class="h-full gap-0 bg-transparent p-0">
          <UiTabsTrigger value="params" class="flex-none rounded-none border-x border-transparent px-3 font-mono text-[10px] font-black uppercase tracking-widest data-active:border-border data-active:bg-muted/35 data-active:text-primary" @click="activeTab = 'params'">
            <PhSlidersHorizontal class="size-3" />
            Params
            <span class="ml-1 text-[8px] text-muted-foreground">{{ activeParamCount }}</span>
          </UiTabsTrigger>
          <UiTabsTrigger value="headers" class="flex-none rounded-none border-x border-transparent px-3 font-mono text-[10px] font-black uppercase tracking-widest data-active:border-border data-active:bg-muted/35 data-active:text-primary" @click="activeTab = 'headers'">
            <PhListChecks class="size-3" />
            Headers
            <span class="ml-1 text-[8px] text-muted-foreground">{{ activeHeaderCount }}</span>
          </UiTabsTrigger>
          <UiTabsTrigger value="body" class="flex-none rounded-none border-x border-transparent px-3 font-mono text-[10px] font-black uppercase tracking-widest data-active:border-border data-active:bg-muted/35 data-active:text-primary" @click="activeTab = 'body'">
            <PhKey class="size-3" />
            {{ bodyLabel }}
          </UiTabsTrigger>

          <UiSheet v-model:open="authOpen">
            <UiSheetTrigger as-child>
              <button
                class="inline-flex h-full flex-none items-center justify-center gap-1.5 rounded-none border-x border-transparent px-3 font-mono text-[10px] font-black uppercase tracking-widest transition-colors"
                :class="authMode === 'none' ? 'text-muted-foreground hover:border-border hover:bg-muted/35 hover:text-foreground' : 'text-primary hover:border-border hover:bg-muted/35'"
                type="button"
              >
                <PhShieldCheck class="size-3" />
                Auth
                <span class="ml-1 text-[8px] text-muted-foreground">
                  {{ authMode === 'api-key' ? 'API Key' : authMode === 'oauth2' ? 'OAuth2' : authMode === 'oidc' ? 'OIDC' : authMode === 'bearer' ? 'Bearer' : authMode }}
                </span>
              </button>
            </UiSheetTrigger>
            <UiSheetContent
              accessibility-title="Request authorization"
              accessibility-description="Configure request authorization and generated auth header."
              :show-close-button="false"
              side="bottom"
              class="!inset-auto !left-1/2 !top-1/2 !h-auto max-h-[82vh] !w-[min(600px,calc(100vw-32px))] !max-w-none -translate-x-1/2 -translate-y-1/2 border shadow-2xl"
            >
              <WorkbenchAuthPanel />
            </UiSheetContent>
          </UiSheet>
        </UiTabsList>

        <div class="flex items-center gap-1.5">
          <button
            v-if="activeTab === 'params' || activeTab === 'headers'"
            class="flex h-5 items-center gap-1 border border-primary/25 bg-primary/8 px-2 font-mono text-[9px] font-black uppercase tracking-widest text-primary transition-colors hover:bg-primary/12"
            type="button"
            @click="addActiveItem"
          >
            <PhPlus class="size-3" />
            {{ activeTab === 'params' ? 'Param' : 'Header' }}
          </button>

          <UiBadge variant="outline" class="h-5 rounded-none px-1.5 font-mono text-[9px] opacity-60">
            <PhDatabase class="mr-1 size-2.5" />
            Local
          </UiBadge>
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
