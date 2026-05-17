<script setup lang="ts">
import {
  PhClock,
  PhCommand,
  PhList,
  PhTerminalWindow,
  PhCaretDown,
  PhCube,
  PhBuildings,
} from '@phosphor-icons/vue'
import { useMagicKeys } from '@vueuse/core'

const props = defineProps<{
  workspaceName: string
  projectName: string
}>()

defineEmits<{
  openMobile: []
}>()

const now = ref(new Date())
let clockTimer: ReturnType<typeof setInterval> | undefined

const currentTime = computed(() => new Intl.DateTimeFormat(undefined, {
  hour: '2-digit',
  minute: '2-digit',
  second: '2-digit',
  hour12: false,
}).format(now.value))

const currentDate = computed(() => new Intl.DateTimeFormat(undefined, {
  month: 'short',
  day: '2-digit',
}).format(now.value))

const { meta, ctrl } = useMagicKeys()
const isMac = ref(false)

onMounted(() => {
  isMac.value = /Mac|iPod|iPhone|iPad/.test(navigator.userAgent)
  clockTimer = setInterval(() => {
    now.value = new Date()
  }, 1000)
})

const shortcutHint = computed(() => isMac.value ? '⌘ K' : 'CTRL K')

const { meta_k, ctrl_k } = useMagicKeys({
  passive: false,
  onEventFired(e) {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
    }
  },
})

watch([meta_k, ctrl_k], ([meta, ctrl]) => {
  if (meta || ctrl) {
    console.log('Search Palette Triggered')
  }
})

onBeforeUnmount(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>

<template>
  <header class="flex h-14 shrink-0 items-center justify-between border-b bg-background/95 px-3 backdrop-blur md:px-4">
    <div class="flex flex-1 items-center gap-3 overflow-hidden">
      <UiButton class="md:hidden rounded-none" size="icon-sm" variant="ghost" @click="$emit('openMobile')">
        <PhList class="size-5" />
        <span class="sr-only">Open navigation</span>
      </UiButton>

      <!-- Tactile Breadcrumb Selectors -->
      <nav class="flex min-w-0 items-center gap-0.5 font-mono text-[10px] font-black uppercase tracking-widest">
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <button class="group flex items-center gap-1.5 border border-transparent px-2 py-1 transition-all hover:border-primary/20 hover:bg-primary/5 active:bg-primary/10">
              <PhBuildings class="size-2.5 text-muted-foreground/40 group-hover:text-primary transition-colors" />
              <span class="truncate opacity-40 group-hover:opacity-100 transition-opacity">{{ workspaceName }}</span>
              <PhCaretDown class="size-2 opacity-20 group-hover:opacity-100" />
            </button>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="start" class="w-56 rounded-none border-2 p-1">
             <UiDropdownMenuLabel class="text-[9px] font-black uppercase tracking-[0.2em] text-primary/40 p-2">Switch Workspace</UiDropdownMenuLabel>
             <UiDropdownMenuItem class="rounded-none font-bold uppercase text-[10px]">{{ workspaceName }}</UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>

        <span class="text-border px-1 opacity-20">/</span>

        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <button class="group flex items-center gap-1.5 border border-transparent px-2 py-1 transition-all hover:border-primary/20 hover:bg-primary/5 active:bg-primary/10">
              <PhCube class="size-2.5 text-primary group-hover:scale-110 transition-transform" />
              <span class="truncate text-primary">{{ projectName }}</span>
              <PhCaretDown class="size-2 opacity-20 group-hover:opacity-100 text-primary" />
            </button>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="start" class="w-56 rounded-none border-2 p-1">
             <UiDropdownMenuLabel class="text-[9px] font-black uppercase tracking-[0.2em] text-primary/40 p-2">Switch Project</UiDropdownMenuLabel>
             <UiDropdownMenuItem class="rounded-none font-bold uppercase text-[10px]">{{ projectName }}</UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </nav>
    </div>

    <!-- Center Command Palette -->
    <div class="hidden flex-1 justify-center lg:flex px-4">
      <button class="group relative flex h-9 w-full max-w-md items-center justify-between gap-3 rounded-none border-2 border-primary/10 bg-muted/20 px-3 text-muted-foreground transition-all hover:border-primary/40 hover:bg-background hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.2)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none">
        <div class="flex items-center gap-2">
          <PhCommand class="size-3.5 group-hover:text-primary transition-colors" />
          <span class="font-mono text-[10px] font-bold uppercase tracking-[0.15em] group-hover:text-foreground transition-colors">Search Terminal...</span>
        </div>
        <div class="flex items-center gap-1 rounded-none border-2 border-primary/10 bg-background px-1.5 py-0.5 font-mono text-[9px] font-black shadow-xs group-hover:border-primary/30">
          {{ shortcutHint }}
        </div>
      </button>
    </div>

    <!-- Right Actions -->
    <div class="flex flex-1 items-center justify-end gap-3">
      <!-- Environment Switcher -->
      <UiDropdownMenu>
        <UiDropdownMenuTrigger as-child>
          <UiButton class="group h-9 rounded-none border-2 border-primary/20 bg-primary/3 px-3 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.1)]" variant="ghost">
            <PhTerminalWindow class="mr-2 size-3.5 group-hover:scale-110 transition-transform" />
            Local
            <PhCaretDown class="ml-2 size-2 opacity-40 group-hover:opacity-100" />
          </UiButton>
        </UiDropdownMenuTrigger>
        <UiDropdownMenuContent align="end" class="w-48 rounded-none border-2 p-1">
           <UiDropdownMenuLabel class="text-[9px] font-black uppercase tracking-[0.2em] text-primary/40 p-2">Active Environment</UiDropdownMenuLabel>
           <UiDropdownMenuItem class="rounded-none font-bold uppercase text-[10px]">Local</UiDropdownMenuItem>
           <UiDropdownMenuItem class="rounded-none font-bold uppercase text-[10px] opacity-40">Staging</UiDropdownMenuItem>
        </UiDropdownMenuContent>
      </UiDropdownMenu>

      <div class="text-border hidden h-4 w-px bg-border sm:block mx-1" />

      <!-- Date/Time Display -->
      <div class="hidden items-center gap-2 font-mono text-[11px] tabular-nums tracking-widest text-muted-foreground sm:flex">
        <PhClock class="size-3.5 opacity-50" />
        <span class="font-bold text-foreground/40">{{ currentDate }}</span>
        <span class="text-border opacity-30">/</span>
        <span class="font-bold text-foreground/80">{{ currentTime }}</span>
      </div>
    </div>
  </header>
</template>
