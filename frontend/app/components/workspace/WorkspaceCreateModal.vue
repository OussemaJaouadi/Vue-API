<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '~/components/ui/sheet'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const { createWorkspace } = useWorkspace()
const name = ref('')
const loading = ref(false)
const error = ref('')

const handleCreate = async () => {
  if (!name.value.trim()) return
  loading.value = true
  error.value = ''
  try {
    await createWorkspace(name.value.trim())
    emit('update:open', false)
    name.value = ''
  } catch (err: any) {
    error.value = err?.data?.error || 'Failed to create workspace'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Sheet :open="open" @update:open="(val: boolean) => emit('update:open', val)">
    <SheetContent
      class="w-[min(420px,100vw)] max-w-none border-l-2 border-primary/20 bg-background p-0 sm:max-w-none select-none"
      accessibility-title="Create workspace"
      accessibility-description="Create a new isolated workspace"
    >
      <SheetHeader class="border-b bg-muted/30 p-6">
        <div class="flex items-center gap-4">
          <div class="grid size-10 place-items-center border-2 border-primary bg-primary/10 text-primary">
            <PhPlus class="size-5" />
          </div>
          <div>
            <SheetTitle class="font-heading text-lg font-black uppercase tracking-tight">
              Create workspace
            </SheetTitle>
            <SheetDescription class="font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">
              New isolated environment
            </SheetDescription>
          </div>
        </div>
      </SheetHeader>

      <div class="p-6">
        <label class="grid gap-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Workspace name</span>
          <input
            v-model="name"
            class="h-11 rounded-none border-2 border-primary/10 bg-background/50 px-3 font-mono text-sm outline-none transition-all placeholder:text-muted-foreground/70 hover:border-primary/40 hover:bg-background focus:border-primary shadow-[inset_0_1px_2px_rgba(0,0,0,0.05)]"
            placeholder="My workspace"
            type="text"
            @keyup.enter="handleCreate"
          >
        </label>
        <p v-if="error" class="mt-2 font-mono text-[10px] font-bold uppercase tracking-widest text-destructive">{{ error }}</p>
      </div>

      <SheetFooter class="border-t bg-muted/30 p-6">
        <div class="flex w-full gap-3">
          <button
            class="flex h-10 flex-1 items-center justify-center border-2 border-border/40 bg-background px-4 font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground transition-all hover:border-destructive/40 hover:bg-destructive/5 hover:text-destructive active:translate-x-0.5 active:translate-y-0.5"
            type="button"
            @click="emit('update:open', false)"
          >
            Cancel
          </button>
          <button
            class="flex h-10 flex-1 items-center justify-center border-2 border-primary/20 bg-primary/5 px-4 font-mono text-[10px] font-black uppercase tracking-widest text-primary transition-all hover:border-primary/50 hover:bg-primary/10 hover:shadow-[4px_4px_0_0_rgba(16,185,129,0.15)] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none disabled:opacity-50"
            :disabled="!name.trim() || loading"
            type="button"
            @click="handleCreate"
          >
            {{ loading ? 'Creating...' : 'Create' }}
          </button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
