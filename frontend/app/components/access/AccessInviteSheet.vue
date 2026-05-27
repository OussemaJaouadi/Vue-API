<script setup lang="ts">
import { PhCaretDown, PhDatabase, PhFolderOpen, PhKey, PhUsersThree } from '@phosphor-icons/vue'
import { type AccessLevel, accessTone, accessWeight } from '~/types/access'
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
  roleOptions: Array<{ value: string, label: string }>
  accessOptions: Array<{ value: AccessLevel, label: string }>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'send', invite: any): void
}>()

const inviteDraft = reactive({
  email: '',
  role: 'developer',
  collectionAccess: 'read' as AccessLevel,
  environmentAccess: 'read' as AccessLevel,
  secretAccess: 'none' as AccessLevel,
})

const inviteExecutionAllowed = computed(() => (accessWeight[inviteDraft.collectionAccess] ?? 0) > 0 && (accessWeight[inviteDraft.environmentAccess] ?? 0) > 0)

const getOptionLabel = (value: AccessLevel) => props.accessOptions.find(opt => opt.value === value)?.label ?? value
</script>

<template>
  <Sheet :open="open" @update:open="val => emit('update:open', val)">
    <SheetContent
      accessibility-title="Invite project user"
      accessibility-description="Prepare a project invite with role and initial resource access."
      class="w-[min(580px,100vw)] max-w-none border-l-2 border-primary/20 bg-background p-0 sm:max-w-none select-none"
    >
      <!-- Tactile Header -->
      <SheetHeader class="border-b bg-muted/30 p-6">
        <div class="flex items-center gap-5 pr-8">
          <div class="grid size-12 place-items-center border-2 border-primary shadow-[4px_4px_0_0_rgba(16,185,129,0.2)] bg-primary/10 text-primary transition-transform hover:scale-105">
            <PhUsersThree class="size-6" />
          </div>
          <div class="min-w-0">
            <SheetTitle class="truncate font-heading text-xl font-black uppercase tracking-tight text-foreground">
              Invite user
            </SheetTitle>
            <SheetDescription class="font-mono text-[10px] font-black uppercase tracking-widest text-primary/60">
              Provisioning / access initialization
            </SheetDescription>
          </div>
        </div>
      </SheetHeader>

      <div class="min-h-0 flex-1 overflow-y-auto custom-scrollbar">
        <!-- Email Input -->
        <div class="border-b bg-muted/5 p-6">
          <label class="grid gap-2">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Target Email</span>
            <input
              v-model="inviteDraft.email"
              class="h-11 rounded-none border-2 border-primary/10 bg-background/50 px-3 font-mono text-sm outline-none transition-all placeholder:text-muted-foreground/70 hover:border-primary/40 hover:bg-background focus:border-primary shadow-[inset_0_1px_2px_rgba(0,0,0,0.05)]"
              placeholder="teammate@example.com"
              type="email"
            >
          </label>
        </div>

        <!-- Role & Stats -->
        <div class="grid border-b md:grid-cols-2">
          <div class="grid gap-2 border-b p-6 md:border-b-0 md:border-r border-border/30 bg-muted/5">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">System Role</span>
            
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <button class="flex h-11 w-full items-center justify-between border-2 border-primary/10 bg-background/50 px-3 font-mono text-[10px] font-black uppercase tracking-widest text-foreground outline-none transition-all hover:border-primary/40 hover:bg-background focus:border-primary shadow-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none">
                  <span class="truncate">{{ roleOptions.find(r => r.value === inviteDraft.role)?.label }}</span>
                  <PhCaretDown class="ml-2 size-3 text-muted-foreground/80" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="min-w-48 rounded-none border-2 border-primary/20 bg-background p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="start">
                <DropdownMenuItem
                  v-for="role in roleOptions"
                  :key="role.value"
                  class="font-mono text-[10px] font-black uppercase tracking-widest px-3 py-2.5 rounded-none border-l-2 border-transparent focus:bg-primary/10 focus:text-primary focus:border-primary/60 transition-all mb-0.5 last:mb-0"
                  @select="inviteDraft.role = role.value"
                >
                  {{ role.label }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div class="p-6 bg-muted/5">
            <p class="mb-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Initial state</p>
            <div class="grid grid-cols-2 gap-3">
              <div class="border-2 border-border/40 bg-background/30 px-3 py-2 transition-all hover:border-primary/20">
                <p class="font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/80">Status</p>
                <p class="mt-1 font-mono text-[10px] font-black uppercase tracking-widest text-amber-500 shadow-[0_0_8px_rgba(245,158,11,0.1)]">Draft</p>
              </div>
              <div class="border-2 border-border/40 bg-background/30 px-3 py-2 transition-all hover:border-primary/20">
                <p class="font-mono text-[8px] font-black uppercase tracking-widest text-muted-foreground/80">Execution</p>
                <p class="mt-1 font-mono text-[10px] font-black uppercase tracking-widest transition-colors" :class="inviteExecutionAllowed ? 'text-primary' : 'text-destructive'">
                  {{ inviteExecutionAllowed ? 'Allowed' : 'Blocked' }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Permissions Section -->
        <div class="divide-y divide-border/30">
          <div class="p-6 bg-muted/10 font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground/70">
            Baseline Permissions
          </div>

          <div class="group grid grid-cols-[minmax(0,1fr)_180px] items-center gap-4 p-6 transition-all duration-300 hover:bg-primary/[0.02]">
            <div class="min-w-0">
              <div class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest text-foreground group-hover:text-foreground">
                <PhFolderOpen class="size-4 text-primary/60" />
                Collections
              </div>
              <p class="mt-1 font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/80 transition-colors group-hover:text-muted-foreground">Initial access applied to all nodes.</p>
            </div>
            
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <button 
                  class="flex h-10 w-full items-center justify-between border-2 px-3 text-center font-mono text-[9px] font-black uppercase tracking-widest outline-none transition-all duration-300 hover:shadow-md active:translate-x-0.5 active:translate-y-0.5"
                  :class="accessTone(inviteDraft.collectionAccess)"
                >
                  <span class="truncate flex-1">{{ getOptionLabel(inviteDraft.collectionAccess) }}</span>
                  <PhCaretDown class="ml-1 size-3 opacity-70 group-hover:opacity-60" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="min-w-44 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="end">
                <DropdownMenuItem
                  v-for="option in accessOptions"
                  :key="option.value"
                  class="font-mono text-[9px] font-black uppercase tracking-widest transition-all mb-0.5 last:mb-0 px-3 py-2.5 rounded-none border-l-2"
                  :class="[
                    inviteDraft.collectionAccess === option.value 
                      ? 'bg-primary/10 text-primary border-primary' 
                      : 'border-transparent text-muted-foreground',
                    option.value === 'admin' ? 'focus:bg-amber-500/10 focus:text-amber-600 focus:border-amber-500/60' : 
                    option.value === 'write' ? 'focus:bg-blue-500/10 focus:text-blue-600 focus:border-blue-500/60' :
                    option.value === 'read' ? 'focus:bg-emerald-500/10 focus:text-emerald-600 focus:border-emerald-500/60' :
                    'focus:bg-destructive/10 focus:text-destructive focus:border-destructive/60'
                  ]"
                  @select="inviteDraft.collectionAccess = option.value"
                >
                  {{ option.label }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div class="group grid grid-cols-[minmax(0,1fr)_180px] items-center gap-4 p-6 transition-all duration-300 hover:bg-primary/[0.02]">
            <div class="min-w-0">
              <div class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest text-foreground group-hover:text-foreground">
                <PhDatabase class="size-4 text-primary/60" />
                Environments
              </div>
              <p class="mt-1 font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/80 transition-colors group-hover:text-muted-foreground">Controls variable injection scope.</p>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <button 
                  class="flex h-10 w-full items-center justify-between border-2 px-3 text-center font-mono text-[9px] font-black uppercase tracking-widest outline-none transition-all duration-300 hover:shadow-md active:translate-x-0.5 active:translate-y-0.5"
                  :class="accessTone(inviteDraft.environmentAccess)"
                >
                  <span class="truncate flex-1">{{ getOptionLabel(inviteDraft.environmentAccess) }}</span>
                  <PhCaretDown class="ml-1 size-3 opacity-70 group-hover:opacity-60" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="min-w-44 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="end">
                <DropdownMenuItem
                  v-for="option in accessOptions"
                  :key="option.value"
                  class="font-mono text-[9px] font-black uppercase tracking-widest transition-all mb-0.5 last:mb-0 px-3 py-2.5 rounded-none border-l-2"
                  :class="[
                    inviteDraft.environmentAccess === option.value 
                      ? 'bg-primary/10 text-primary border-primary' 
                      : 'border-transparent text-muted-foreground',
                    option.value === 'admin' ? 'focus:bg-amber-500/10 focus:text-amber-600 focus:border-amber-500/60' : 
                    option.value === 'write' ? 'focus:bg-blue-500/10 focus:text-blue-600 focus:border-blue-500/60' :
                    option.value === 'read' ? 'focus:bg-emerald-500/10 focus:text-emerald-600 focus:border-emerald-500/60' :
                    'focus:bg-destructive/10 focus:text-destructive focus:border-destructive/60'
                  ]"
                  @select="inviteDraft.environmentAccess = option.value"
                >
                  {{ option.label }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div class="group grid grid-cols-[minmax(0,1fr)_180px] items-center gap-4 p-6 transition-all duration-300 hover:bg-primary/[0.02]">
            <div class="min-w-0">
              <div class="flex items-center gap-2 font-mono text-[10px] font-black uppercase tracking-widest text-foreground group-hover:text-foreground">
                <PhKey class="size-4 text-primary/60" />
                Secrets
              </div>
              <p class="mt-1 font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/80 transition-colors group-hover:text-muted-foreground">UI visibility of masked variables.</p>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <button 
                  class="flex h-10 w-full items-center justify-between border-2 px-3 text-center font-mono text-[9px] font-black uppercase tracking-widest outline-none transition-all duration-300 hover:shadow-md active:translate-x-0.5 active:translate-y-0.5"
                  :class="accessTone(inviteDraft.secretAccess)"
                >
                  <span class="truncate flex-1">{{ getOptionLabel(inviteDraft.secretAccess) }}</span>
                  <PhCaretDown class="ml-1 size-3 opacity-70 group-hover:opacity-60" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="min-w-44 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[6px_6px_0_0_rgba(16,185,129,0.1)]" align="end">
                <DropdownMenuItem
                  v-for="option in accessOptions"
                  :key="option.value"
                  class="font-mono text-[9px] font-black uppercase tracking-widest transition-all mb-0.5 last:mb-0 px-3 py-2.5 rounded-none border-l-2"
                  :class="[
                    inviteDraft.secretAccess === option.value 
                      ? 'bg-primary/10 text-primary border-primary' 
                      : 'border-transparent text-muted-foreground',
                    option.value === 'admin' ? 'focus:bg-amber-500/10 focus:text-amber-600 focus:border-amber-500/60' : 
                    option.value === 'write' ? 'focus:bg-blue-500/10 focus:text-blue-600 focus:border-blue-500/60' :
                    option.value === 'read' ? 'focus:bg-emerald-500/10 focus:text-emerald-600 focus:border-emerald-500/60' :
                    'focus:bg-destructive/10 focus:text-destructive focus:border-destructive/60'
                  ]"
                  @select="inviteDraft.secretAccess = option.value"
                >
                  {{ option.label }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      <!-- Tactile Footer -->
      <SheetFooter class="border-t bg-muted/30 p-6">
        <div class="grid w-full gap-4 sm:grid-cols-2">
          <!-- Matches Local Button (Green) - Safe Primary Action -->
          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-primary/20 bg-primary/3 px-4 text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.15)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none opacity-50 cursor-not-allowed"
            disabled
            type="button"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Dispatch Invitation</span>
          </button>

          <!-- Destructive Action - Placed Far Right -->
          <button
            class="group relative flex h-12 items-center justify-center gap-3 rounded-none border-2 border-destructive/20 bg-destructive/5 px-4 text-destructive transition-all hover:border-destructive/50 hover:bg-destructive/10 hover:shadow-[4px_4px_0_0_rgba(239,68,68,0.22)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
            type="button"
            @click="emit('update:open', false)"
          >
            <span class="font-mono text-[10px] font-black uppercase tracking-widest">Abort Invite</span>
          </button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
