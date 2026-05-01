<script setup lang="ts">
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
  <main class="grid min-h-screen place-items-center bg-background px-4">
    <form class="w-full max-w-sm space-y-5 rounded-lg border bg-card p-6 shadow-sm" @submit.prevent="submit">
      <div class="space-y-1">
        <h1 class="font-heading text-2xl font-bold">Create account</h1>
        <p class="text-sm text-muted-foreground">Register first, then a manager grants workspace access.</p>
      </div>

      <label class="block space-y-2 text-sm font-medium">
        <span>Email</span>
        <input
          v-model="email"
          class="h-10 w-full rounded-md border bg-background px-3 text-sm outline-none transition-colors focus:border-primary"
          required
          type="email"
        >
      </label>

      <label class="block space-y-2 text-sm font-medium">
        <span>Username</span>
        <input
          v-model="username"
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
          minlength="12"
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
        {{ submitting ? 'Creating...' : 'Create account' }}
      </button>

      <p class="text-center text-sm text-muted-foreground">
        Already registered?
        <NuxtLink class="font-medium text-primary hover:underline" to="/login">Sign in</NuxtLink>
      </p>
    </form>
  </main>
</template>
