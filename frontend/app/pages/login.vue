<script setup lang="ts">
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
  <main class="grid min-h-screen place-items-center bg-background px-4">
    <form class="w-full max-w-sm space-y-5 rounded-lg border bg-card p-6 shadow-sm" @submit.prevent="submit">
      <div class="space-y-1">
        <h1 class="font-heading text-2xl font-bold">Sign in</h1>
        <p class="text-sm text-muted-foreground">Use your email or username to continue.</p>
      </div>

      <label class="block space-y-2 text-sm font-medium">
        <span>Email or username</span>
        <input
          v-model="login"
          class="h-10 w-full rounded-md border bg-background px-3 text-sm outline-none transition-colors focus:border-primary"
          required
          type="text"
        >
      </label>

      <label class="block space-y-2 text-sm font-medium">
        <span>Password</span>
        <input
          v-model="password"
          class="h-10 w-full rounded-md border bg-background px-3 text-sm outline-none transition-colors focus:border-primary"
          required
          type="password"
        >
      </label>

      <p v-if="error" class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
        {{ error }}
      </p>

      <button
        class="h-10 w-full rounded-md bg-primary px-4 text-sm font-semibold text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-60"
        :disabled="submitting"
        type="submit"
      >
        {{ submitting ? 'Signing in...' : 'Sign in' }}
      </button>

      <p class="text-center text-sm text-muted-foreground">
        No account?
        <NuxtLink class="font-medium text-primary hover:underline" to="/register">Create one</NuxtLink>
      </p>
    </form>
  </main>
</template>
