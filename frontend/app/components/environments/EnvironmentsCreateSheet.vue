<script setup lang="ts">
import { PhCaretDown, PhDatabase, PhGlobe, PhLockKey, PhShieldCheck } from '@phosphor-icons/vue'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '~/components/ui/sheet'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'create', environment: any): void
}>()

const draft = reactive({
  name: '',
  visibility: 'project' as 'project' | 'restricted',
  description: '',
})

const visibilityOptions = [
  { value: 'project', label: 'Project', description: 'Visible to all project members' },
  { value: 'restricted', label: 'Restricted', description: 'Requires explicit grants to view/use' },
]

const reset = () => {
  draft.name = ''
  draft.visibility = 'project'
  draft.description = ''
}

const handleCreate = () => {
  emit('create', { ...draft })
  emit('update:open', false)
  reset()
}
</script>

<template>
  <Sheet :open="open" @update:open="val => emit('update:open', val)">
    <SheetContent
      accessibility-title="Create new environment"
      accessibility-description="Initialize a new environment for variable resolution and secret storage."
      class="w-[min(540px,100vw)] max-w-none border-l-2 border-primary/20 bg-background p-0 sm:max-w-none select-none"
    >
      <!-- Tactile Header -->
      <SheetHeader class="border-b bg-muted/30 p-6">
        <div class="flex items-center gap-5 pr-8">
          <div class="grid size-12 place-items-center border-2 border-primary shadow-[4px_4px_0_0_rgba(16,185,129,0.2)] bg-primary/10 text-primary transition-transform hover:scale-105">
            <PhDatabase class="size-6" />
          </div>
          <div class="min-w-0">
            <SheetTitle class="truncate font-heading text-xl font-black uppercase tracking-tight text-foreground">
              New Environment
            </SheetTitle>
            <SheetDescription class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
              Infrastructure / Variable Scoping
            </SheetDescription>
          </div>
        </div>
      </SheetHeader>

      <div class="min-h-0 flex-1 overflow-y-auto custom-scrollbar">
        <!-- Name Input -->
        <div class="border-b bg-muted/5 p-6">
          <label class="grid gap-2">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Environment Name</span>
            <input
              v-model="draft.name"
              class="h-11 rounded-none border-2 border-primary/10 bg-background/50 px-3 font-mono text-sm outline-none transition-all placeholder:text-muted-foreground/70 hover:border-primary/40 hover:bg-background focus:border-primary shadow-[inset_0_1px_2px_rgba(0,0,0,0.05)]"
              placeholder="e.g. Staging, Production-Internal"
              type="text"
              autofocus
            >
          </label>
        </div>

        <!-- Visibility Selection -->
        <div class="border-b bg-muted/5 p-6 space-y-4">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Visibility Tier</span>
          
          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="opt in visibilityOptions"
              :key="opt.value"
              class="flex flex-col border-2 p-4 text-left transition-all duration-300 outline-none"
              :class="draft.visibility === opt.value 
                ? 'border-primary bg-primary/8 shadow-[4px_4px_0_0_rgba(16,185,129,0.15)]' 
                : 'border-border/40 bg-card/30 text-muted-foreground hover:border-primary/20'"
              type="button"
              @click="draft.visibility = opt.value as any"
            >
              <div class="flex items-center justify-between mb-2">
                <span class="font-mono text-[10px] font-black uppercase tracking-widest" :class="draft.visibility === opt.value ? 'text-primary' : ''">
                  {{ opt.label }}
                </span>
                <component :is="opt.value === 'restricted' ? PhLockKey : PhGlobe" class="size-4" :class="draft.visibility === opt.value ? 'text-primary' : 'opacity-70'" />
              </div>
              <p class="font-mono text-[9px] leading-relaxed transition-colors" :class="draft.visibility === opt.value ? 'text-primary/70' : 'text-muted-foreground/80'">
                {{ opt.description }}
              </p>
            </button>
          </div>
        </div>

        <!-- Role & Initial Access (Read-only Preview) -->
        <div class="p-6 bg-muted/5 space-y-4">
          <div class="flex items-center gap-2">
            <PhShieldCheck class="size-4 text-primary/60" />
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Inherited Baseline</span>
          </div>
          
          <div class="grid grid-cols-2 gap-3 opacity-60">
            <div class="border-2 border-border/40 bg-background/30 px-3 py-2">
              <p class="font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/80">Default Role Access</p>
              <p class="mt-1 font-mono text-[10px] font-black uppercase tracking-widest text-foreground">Developers + Managers</p>
            </div>
            <div class="border-2 border-border/40 bg-background/30 px-3 py-2">
              <p class="font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/80">Secret Access</p>
              <p class="mt-1 font-mono text-[10px] font-black uppercase tracking-widest text-amber-500">Restricted</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Tactile Footer -->
      <SheetFooter class="border-t bg-muted/30 p-6">
        <div class="grid w-full gap-4 sm:grid-cols-2">
          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-primary/20 bg-primary/3 px-4 text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.15)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            :disabled="!draft.name"
            :class="!draft.name && 'opacity-50 cursor-not-allowed'"
            type="button"
            @click="handleCreate"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Create Environment</span>
          </button>

          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-primary/10 bg-muted/20 px-4 text-muted-foreground transition-all hover:border-primary/40 hover:bg-background hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.1)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            type="button"
            @click="emit('update:open', false)"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Cancel</span>
          </button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
