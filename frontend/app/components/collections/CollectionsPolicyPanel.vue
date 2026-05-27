<script setup lang="ts">
import {
  PhGlobe,
  PhLightning,
  PhLockKey,
  PhShieldCheck,
  PhWarning,
} from '@phosphor-icons/vue'

interface Policy {
  defaultEnvironment: string
  allowedEnvironments: string[]
  visibility: 'project' | 'restricted'
  roles: string[]
}

defineProps<{
  activeCollectionName: string
  policy: Policy
}>()
</script>

<template>
  <aside class="flex flex-col w-[340px] gap-3 shrink-0 select-none overflow-hidden bg-card/10">
    <div class="flex h-10 items-center justify-between border-b bg-muted/30 px-3 shrink-0">
      <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Policy & Constraints</span>
      <span
        class="font-mono text-[8px] font-black uppercase tracking-widest px-1.5 py-0.5 border"
        :class="policy.visibility === 'restricted' ? 'border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400' : 'border-border text-muted-foreground/80'"
      >
        {{ policy.visibility }}
      </span>
    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar p-3 space-y-3">
      <!-- Visibility Alert -->
      <section v-if="policy.visibility === 'restricted'" class="border-2 border-amber-500/20 bg-amber-500/5 p-4 transition-all hover:bg-amber-500/10">
        <div class="flex gap-4">
          <PhWarning class="size-5 shrink-0 text-amber-500 animate-pulse" />
          <div class="space-y-1">
            <h4 class="font-mono text-[10px] font-black uppercase tracking-widest text-amber-700 dark:text-amber-300">Restricted Authority</h4>
            <p class="font-mono text-[9px] leading-relaxed text-amber-700/60 dark:text-amber-400/60">
              Usage in non-local environments requires explicit permission grants. Visibility is limited to authorized members only.
            </p>
          </div>
        </div>
      </section>

      <!-- Primary Policy Section -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhShieldCheck class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Access Controls</span>
        </div>

        <div class="p-4 space-y-5">
          <div class="space-y-2">
            <label class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80">Default environment</label>
            <div class="flex h-10 items-center gap-3 border-2 border-primary/10 bg-muted/5 px-3 transition-colors hover:border-primary/30">
              <PhLightning class="size-4 text-primary" />
              <span class="font-mono text-[10px] font-black uppercase tracking-widest text-foreground">{{ policy.defaultEnvironment }}</span>
            </div>
          </div>

          <div class="space-y-2">
            <label class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80">Target Environments</label>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="env in policy.allowedEnvironments"
                :key="env"
                class="border-2 border-primary/15 bg-primary/5 px-2.5 py-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-primary transition-all hover:scale-105"
              >
                {{ env }}
              </span>
            </div>
          </div>

          <div class="space-y-2">
            <label class="block font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80">Role Gates</label>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="role in policy.roles"
                :key="role"
                class="border-2 border-border/60 bg-muted px-2.5 py-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80 transition-all hover:border-primary/20"
              >
                {{ role }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- Safety Section -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhLockKey class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Execution Safety</span>
        </div>
        <div class="p-4 space-y-3">
          <div v-for="item in [
            { label: 'TLS Verification', value: 'Enabled', status: 'primary' },
            { label: 'Auto Redirects', value: 'Follow', status: 'muted' },
            { label: 'Network Timeout', value: '30s', status: 'muted' }
          ]" :key="item.label" class="flex items-center justify-between">
            <span class="font-mono text-[10px] font-bold uppercase tracking-widest text-muted-foreground">{{ item.label }}</span>
            <span class="font-mono text-[10px] font-black uppercase tracking-widest" :class="item.status === 'primary' ? 'text-primary' : 'text-foreground/60'">{{ item.value }}</span>
          </div>
        </div>
      </section>

      <!-- Global Headers -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhGlobe class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Global Headers</span>
        </div>
        <div class="p-4 space-y-3">
          <div v-for="item in [
            { label: 'Accept', value: 'application/json' },
            { label: 'User-Agent', value: 'API-Workbench/1.0' }
          ]" :key="item.label" class="flex items-center justify-between">
            <span class="font-mono text-[10px] font-bold uppercase tracking-widest text-muted-foreground">{{ item.label }}</span>
            <span class="font-mono text-[10px] font-black tracking-tight text-foreground/60">{{ item.value }}</span>
          </div>
        </div>
      </section>
    </div>
  </aside>
</template>
e>
