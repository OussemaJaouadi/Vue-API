<script setup lang="ts">
import {
  PhActivity as Activity,
  PhBellRinging as BellRinging,
  PhCheckCircle as CheckCircle,
  PhClock as Clock,
  PhShieldCheck as ShieldCheck,
  PhUserCircle as UserCircle,
  PhUsersThree as UsersThree,
} from '@phosphor-icons/vue'

const auth = useAuthSession()

const loading = ref(true)
const isManager = computed(() => auth.isManager.value)
const eventLabel = computed(() => auth.lastEvent.value?.type || 'connecting')

onMounted(async () => {
  const ok = await auth.loadMe()
  loading.value = false
  if (!ok) {
    await navigateTo('/login')
  }
})
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6">
    <div v-if="loading" class="grid gap-4 md:grid-cols-3">
      <UiSkeleton class="h-28" />
      <UiSkeleton class="h-28" />
      <UiSkeleton class="h-28" />
    </div>

    <template v-else-if="auth.user.value">
      <section class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <UiBadge :variant="isManager ? 'default' : 'secondary'">
            <ShieldCheck class="size-3" />
            {{ auth.user.value.globalRole }}
          </UiBadge>
          <h1 class="font-heading text-2xl font-semibold">
            {{ isManager ? 'Manager console' : 'Workspace access pending' }}
          </h1>
          <p class="max-w-2xl text-sm text-muted-foreground">
            <span v-if="isManager">
              Review new users and assign workspace roles when the next backend slice exposes the queue.
            </span>
            <span v-else>
              Your account exists. A manager still needs to add you to a workspace.
            </span>
          </p>
        </div>
        <UiBadge variant="outline">
          <Activity class="size-3" />
          {{ eventLabel }}
        </UiBadge>
      </section>

      <section class="grid gap-4 md:grid-cols-3">
        <UiCard>
          <UiCardHeader>
            <UiCardTitle class="flex items-center gap-2 text-sm">
              <UserCircle class="size-4" />
              Signed in
            </UiCardTitle>
            <UiCardDescription>{{ auth.user.value.email }}</UiCardDescription>
          </UiCardHeader>
          <UiCardContent>
            <div class="text-lg font-semibold">{{ auth.user.value.username }}</div>
          </UiCardContent>
        </UiCard>

        <UiCard>
          <UiCardHeader>
            <UiCardTitle class="flex items-center gap-2 text-sm">
              <ShieldCheck class="size-4" />
              Global role
            </UiCardTitle>
            <UiCardDescription>Application-level access</UiCardDescription>
          </UiCardHeader>
          <UiCardContent>
            <UiBadge>{{ auth.user.value.globalRole }}</UiBadge>
          </UiCardContent>
        </UiCard>

        <UiCard>
          <UiCardHeader>
            <UiCardTitle class="flex items-center gap-2 text-sm">
              <BellRinging class="size-4" />
              Realtime
            </UiCardTitle>
            <UiCardDescription>Latest event channel state</UiCardDescription>
          </UiCardHeader>
          <UiCardContent>
            <UiBadge variant="outline">{{ eventLabel }}</UiBadge>
          </UiCardContent>
        </UiCard>
      </section>

      <UiTabs :default-value="isManager ? 'queue' : 'access'">
        <UiTabsList>
          <UiTabsTrigger :value="isManager ? 'queue' : 'access'">
            {{ isManager ? 'Pending users' : 'Access' }}
          </UiTabsTrigger>
          <UiTabsTrigger value="events">Events</UiTabsTrigger>
        </UiTabsList>

        <UiTabsContent value="access">
          <UiAlert>
            <Clock class="size-4" />
            <UiAlertTitle>Access pending</UiAlertTitle>
            <UiAlertDescription>
              Keep this page open. When a manager grants workspace access, the app receives a realtime event and can refresh your workspace list.
            </UiAlertDescription>
          </UiAlert>
        </UiTabsContent>

        <UiTabsContent value="queue">
          <UiCard>
            <UiCardHeader>
              <UiCardTitle class="flex items-center gap-2">
                <UsersThree class="size-4" />
                Pending users
              </UiCardTitle>
              <UiCardDescription>Workspace assignment actions will land in the next backend slice.</UiCardDescription>
            </UiCardHeader>
            <UiCardContent class="space-y-3">
              <div class="flex items-center justify-between border px-3 py-2">
                <div>
                  <div class="text-sm font-medium">Registration stream</div>
                  <div class="text-xs text-muted-foreground">{{ eventLabel }}</div>
                </div>
                <UiBadge variant="secondary">
                  <CheckCircle class="size-3" />
                  listening
                </UiBadge>
              </div>
            </UiCardContent>
          </UiCard>
        </UiTabsContent>

        <UiTabsContent value="events">
          <UiCard>
            <UiCardHeader>
              <UiCardTitle>Event channel</UiCardTitle>
              <UiCardDescription>Current browser connection state.</UiCardDescription>
            </UiCardHeader>
            <UiCardContent>
              <pre class="overflow-auto border bg-muted p-3 text-xs">{{ auth.lastEvent.value || { type: 'connecting' } }}</pre>
            </UiCardContent>
          </UiCard>
        </UiTabsContent>
      </UiTabs>
    </template>
  </div>
</template>
