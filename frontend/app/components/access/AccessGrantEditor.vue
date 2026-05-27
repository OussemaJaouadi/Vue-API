<script setup lang="ts">
import { PhCaretDown, PhUserMinus } from '@phosphor-icons/vue'
import { type AccessLevel, type AccessUser, type GrantTarget, accessTone, roleTone } from '~/types/access'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

const props = defineProps<{
  user: AccessUser
  roleOptions: Array<{ value: string, label: string }>
  sections: Array<{
    key: GrantTarget
    label: string
    icon: any
    rows: Array<{
      name: string
      meta: string
      level: AccessLevel
    }>
  }>
  accessOptions: Array<{ value: AccessLevel, label: string }>
}>()

const emit = defineEmits<{
  (e: 'updateRole', role: string): void
  (e: 'updateGrant', target: GrantTarget, name: string, level: AccessLevel): void
  (e: 'kickUser', id: string): void
}>()

const getOptionLabel = (value: AccessLevel) => props.accessOptions.find(opt => opt.value === value)?.label ?? value
const getRoleLabel = (value: string) => props.roleOptions.find(opt => opt.value === value)?.label ?? value
</script>

<template>
  <section class="min-w-[660px] flex-1 border-r bg-card/10 select-none flex flex-col overflow-hidden">
    <!-- Section Header -->
    <div class="grid border-b bg-muted/30 px-4 py-3 sm:grid-cols-[minmax(0,1fr)_340px] sm:items-center shrink-0">
      <div class="min-w-0">
        <h2 class="truncate font-mono text-[11px] font-black uppercase tracking-tight text-primary">{{ user.username }}</h2>
        <p class="font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/90">Access Control List / Permissions</p>
      </div>
      
      <div class="mt-3 flex items-center gap-3 sm:mt-0">
        <div class="flex flex-1 items-center gap-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Role</span>
          
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <button 
                class="flex h-8 w-full items-center justify-between border-2 px-2 font-mono text-[10px] font-black uppercase tracking-widest outline-none transition-all hover:shadow-md active:translate-x-0.5 active:translate-y-0.5"
                :class="roleTone(user.role)"
              >
                <span class="truncate">{{ getRoleLabel(user.role) }}</span>
                <PhCaretDown class="ml-2 size-3 opacity-70" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent class="min-w-40 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="end">
              <DropdownMenuItem
                v-for="role in roleOptions"
                :key="role.value"
                class="font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2 rounded-none border-l-2 border-transparent transition-all mb-0.5 last:mb-0"
                :class="[
                  user.role === role.value ? 'bg-primary/10 text-primary border-primary' : 'text-muted-foreground',
                  role.value === 'manager' ? 'focus:bg-indigo-500/10 focus:text-indigo-600 focus:border-indigo-500/60' :
                  role.value === 'developer' ? 'focus:bg-sky-500/10 focus:text-sky-600 focus:border-sky-500/60' :
                  'focus:bg-teal-500/10 focus:text-teal-600 focus:border-teal-500/60'
                ]"
                @select="emit('updateRole', role.value)"
              >
                {{ role.label }}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div class="h-6 w-px bg-border/40" />

        <button
          class="group flex size-8 items-center justify-center border-2 border-destructive/20 bg-destructive/5 text-destructive transition-all hover:border-destructive/50 hover:bg-destructive/10 hover:shadow-[3px_3px_0_0_rgba(239,68,68,0.15)] active:translate-x-0.5 active:translate-y-0.5"
          title="Kick User"
          type="button"
          @click="emit('kickUser', user.id)"
        >
          <PhUserMinus class="size-4 group-hover:scale-110 transition-transform" />
        </button>
      </div>
    </div>

    <div class="flex-1 divide-y divide-border/30 overflow-y-auto custom-scrollbar">
      <div
        v-for="section in sections"
        :key="section.key"
        class="transition-all duration-300"
      >
        <!-- Group Header -->
        <div class="grid h-10 grid-cols-[minmax(0,1fr)_128px] items-center border-b bg-muted/10 px-4 font-mono text-[10px] font-black uppercase tracking-widest">
          <span class="flex items-center gap-2 text-muted-foreground/70">
            <component :is="section.icon" class="size-3.5 text-primary/60" />
            {{ section.label }}
          </span>
          <span class="text-right text-muted-foreground/80">Authority</span>
        </div>

        <div class="divide-y divide-border/20">
          <div
            v-for="row in section.rows"
            :key="`${section.key}-${row.name}`"
            class="group grid min-h-14 grid-cols-[minmax(0,1fr)_128px] items-center gap-4 px-4 py-2 transition-all duration-200 hover:bg-primary/[0.02]"
          >
            <div class="min-w-0">
              <p class="truncate font-mono text-[11px] font-black uppercase tracking-tight text-foreground transition-colors group-hover:text-foreground">{{ row.name }}</p>
              <p class="truncate font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/80 transition-colors group-hover:text-muted-foreground">{{ row.meta }}</p>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <button 
                  class="flex h-8 w-full items-center justify-between border-2 px-2 text-center font-mono text-[9px] font-black uppercase tracking-widest outline-none transition-all duration-300 hover:shadow-md active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
                  :class="accessTone(row.level)"
                >
                  <span class="truncate flex-1">{{ getOptionLabel(row.level) }}</span>
                  <PhCaretDown class="ml-1 size-2.5 opacity-70 group-hover:opacity-100" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="min-w-32 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="end">
                <DropdownMenuItem
                  v-for="option in accessOptions"
                  :key="option.value"
                  class="font-mono text-[9px] font-black uppercase tracking-widest transition-all mb-0.5 last:mb-0 px-3 py-2.5 rounded-none border-l-2"
                  :class="[
                    row.level === option.value 
                      ? 'bg-primary/10 text-primary border-primary' 
                      : 'border-transparent text-muted-foreground',
                    option.value === 'admin' ? 'focus:bg-amber-500/10 focus:text-amber-600 focus:border-amber-500/60' : 
                    option.value === 'write' ? 'focus:bg-blue-500/10 focus:text-blue-600 focus:border-blue-500/60' :
                    option.value === 'read' ? 'focus:bg-emerald-500/10 focus:text-emerald-600 focus:border-emerald-500/60' :
                    'focus:bg-destructive/10 focus:text-destructive focus:border-destructive/60'
                  ]"
                  @select="emit('updateGrant', section.key, row.name, option.value)"
                >
                  {{ option.label }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
