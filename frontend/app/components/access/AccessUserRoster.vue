<script setup lang="ts">
import { PhUserCircle } from '@phosphor-icons/vue'
import { type AccessUser, roleTone } from '~/types/access'
import { Avatar, AvatarFallback } from '~/components/ui/avatar'

defineProps<{
  users: AccessUser[]
  selectedUserId: string
}>()

defineEmits<{
  (e: 'select', id: string): void
}>()
</script>

<template>
  <aside class="w-64 shrink-0 border-r bg-card/30 select-none flex flex-col overflow-hidden">
    <div class="flex h-10 items-center justify-between border-b bg-muted/30 px-3 shrink-0">
      <span class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">Roster</span>
      <span class="font-mono text-[9px] font-black uppercase tracking-widest text-primary/40">{{ users.length }} users</span>
    </div>

    <div class="flex-1 p-1 space-y-0.5 overflow-y-auto custom-scrollbar">
      <button
        v-for="user in users"
        :key="user.id"
        class="group relative flex h-14 w-full items-center gap-3 px-2 transition-all duration-200 outline-none shrink-0"
        :class="selectedUserId === user.id ? 'bg-primary/10 text-foreground' : 'text-muted-foreground/70 hover:bg-primary/5 hover:text-foreground'"
        type="button"
        @click="$emit('select', user.id)"
      >
        <div v-if="selectedUserId === user.id" class="wb-active-indicator" />
        
        <Avatar
          size="sm"
          class="size-8 rounded-none border transition-all duration-300 group-hover:scale-105 shrink-0"
          :class="selectedUserId === user.id ? 'border-primary/40 bg-primary/20 shadow-sm' : 'border-border/60 bg-muted/10 group-hover:border-primary/30'"
        >
          <AvatarFallback class="rounded-none bg-transparent">
            <PhUserCircle class="size-5 transition-colors" :class="selectedUserId === user.id ? 'text-primary' : 'text-muted-foreground/80 group-hover:text-primary/60'" />
          </AvatarFallback>
        </Avatar>

        <span class="min-w-0 flex-1 text-left">
          <span class="flex items-center gap-1.5">
            <span class="truncate font-mono text-[11px] font-black uppercase tracking-tight transition-colors" :class="selectedUserId === user.id ? 'text-primary' : 'text-foreground group-hover:text-foreground'">
              {{ user.username }}
            </span>
            <span
              class="size-1.5 shrink-0 rounded-full transition-transform duration-500 group-hover:scale-125"
              :class="user.status === 'active' ? 'bg-primary shadow-[0_0_8px_theme(colors.primary.DEFAULT)]' : 'bg-amber-500 shadow-[0_0_8px_theme(colors.amber.500)]'"
              :title="user.status"
            />
          </span>
          <span class="mt-0.5 flex items-center justify-between gap-2 min-w-0">
            <span class="truncate font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/90 group-hover:text-muted-foreground/70">{{ user.email }}</span>
            <span 
              class="shrink-0 px-1 py-0.5 border rounded-[2px] font-mono text-[8px] font-black uppercase tracking-tighter transition-all"
              :class="roleTone(user.role)"
            >
              {{ user.role }}
            </span>
          </span>
        </span>

        <div v-if="selectedUserId === user.id" class="absolute right-1 size-1.5 bg-primary/40 animate-pulse" />
      </button>
    </div>
  </aside>
</template>
