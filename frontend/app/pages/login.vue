<script setup lang="ts">
import { PhSignIn as SignIn } from '@phosphor-icons/vue'

const auth = useAuthSession()

const login = ref('')
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
    await auth.login({
      login: login.value,
      password: password.value,
    })
    await navigateTo('/')
  } catch (err: any) {
    error.value = err?.data?.error || err?.message || 'Login failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="grid min-h-dvh bg-background lg:grid-cols-[1fr_420px]">
    <section class="hidden border-r bg-muted/30 p-8 lg:flex lg:flex-col lg:justify-between">
      <div class="space-y-2">
        <AppLogo mode="lockup" class="h-12 w-auto" />
        <div class="text-xs text-muted-foreground">Collections, environments, execution.</div>
      </div>

      <div class="max-w-xl space-y-4">
        <AppLogo class="h-28 w-28 opacity-90" />
        <UiBadge variant="secondary">Auth gateway</UiBadge>
        <h1 class="font-heading text-4xl font-semibold leading-tight">
          Sign in to continue working with protected APIs.
        </h1>
        <p class="text-sm text-muted-foreground">
          Access tokens stay in memory. Refresh is handled through the backend cookie lifecycle.
        </p>
      </div>
    </section>

    <section class="grid place-items-center px-4 py-10">
      <form class="w-full max-w-sm" @submit.prevent="submit">
        <UiCard>
          <UiCardHeader>
            <UiCardTitle>Sign in</UiCardTitle>
            <UiCardDescription>Use your email or username to continue.</UiCardDescription>
          </UiCardHeader>
          <UiCardContent class="space-y-4">
            <div class="space-y-2">
              <UiLabel for="login">Email or username</UiLabel>
              <UiInput
                id="login"
                v-model="login"
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
                autocomplete="current-password"
                required
                type="password"
              />
            </div>

            <UiAlert v-if="error" variant="destructive">
              <UiAlertDescription>{{ error }}</UiAlertDescription>
            </UiAlert>
          </UiCardContent>
          <UiCardFooter class="flex-col gap-3">
            <UiButton
              class="w-full"
              :disabled="submitting"
              type="submit"
            >
              <SignIn class="size-4" />
              {{ submitting ? 'Signing in...' : 'Sign in' }}
            </UiButton>
            <p class="text-center text-xs text-muted-foreground">
              No account?
              <NuxtLink class="font-medium text-primary hover:underline" to="/register">Create one</NuxtLink>
            </p>
          </UiCardFooter>
        </UiCard>
      </form>
    </section>
  </main>
</template>
