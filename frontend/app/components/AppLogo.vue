<script setup lang="ts">
import logoMark from '~/assets/brand/logo-mark.png'

const props = withDefaults(defineProps<{
  mode?: 'mark' | 'lockup'
  class?: string
}>(), {
  mode: 'mark',
})
</script>

<template>
  <div :class="props.class">
    <img
      v-if="mode === 'mark'"
      alt="Vue API Workbench"
      class="h-full w-full object-contain"
      :src="logoMark"
    >
    <div v-else class="flex items-center gap-3">
      <img
        alt="Vue API Workbench"
        class="size-12 object-contain"
        :src="logoMark"
      >
      <div class="leading-none">
        <div class="font-heading text-sm font-semibold text-foreground">Vue API</div>
        <div class="mt-1 text-xs font-medium text-muted-foreground">Workbench</div>
      </div>
    </div>
  </div>
</template>
