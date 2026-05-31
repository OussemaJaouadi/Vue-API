<script setup lang="ts">
import {
  PhCheckCircle,
  PhFileSearch,
  PhUploadSimple,
  PhWarning,
} from '@phosphor-icons/vue'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '~/components/ui/sheet'

interface ImportPreview {
  fileName: string
  format: string
  status: 'ready' | 'unsupported' | 'error'
  summary: string
  details: string[]
}

interface ImportResult {
  format: string
  collectionsCreated: number
  requestsCreated: number
  warnings: string[]
}

const props = defineProps<{
  preview: ImportPreview | null
  result?: ImportResult | null
  importing?: boolean
}>()

const open = defineModel<boolean>('open', { default: false })

defineEmits<{
  import: []
}>()

const statusIcon = computed(() => {
  if (props.result) return PhCheckCircle
  if (props.preview?.status === 'ready') return PhCheckCircle
  if (props.preview?.status === 'unsupported') return PhWarning
  return PhFileSearch
})

const statusColor = computed(() => {
  if (props.result) return 'text-emerald-500'
  if (props.preview?.status === 'ready') return 'text-emerald-500'
  if (props.preview?.status === 'unsupported') return 'text-amber-500'
  return 'text-destructive'
})
</script>

<template>
  <Sheet v-model:open="open">
    <SheetContent
      accessibility-title="Import Collection"
      accessibility-description="Preview and confirm collection import from external formats."
      class="w-[min(540px,100vw)] max-w-none border-l-2 border-primary/20 bg-background p-0 sm:max-w-none select-none"
    >
      <!-- Tactile Header -->
      <SheetHeader class="border-b bg-muted/30 p-6">
        <div class="flex items-center gap-5 pr-8">
          <div class="grid size-12 place-items-center border-2 border-primary shadow-[4px_4px_0_0_rgba(16,185,129,0.2)] bg-primary/10 text-primary transition-transform hover:scale-105">
            <PhUploadSimple class="size-6" />
          </div>
          <div class="min-w-0">
            <SheetTitle class="truncate font-heading text-xl font-black uppercase tracking-tight text-foreground">
              Import Collection
            </SheetTitle>
            <SheetDescription class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
              Migration / Data Ingestion
            </SheetDescription>
          </div>
        </div>
      </SheetHeader>

      <div v-if="preview" class="min-h-0 flex-1 overflow-y-auto custom-scrollbar">
        <!-- Source File -->
        <div class="border-b bg-muted/5 p-6 space-y-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Selected Source</span>
          <div class="flex items-center gap-3 border-2 border-primary/10 bg-background/50 px-3 py-3 font-mono text-sm shadow-sm">
            <PhFileSearch class="size-5 text-primary/60" />
            <span class="truncate font-bold">{{ preview.fileName }}</span>
          </div>
        </div>

        <!-- Format & Status -->
        <div class="border-b bg-muted/5 p-6">
          <div class="flex items-start gap-5">
            <div 
              class="grid size-10 place-items-center border-2 shrink-0 transition-all shadow-sm"
              :class="[
                preview.status === 'ready' ? 'border-emerald-500/30 bg-emerald-500/5 text-emerald-500' :
                preview.status === 'unsupported' ? 'border-amber-500/30 bg-amber-500/5 text-amber-500' :
                'border-destructive/30 bg-destructive/5 text-destructive'
              ]"
            >
              <component :is="statusIcon" class="size-5" />
            </div>
            
            <div class="min-w-0 flex-1">
              <h4 class="font-mono text-[11px] font-black uppercase tracking-widest text-foreground">{{ preview.format }}</h4>
              <p class="mt-1 font-mono text-[10px] font-bold text-muted-foreground">{{ preview.summary }}</p>
              
              <div class="mt-4 space-y-2 border-l-2 border-primary/10 pl-4">
                <div
                  v-for="(detail, i) in preview.details"
                  :key="i"
                  class="font-mono text-[10px] text-muted-foreground/80 leading-relaxed"
                >
                  • {{ detail }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Import Result -->
        <div v-if="result" class="border-b bg-emerald-500/5 p-6">
          <div class="border-2 border-emerald-500/25 bg-background/60 p-4 shadow-[4px_4px_0_0_rgba(16,185,129,0.1)]">
            <div class="flex items-center gap-2 text-emerald-600 dark:text-emerald-400">
              <PhCheckCircle class="size-4" />
              <span class="font-mono text-[10px] font-black uppercase tracking-widest">Import Complete</span>
            </div>

            <div class="mt-4 grid grid-cols-2 gap-3">
              <div class="border border-emerald-500/20 bg-emerald-500/5 p-3">
                <div class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Collections</div>
                <div class="mt-1 font-mono text-2xl font-black text-foreground">{{ result.collectionsCreated }}</div>
              </div>
              <div class="border border-emerald-500/20 bg-emerald-500/5 p-3">
                <div class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Requests</div>
                <div class="mt-1 font-mono text-2xl font-black text-foreground">{{ result.requestsCreated }}</div>
              </div>
            </div>

            <div v-if="result.warnings.length > 0" class="mt-4 border-l-2 border-amber-500/60 pl-3">
              <div
                v-for="(warning, i) in result.warnings"
                :key="i"
                class="font-mono text-[10px] leading-relaxed text-amber-700/80 dark:text-amber-300/80"
              >
                • {{ warning }}
              </div>
            </div>
          </div>
        </div>

        <!-- Warning for Unsupported -->
        <div v-if="preview.status === 'unsupported'" class="p-6 bg-amber-500/5">
          <div class="border-2 border-amber-500/20 p-4 space-y-2">
            <div class="flex items-center gap-2 text-amber-600 dark:text-amber-400">
              <PhWarning class="size-4" />
              <span class="font-mono text-[10px] font-black uppercase tracking-widest">Partial Compatibility</span>
            </div>
            <p class="font-mono text-[9px] leading-relaxed text-amber-700/60 dark:text-amber-400/60">
              Direct import for this format is currently being optimized. The backend parser will handle deep mapping of headers, auth, and complex body structures in the next release.
            </p>
          </div>
        </div>
      </div>

      <div v-else class="flex-1 grid place-items-center p-12 text-center">
        <div class="space-y-4 max-w-xs">
          <div class="grid size-16 mx-auto place-items-center border-2 border-dashed border-border/60 bg-muted/5 text-muted-foreground/20">
            <PhUploadSimple class="size-8" />
          </div>
          <p class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground/65 leading-relaxed">
            Select a JSON or YAML file to begin the validation process
          </p>
        </div>
      </div>

      <!-- Tactile Footer -->
      <SheetFooter class="border-t bg-muted/30 p-6">
        <div class="grid w-full gap-4 sm:grid-cols-2">
          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-primary/20 bg-primary/3 px-4 text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.15)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            :disabled="(!result && preview?.status !== 'ready') || importing"
            :class="((!result && preview?.status !== 'ready') || importing) && 'opacity-50 cursor-not-allowed'"
            type="button"
            @click="result ? open = false : $emit('import')"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">{{ importing ? 'Importing...' : result ? 'Done' : 'Confirm Ingestion' }}</span>
          </button>

          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-primary/10 bg-muted/20 px-4 text-muted-foreground transition-all hover:border-primary/40 hover:bg-background hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.1)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            type="button"
            @click="open = false"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Cancel</span>
          </button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
