<script setup lang="ts">
import {
  PhCheck,
  PhDesktop,
  PhGearSix,
  PhMoon,
  PhSun,
} from '@phosphor-icons/vue'
import type { Component } from 'vue'
import type { ThemePreference } from '~/composables/useThemePreference'

const theme = useThemePreference()

const themeItems: Array<{
  value: ThemePreference
  label: string
  description: string
  icon: Component
}> = [
  { value: 'light', label: 'Light', description: 'Use the light workbench palette.', icon: PhSun },
  { value: 'system', label: 'System', description: 'Follow the operating system preference.', icon: PhDesktop },
  { value: 'dark', label: 'Dark', description: 'Use the darker workbench palette.', icon: PhMoon },
]
</script>

<template>
  <div class="min-h-[calc(100dvh-5.5rem)] border bg-card">
    <header class="flex h-14 items-center justify-between border-b bg-muted/20 px-4">
      <div class="flex items-center gap-3">
        <div class="grid size-8 place-items-center border-2 border-primary/25 bg-primary/10 text-primary">
          <PhGearSix class="size-4" />
        </div>
        <div>
          <h1 class="font-heading text-base font-bold">Settings</h1>
          <p class="font-mono text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
            User preferences
          </p>
        </div>
      </div>
    </header>

    <div class="grid gap-3 p-3 lg:grid-cols-[260px_minmax(0,1fr)]">
      <aside class="border bg-background">
        <div class="border-b bg-muted/15 px-3 py-2 font-mono text-[10px] font-black uppercase tracking-widest text-muted-foreground">
          Sections
        </div>
        <div class="border-b bg-primary/8 px-3 py-3 font-mono text-[11px] font-black uppercase tracking-widest text-primary">
          Appearance
        </div>
      </aside>

      <section class="border bg-background">
        <div class="border-b bg-muted/15 px-4 py-3">
          <h2 class="font-mono text-[11px] font-black uppercase tracking-widest">Theme preference</h2>
          <p class="mt-1 text-sm text-muted-foreground">Stored locally for now. Backend user preferences can replace this later.</p>
        </div>

        <div class="grid gap-2 p-4 md:grid-cols-3">
          <button
            v-for="item in themeItems"
            :key="item.value"
            class="relative border p-4 text-left transition-colors"
            :class="theme.preference.value === item.value ? 'border-primary bg-primary/8 text-foreground' : 'border-border bg-card text-muted-foreground hover:border-primary/35 hover:text-foreground'"
            type="button"
            @click="theme.setTheme(item.value)"
          >
            <div class="mb-4 flex items-center justify-between">
              <component :is="item.icon" class="size-5" :class="theme.preference.value === item.value ? 'text-primary' : 'text-muted-foreground'" />
              <PhCheck v-if="theme.preference.value === item.value" class="size-4 text-primary" />
            </div>
            <h3 class="font-mono text-[11px] font-black uppercase tracking-widest">{{ item.label }}</h3>
            <p class="mt-2 text-sm leading-5 text-muted-foreground">{{ item.description }}</p>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>
