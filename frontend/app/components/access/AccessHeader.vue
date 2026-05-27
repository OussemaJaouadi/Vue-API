<script setup lang="ts">
import { PhPlus, PhUsersThree } from '@phosphor-icons/vue'
import { Button } from '~/components/ui/button'

defineEmits<{
  (e: 'invite'): void
}>()
</script>

<template>
  <header class="flex h-14 items-center justify-between border-b bg-muted/30 px-4 select-none">
    <div class="flex items-center gap-4">
      <div class="grid size-9 place-items-center border-2 border-primary shadow-[4px_4px_0_0_rgba(16,185,129,0.1)] bg-primary/10 text-primary transition-transform hover:scale-105">
        <PhUsersThree class="size-5" />
      </div>
      <div class="flex flex-col leading-none">
        <h1 class="font-heading text-lg font-black uppercase tracking-tight text-foreground">Access</h1>
        <span class="font-mono text-[9px] font-bold uppercase tracking-widest text-primary/60">
          User grants / permissions & policy
        </span>
      </div>
    </div>

    <button
      class="group relative h-9 rounded-none border-2 border-primary/20 bg-primary/3 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[3px_3px_0_0_rgba(16,185,129,0.1)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
      type="button"
      @click="$emit('invite')"
    >
      <div class="flex items-center gap-2.5">
        <PhPlus class="size-3.5 group-hover:rotate-90 transition-transform duration-500" />
        <span>Invite user</span>
      </div>
    </button>
  </header>
</template>

