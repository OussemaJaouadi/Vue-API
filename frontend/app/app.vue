<script setup lang="ts">
useThemePreference()
</script>

<template>
  <UiTooltipProvider>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiToaster position="top-right" />
  </UiTooltipProvider>
</template>
