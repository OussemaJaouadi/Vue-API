<script setup lang="ts">
import {
  PhGlobe,
  PhLockKey,
  PhShieldCheck,
  PhUsersThree,
  PhTrash,
  PhWarning,
} from '@phosphor-icons/vue'

interface Environment {
  name: string
  visibility: string
  allowedRoles: string[]
  variables: any[]
}

defineProps<{
  environment: Environment
  secretCount: number
}>()

defineEmits<{
  delete: []
}>()
</script>

<template>
  <aside class="flex flex-col w-[300px] gap-3 shrink-0 select-none overflow-hidden bg-card/10">
    <div class="flex h-10 items-center justify-between border-b bg-muted/30 px-3 shrink-0">
      <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Node Profile</span>
      <span
        class="font-mono text-[8px] font-black uppercase tracking-widest px-1.5 py-0.5 border"
        :class="environment.visibility === 'restricted' ? 'border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400' : 'border-border text-muted-foreground/80'"
      >
        {{ environment.visibility }}
      </span>
    </div>

    <div class="flex-1 overflow-y-auto custom-scrollbar p-3 space-y-3">
      <!-- Visibility Section -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhShieldCheck class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Visibility Tier</span>
        </div>
        <div class="p-4 space-y-4">
          <div
            class="flex h-10 items-center gap-3 border-2 px-3 transition-all"
            :class="environment.visibility === 'restricted' ? 'border-amber-500/20 bg-amber-500/5' : 'border-primary/10 bg-muted/5'"
          >
            <component :is="environment.visibility === 'restricted' ? PhLockKey : PhGlobe" class="size-4" :class="environment.visibility === 'restricted' ? 'text-amber-500' : 'text-primary'" />
            <span
              class="font-mono text-[10px] font-black uppercase tracking-widest"
              :class="environment.visibility === 'restricted' ? 'text-amber-600 dark:text-amber-400' : 'text-foreground'"
            >
              {{ environment.visibility }}
            </span>
          </div>
          <p class="font-mono text-[9px] leading-relaxed text-muted-foreground italic">
            <span v-if="environment.visibility === 'restricted'">
              Access to these variables and their secrets requires explicit permission grants.
            </span>
            <span v-else>
              Visible to all project members by default.
            </span>
          </p>
        </div>
      </section>

      <!-- Roles Section -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhUsersThree class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Allowed Roles</span>
        </div>
        <div class="p-4">
          <div class="flex flex-wrap gap-2">
            <span
              v-for="role in environment.allowedRoles"
              :key="role"
              class="border-2 border-border/60 bg-muted px-2.5 py-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground/80 transition-all hover:border-primary/20"
            >
              {{ role }}
            </span>
          </div>
        </div>
      </section>

      <!-- Secrets Section -->
      <section class="border bg-background shadow-sm overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b bg-muted/10 px-3">
          <PhLockKey class="size-3.5 text-primary/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest">Authority Monitor</span>
        </div>
        <div class="p-4 space-y-3">
          <div class="flex items-center justify-between">
            <span class="font-mono text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Active Secrets</span>
            <span 
              class="font-mono text-[11px] font-black uppercase tracking-widest px-2 py-0.5 border-2"
              :class="secretCount > 0 ? 'border-amber-500/20 bg-amber-500/5 text-amber-500' : 'border-border/40 text-muted-foreground/50'"
            >
              {{ secretCount }}
            </span>
          </div>
          <div class="mt-2 border-t-2 border-dashed border-border/20 pt-3">
            <p class="font-mono text-[8px] leading-relaxed text-muted-foreground/80 uppercase tracking-widest">
              Sensitive data is encrypted at rest. Decryption is gated by project authority levels.
            </p>
          </div>
        </div>
      </section>

      <!-- Danger Zone -->
      <section class="mt-6 border-2 border-destructive/20 bg-destructive/5 overflow-hidden">
        <div class="flex h-10 items-center gap-2 border-b border-destructive/10 bg-destructive/5 px-3">
          <PhWarning class="size-3.5 text-destructive/60" />
          <span class="font-mono text-[10px] font-black uppercase tracking-widest text-destructive/80">Danger Zone</span>
        </div>
        <div class="p-4">
          <p class="mb-4 font-mono text-[9px] leading-relaxed text-destructive/60">
            Deleting this environment will permanently remove all associated variables and secret configurations. This action cannot be reversed.
          </p>
          <button
            class="group relative flex h-10 w-full items-center justify-center gap-2 border-2 border-destructive/20 bg-background px-4 text-destructive transition-all hover:border-destructive/50 hover:bg-destructive/5 active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            type="button"
            @click="$emit('delete')"
          >
            <PhTrash class="size-3.5 group-hover:scale-110 transition-transform" />
            <span class="font-mono text-[9px] font-black uppercase tracking-widest">Destroy Node</span>
          </button>
        </div>
      </section>
    </div>
  </aside>
</template>
