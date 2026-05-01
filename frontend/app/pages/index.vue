<script setup lang="ts">
const auth = useAuthSession()

const loading = ref(true)

onMounted(async () => {
  const ok = await auth.loadMe()
  loading.value = false
  if (!ok) {
    await navigateTo('/login')
  }
})
</script>

<template>
  <div class="mx-auto max-w-5xl space-y-6">
    <div v-if="loading" class="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
      Restoring session...
    </div>

    <template v-else-if="auth.user.value">
      <section class="space-y-2">
        <h1 class="font-heading text-2xl font-bold">
          {{ auth.isManager.value ? 'Manager console' : 'Waiting for workspace access' }}
        </h1>
        <p class="max-w-2xl text-sm text-muted-foreground">
          <span v-if="auth.isManager.value">
            New users will appear here in real time once the manager workflow is connected.
          </span>
          <span v-else>
            Your account is active. A manager still needs to add you to a workspace before you can work.
          </span>
        </p>
      </section>

      <section class="grid gap-4 md:grid-cols-3">
        <div class="rounded-lg border bg-card p-4">
          <div class="text-xs font-semibold uppercase text-muted-foreground">Signed in as</div>
          <div class="mt-2 text-sm font-medium">{{ auth.user.value.username }}</div>
          <div class="text-sm text-muted-foreground">{{ auth.user.value.email }}</div>
        </div>

        <div class="rounded-lg border bg-card p-4">
          <div class="text-xs font-semibold uppercase text-muted-foreground">Global role</div>
          <div class="mt-2 text-sm font-medium">{{ auth.user.value.globalRole }}</div>
        </div>

        <div class="rounded-lg border bg-card p-4">
          <div class="text-xs font-semibold uppercase text-muted-foreground">Realtime</div>
          <div class="mt-2 text-sm font-medium">
            {{ auth.lastEvent.value?.type || 'connecting' }}
          </div>
        </div>
      </section>

      <section v-if="!auth.isManager.value" class="rounded-lg border bg-card p-6">
        <div class="space-y-2">
          <h2 class="font-heading text-lg font-semibold">Access pending</h2>
          <p class="text-sm text-muted-foreground">
            Keep this page open. When a manager grants workspace access, the app receives a realtime event and can refresh your workspace list.
          </p>
        </div>
      </section>

      <section v-else class="rounded-lg border bg-card p-6">
        <div class="space-y-2">
          <h2 class="font-heading text-lg font-semibold">Pending users</h2>
          <p class="text-sm text-muted-foreground">
            The next backend slice will expose pending users and workspace assignment actions.
          </p>
        </div>
      </section>
    </template>
  </div>
</template>
