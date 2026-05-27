<script setup lang="ts">
import { PhCircleNotch, PhSignIn as SignIn } from '@phosphor-icons/vue'

definePageMeta({
  layout: 'auth',
})

const auth = useAuthSession()

const login = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)
const inputClass = 'rounded-none border-2 border-muted-foreground/20 h-12 px-4 focus-visible:ring-0 focus-visible:border-primary transition-all duration-200 bg-muted/20 font-mono text-sm placeholder:text-muted-foreground/80 focus-visible:bg-background focus-visible:shadow-[0_0_0_4px_oklch(0.508_0.118_165.612_/_0.10)] disabled:opacity-60 disabled:cursor-not-allowed'

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
    const backendErr = err?.data?.error
    error.value = typeof backendErr === 'string' ? backendErr : err?.message || 'Login failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="auth-form-in mx-auto w-full max-w-sm">
    <div class="mb-10">
      <AppLogo mode="lockup" class="h-10 w-auto lg:hidden mb-8" />
      <h2 class="font-heading text-4xl font-bold tracking-tight">Sign in</h2>
      <p class="mt-2 text-sm text-muted-foreground font-medium">
        Use your email or username to continue.
      </p>
    </div>

    <form class="space-y-6" @submit.prevent="submit">
      <div class="group space-y-2">
        <UiLabel for="login" class="font-bold uppercase tracking-widest text-muted-foreground transition-colors duration-200 group-focus-within:text-primary text-[10px]">Identity</UiLabel>
        <UiInput
          id="login"
          v-model="login"
          autocomplete="username"
          :disabled="submitting"
          required
          type="text"
          placeholder="username or email"
          :class="[inputClass, error && 'border-destructive/50']"
        />
      </div>

      <div class="group space-y-2">
        <div class="flex items-center justify-between">
          <UiLabel for="password" class="font-bold uppercase tracking-widest text-muted-foreground transition-colors duration-200 group-focus-within:text-primary text-[10px]">Password</UiLabel>
        </div>
        <UiInput
          id="password"
          v-model="password"
          autocomplete="current-password"
          :disabled="submitting"
          required
          type="password"
          :class="[inputClass, error && 'border-destructive/50']"
        />
      </div>

      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="-translate-y-1 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="-translate-y-1 opacity-0"
      >
        <UiAlert v-if="error" variant="destructive" class="mt-4 rounded-none border-2">
          <UiAlertDescription class="font-medium">{{ error }}</UiAlertDescription>
        </UiAlert>
      </Transition>

      <UiButton
        class="mt-8 h-14 w-full rounded-none text-sm font-bold uppercase tracking-widest shadow-primary/10 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-primary/15 active:translate-y-0 active:scale-[0.99]"
        :disabled="submitting"
        type="submit"
      >
        <PhCircleNotch v-if="submitting" class="mr-2 size-5 animate-spin" />
        <SignIn v-else class="mr-2 size-5" />
        {{ submitting ? 'Authenticating...' : 'Sign in' }}
      </UiButton>
    </form>

    <p class="mt-10 text-center text-sm text-muted-foreground font-medium">
      No account?
      <NuxtLink class="font-bold text-primary hover:text-primary/80 transition-colors uppercase text-xs tracking-wider ml-1" to="/register">Create one</NuxtLink>
    </p>
  </div>
</template>
