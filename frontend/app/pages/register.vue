<script setup lang="ts">
import { PhUserPlus as UserPlus } from '@phosphor-icons/vue'

const auth = useAuthSession()

const email = ref('')
const username = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)

onMounted(async () => {
  if (await auth.loadMe()) {
    await navigateTo('/')
  }
})

const submit = async () => {
  error.value = ''
  submitting.value = true
  try {
    await auth.register({
      email: email.value,
      username: username.value,
      password: password.value,
    })
    await navigateTo('/')
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || 'Registration failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="grid min-h-dvh bg-background lg:grid-cols-[1fr_440px]">
    <section class="hidden border-r bg-muted/30 p-8 lg:flex lg:flex-col lg:justify-between">
      <div class="space-y-2">
        <AppLogo mode="lockup" class="h-12 w-auto" />
        <div class="text-xs text-muted-foreground">Protected workspace access.</div>
      </div>

      <div class="max-w-xl space-y-4">
        <AppLogo class="h-28 w-28 opacity-90" />
        <UiBadge variant="secondary">Registration queue</UiBadge>
        <h1 class="font-heading text-4xl font-semibold leading-tight">
          Create your account, then wait for workspace assignment.
        </h1>
        <p class="text-sm text-muted-foreground">
          Managers receive the registration event in real time and grant access per workspace.
        </p>
      </div>
    </section>

    <section class="grid place-items-center px-4 py-10">
      <form class="w-full max-w-sm" @submit.prevent="submit">
        <UiCard>
          <UiCardHeader>
            <UiCardTitle>Create account</UiCardTitle>
            <UiCardDescription>Register first, then a manager grants workspace access.</UiCardDescription>
          </UiCardHeader>
          <UiCardContent class="space-y-4">
            <div class="space-y-2">
              <UiLabel for="email">Email</UiLabel>
              <UiInput
                id="email"
                v-model="email"
                autocomplete="email"
                required
                type="email"
              />
            </div>

            <div class="space-y-2">
              <UiLabel for="username">Username</UiLabel>
              <UiInput
                id="username"
                v-model="username"
                autocomplete="username"
                required
                type="text"
              />
            </div>

            <div class="space-y-2">
              <UiLabel for="password">Password</UiLabel>
              <UiInput
                id="password"
                v-model="password"
                autocomplete="new-password"
                minlength="12"
                required
                type="password"
              />
            </div>

            <UiAlert v-if="error" variant="destructive">
              <UiAlertDescription>{{ error }}</UiAlertDescription>
            </UiAlert>
          </UiCardContent>
          <UiCardFooter class="flex-col gap-3">
            <UiButton class="w-full" :disabled="submitting" type="submit">
              <UserPlus class="size-4" />
              {{ submitting ? 'Creating...' : 'Create account' }}
            </UiButton>
            <p class="text-center text-xs text-muted-foreground">
              Already registered?
              <NuxtLink class="font-medium text-primary hover:underline" to="/login">Sign in</NuxtLink>
            </p>
          </UiCardFooter>
        </UiCard>
      </form>
    </section>
  </main>
</template>
