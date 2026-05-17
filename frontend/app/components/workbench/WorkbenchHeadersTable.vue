<script setup lang="ts">
import {
  PhCheck,
  PhX,
} from '@phosphor-icons/vue'
import type { HeaderItem } from '~/composables/useWorkbench'

const workbench = useWorkbench()
const authMode = useState<'inherit' | 'bearer' | 'api-key' | 'basic' | 'oauth2' | 'oidc' | 'none'>('workbench:auth-mode', () => 'bearer')
const apiKeyName = useState<string>('workbench:auth-api-key-name', () => 'x-api-key')

const toggleHeader = (header: HeaderItem) => {
  header.enabled = !header.enabled
}

watch(workbench.headers.value, (headers) => {
  const lastHeader = headers[headers.length - 1]
  if (lastHeader && (lastHeader.key !== '' || lastHeader.value !== '')) {
    workbench.addHeader()
  }
}, { deep: true })

onMounted(() => {
  if (workbench.headers.value.length === 0) {
    workbench.addHeader()
  }
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden bg-card font-mono text-xs">
    <div class="grid shrink-0 grid-cols-[42px_minmax(140px,1fr)_minmax(160px,1.4fr)_36px] border-b bg-muted/20 py-1.5 text-[9px] font-black uppercase tracking-[0.18em] text-muted-foreground/60">
      <div class="flex items-center justify-center border-r border-primary/5">On</div>
      <span class="flex items-center border-r border-primary/5 px-3">Header</span>
      <span class="flex items-center border-r border-primary/5 px-3">Value</span>
      <span />
    </div>

    <div class="min-h-0 flex-1 overflow-auto">
      <div
        v-for="(header, index) in workbench.headers.value"
        :key="header.id"
        class="group grid min-h-9 grid-cols-[42px_minmax(140px,1fr)_minmax(160px,1.4fr)_36px] items-stretch border-b text-[11px] transition-colors"
        :class="header.enabled ? 'bg-background hover:bg-primary/3' : 'bg-muted/10 text-muted-foreground/55 hover:bg-muted/25'"
      >
        <div class="flex items-center justify-center border-r border-primary/5">
          <UiTooltip>
            <UiTooltipTrigger as-child>
              <button
                class="grid size-5 place-items-center border transition-all"
                :class="header.enabled ? 'border-primary bg-primary text-primary-foreground shadow-[0_0_10px_rgba(16,185,129,0.25)]' : 'border-muted-foreground/20 bg-background text-transparent hover:text-muted-foreground'"
                type="button"
                @click="toggleHeader(header)"
              >
                <PhCheck class="size-3" />
              </button>
            </UiTooltipTrigger>
            <UiTooltipContent side="top">{{ header.enabled ? 'Disable header' : 'Enable header' }}</UiTooltipContent>
          </UiTooltip>
        </div>

        <div class="flex items-center border-r border-primary/5 focus-within:bg-background">
          <UiBadge
            v-if="((header.key.toLowerCase() === 'authorization' && authMode !== 'none' && authMode !== 'api-key') || (authMode === 'api-key' && header.key.toLowerCase() === apiKeyName.toLowerCase()))"
            variant="outline"
            class="ml-2 h-5 shrink-0 rounded-none border-primary/25 bg-primary/8 px-1.5 font-mono text-[8px] uppercase text-primary"
          >
            auth
          </UiBadge>
          <input
            v-model="header.key"
            class="h-full w-full min-w-0 bg-transparent px-3 font-bold text-foreground/85 outline-none transition-colors placeholder:text-muted-foreground/25 focus:text-primary"
            placeholder="Header-Name"
            spellcheck="false"
          >
        </div>

        <div class="flex items-center border-r border-primary/5 focus-within:bg-background">
          <input
            v-model="header.value"
            class="h-full w-full bg-transparent px-3 text-muted-foreground outline-none transition-colors placeholder:text-muted-foreground/25 focus:text-foreground"
            placeholder="value or {{variable}}"
            spellcheck="false"
          >
        </div>

        <div class="flex items-center justify-center">
          <UiTooltip v-if="index < workbench.headers.value.length - 1 || header.key || header.value">
            <UiTooltipTrigger as-child>
              <button
                class="flex size-full items-center justify-center text-muted-foreground/20 transition-all hover:bg-destructive/10 hover:text-destructive active:scale-90"
                type="button"
                @click="workbench.removeHeader(header.id)"
              >
                <PhX class="size-3" />
              </button>
            </UiTooltipTrigger>
            <UiTooltipContent side="top">Remove header</UiTooltipContent>
          </UiTooltip>
        </div>
      </div>
    </div>
  </div>
</template>
