<script setup lang="ts">
import { PhCheck, PhLockKey, PhShieldCheck } from '@phosphor-icons/vue'
import { type AccessLevel, type AccessUser, type DeniedTarget, type GrantTarget, accessTone } from '~/types/access'

const props = defineProps<{
  user: AccessUser
  executionRows: Array<{
    name: string
    meta: string
    level: AccessLevel
  }>
  deniedTargets: DeniedTarget[]
}>()

const emit = defineEmits<{
  (e: 'resolveDenied', target: GrantTarget, id: string): void
}>()
</script>

<template>
  <aside class="w-80 shrink-0 bg-card/30 select-none overflow-y-auto custom-scrollbar">
    <div class="flex h-10 items-center gap-2 border-b bg-muted/30 px-3 font-mono text-[10px] font-black uppercase tracking-widest text-foreground/70">
      <PhShieldCheck class="size-3.5 text-primary/60" />
      Policy preview
    </div>

    <div class="divide-y divide-border/30">
      <!-- Logic Description -->
      <div class="p-4">
        <div class="mb-2 flex items-center gap-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
          <PhLockKey class="size-3 text-primary/60" />
          Backend logic
        </div>
        <p class="font-mono text-[10px] font-bold leading-relaxed text-muted-foreground/70 uppercase tracking-tight">
          Execution requires combined collection & environment grants. Secret read is restricted to UI visibility only.
        </p>
      </div>

      <!-- Quick Stats -->
      <div class="grid grid-cols-2 divide-x divide-border/30">
        <div class="group p-4 transition-colors hover:bg-primary/[0.02]">
          <p class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/90 transition-colors group-hover:text-muted-foreground/70">Denied</p>
          <p class="mt-1 font-mono text-xl font-black transition-all duration-300 group-hover:scale-110" :class="deniedTargets.length > 0 ? 'text-destructive shadow-[0_0_15px_rgba(239,68,68,0.1)]' : 'text-primary shadow-[0_0_15px_rgba(16,185,129,0.1)]'">
            {{ deniedTargets.length }}
          </p>
        </div>
        <div class="group p-4 transition-colors hover:bg-primary/[0.02]">
          <p class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/90 transition-colors group-hover:text-muted-foreground/70">Status</p>
          <p class="mt-2 truncate font-mono text-[10px] font-black uppercase tracking-widest text-foreground group-hover:text-foreground">{{ user.status }}</p>
        </div>
      </div>

      <!-- Execution Grid -->
      <div class="p-4">
        <p class="mb-3 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Execution Matrix</p>
        <div class="grid gap-2">
          <div
            v-for="row in executionRows"
            :key="row.name"
            class="group grid grid-cols-[minmax(0,1fr)_80px] items-center gap-3 border border-border/40 bg-background/30 px-2 py-2 transition-all duration-300 hover:border-primary/20 hover:bg-background/50"
          >
            <span class="min-w-0">
              <span class="block truncate font-mono text-[9px] font-black uppercase tracking-tight text-foreground/70 group-hover:text-foreground">{{ row.name }}</span>
              <span class="block truncate font-mono text-[8px] font-bold uppercase tracking-widest text-muted-foreground/80 group-hover:text-muted-foreground">{{ row.meta }}</span>
            </span>
            <span
              class="border py-1 text-center font-mono text-[8px] font-black uppercase tracking-widest transition-all duration-300 group-hover:shadow-sm"
              :class="accessTone(row.level)"
            >
              {{ row.level === 'none' ? 'Blocked' : 'Allowed' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Denied Targets -->
      <div class="p-4">
        <p class="mb-3 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Access Gaps</p>
        <div v-if="deniedTargets.length" class="grid gap-2">
          <div
            v-for="target in deniedTargets"
            :key="`${target.target}-${target.id}`"
            class="group grid grid-cols-[minmax(0,1fr)_72px] gap-3 border border-dashed border-border/40 px-2 py-2 transition-all duration-300 hover:border-primary/30 hover:bg-primary/[0.01]"
          >
            <span class="min-w-0">
              <span class="block truncate font-mono text-[9px] font-black uppercase tracking-tight text-foreground/60 group-hover:text-foreground">{{ target.name }}</span>
              <span class="block truncate font-mono text-[8px] font-bold uppercase tracking-widest text-muted-foreground/80">{{ target.section }}</span>
            </span>
            <button
              class="h-7 border-2 border-primary/20 bg-primary/5 px-2 font-mono text-[8px] font-black uppercase tracking-widest text-primary transition-all duration-300 hover:bg-primary/20 hover:shadow-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
              type="button"
              @click="emit('resolveDenied', target.target, target.id)"
            >
              Grant
            </button>
          </div>
        </div>
        <div v-else class="flex items-center gap-3 border border-primary/20 bg-primary/5 px-3 py-3 font-mono text-[9px] font-black uppercase tracking-widest text-primary shadow-[inset_0_0_20px_rgba(16,185,129,0.02)]">
          <PhCheck class="size-3.5" />
          Compliant Policy
        </div>
      </div>
    </div>
  </aside>
</template>
e>
