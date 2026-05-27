<script setup lang="ts">
import {
  PhCaretDown,
  PhCheck,
  PhKey,
  PhLockKey,
  PhPassword,
  PhShieldCheck,
  PhUserCircle,
} from '@phosphor-icons/vue'
import type { Component } from 'vue'

type AuthMode = 'inherit' | 'bearer' | 'api-key' | 'basic' | 'oauth2' | 'oidc' | 'none'
type ApiKeyPlacement = 'header' | 'query'
type OAuthGrant = 'authorization-code-pkce' | 'client-credentials' | 'refresh-token'

const authMode = useState<AuthMode>('workbench:auth-mode', () => 'bearer')
const bearerToken = useState<string>('workbench:auth-bearer-token', () => '{{accessToken}}')
const apiKeyName = useState<string>('workbench:auth-api-key-name', () => 'x-api-key')
const apiKeyValue = useState<string>('workbench:auth-api-key-value', () => '{{apiKey}}')
const apiKeyPlacement = useState<ApiKeyPlacement>('workbench:auth-api-key-placement', () => 'header')
const basicUsername = useState<string>('workbench:auth-basic-username', () => '{{username}}')
const basicPassword = useState<string>('workbench:auth-basic-password', () => '{{password}}')
const oauthGrant = useState<OAuthGrant>('workbench:auth-oauth-grant', () => 'authorization-code-pkce')
const oauthAccessToken = useState<string>('workbench:auth-oauth-access-token', () => '{{oauthAccessToken}}')
const oauthClientId = useState<string>('workbench:auth-oauth-client-id', () => '{{clientId}}')
const oauthTokenUrl = useState<string>('workbench:auth-oauth-token-url', () => '{{tokenUrl}}')
const oauthScopes = useState<string>('workbench:auth-oauth-scopes', () => 'openid profile email')
const oidcIssuerUrl = useState<string>('workbench:auth-oidc-issuer-url', () => '{{issuerUrl}}')
const oidcAudience = useState<string>('workbench:auth-oidc-audience', () => '{{audience}}')

const modes: Array<{
  value: AuthMode
  label: string
  hint: string
  icon: Component
}> = [
  { value: 'inherit', label: 'Inherit', hint: 'Parent profile', icon: PhLockKey },
  { value: 'bearer', label: 'Bearer', hint: 'Authorization header', icon: PhKey },
  { value: 'api-key', label: 'API Key', hint: 'Header or query key', icon: PhShieldCheck },
  { value: 'basic', label: 'Basic', hint: 'Username/password', icon: PhPassword },
  { value: 'oauth2', label: 'OAuth2', hint: 'Token flow output', icon: PhKey },
  { value: 'oidc', label: 'OIDC', hint: 'Discovery token', icon: PhShieldCheck },
  { value: 'none', label: 'None', hint: 'No auth helper', icon: PhUserCircle },
]

const selectedMode = computed(() => modes.find(mode => mode.value === authMode.value) ?? modes[0])

const previewRows = computed(() => {
  if (authMode.value === 'bearer') {
    return [{ key: 'Authorization', value: `Bearer ${bearerToken.value}`, source: 'auth' }]
  }

  if (authMode.value === 'api-key') {
    return [{
      key: apiKeyName.value || 'x-api-key',
      value: apiKeyValue.value,
      source: apiKeyPlacement.value,
    }]
  }

  if (authMode.value === 'basic') {
    return [{ key: 'Authorization', value: `Basic base64(${basicUsername.value}:••••••••)`, source: 'auth' }]
  }

  if (authMode.value === 'oauth2') {
    return [{ key: 'Authorization', value: `Bearer ${oauthAccessToken.value}`, source: oauthGrant.value }]
  }

  if (authMode.value === 'oidc') {
    return [{ key: 'Authorization', value: `Bearer ${oauthAccessToken.value}`, source: 'oidc' }]
  }

  if (authMode.value === 'inherit') {
    return [{ key: 'Authorization', value: 'Inherited from parent auth profile', source: 'parent' }]
  }

  return [{ key: 'Authorization', value: 'Not sent', source: 'off' }]
})
</script>

<template>
  <div class="flex flex-col gap-6 select-none">
    <!-- Mode Selector -->
    <div class="space-y-2">
      <label class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Authentication Protocol</label>
      <UiDropdownMenu>
        <UiDropdownMenuTrigger as-child>
          <button
            class="grid h-12 w-full grid-cols-[24px_minmax(0,1fr)_20px] items-center gap-4 border-2 border-primary/15 bg-background px-4 text-left outline-none transition-all hover:border-primary/40 focus:border-primary shadow-sm"
            type="button"
          >
            <component :is="selectedMode.icon" class="size-5 text-primary" />
            <div class="min-w-0">
              <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight text-foreground">{{ selectedMode.label }}</span>
              <span class="block truncate font-mono text-[9px] font-bold uppercase tracking-widest text-muted-foreground/90">{{ selectedMode.hint }}</span>
            </div>
            <PhCaretDown class="size-4 opacity-70" />
          </button>
        </UiDropdownMenuTrigger>

        <UiDropdownMenuContent align="start" class="w-72 rounded-none border-2 border-primary/20 bg-background/95 backdrop-blur-xl p-1 shadow-[8px_8px_0_0_rgba(16,185,129,0.12)]">
          <UiDropdownMenuItem
            v-for="mode in modes"
            :key="mode.value"
            class="grid grid-cols-[24px_minmax(0,1fr)_20px] gap-4 rounded-none px-3 py-2.5 transition-all mb-0.5 last:mb-0 border-l-2"
            :class="authMode === mode.value ? 'bg-primary/10 text-primary border-primary' : 'border-transparent text-muted-foreground focus:bg-primary/5 focus:text-foreground'"
            @click="authMode = mode.value"
          >
            <component :is="mode.icon" class="size-5 shrink-0" />
            <div class="min-w-0">
              <span class="block truncate font-mono text-[11px] font-black uppercase tracking-tight">{{ mode.label }}</span>
              <span class="block truncate font-mono text-[9px] font-bold uppercase tracking-widest opacity-60">{{ mode.hint }}</span>
            </div>
            <PhCheck v-if="authMode === mode.value" class="size-4 shrink-0" />
          </UiDropdownMenuItem>
        </UiDropdownMenuContent>
      </UiDropdownMenu>
    </div>

    <!-- Configuration Context -->
    <div class="space-y-4 border-t-2 border-dashed border-border/20 pt-6">
      <div v-if="authMode === 'bearer'" class="grid gap-2">
        <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Access Token</span>
        <div class="flex items-center gap-3">
          <input
            v-model="bearerToken"
            class="h-11 flex-1 border-2 border-primary/10 bg-background/50 px-4 font-mono text-sm font-bold outline-none transition-all placeholder:text-muted-foreground/20 hover:border-primary/30 focus:border-primary"
            placeholder="{{accessToken}}"
            spellcheck="false"
          >
          <div class="h-11 border-2 border-primary/20 bg-primary/5 px-3 flex items-center font-mono text-[9px] font-black uppercase tracking-widest text-primary shrink-0">
            Encrypted
          </div>
        </div>
      </div>

      <div v-else-if="authMode === 'api-key'" class="grid gap-4">
        <div class="grid grid-cols-[1fr_2fr] gap-4">
          <label class="grid gap-2">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Header Key</span>
            <input v-model="apiKeyName" class="h-11 border-2 border-primary/10 bg-background/50 px-4 font-mono text-sm font-bold outline-none transition-all focus:border-primary" placeholder="x-api-key">
          </label>
          <label class="grid gap-2">
            <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Key Value</span>
            <input v-model="apiKeyValue" class="h-11 border-2 border-primary/10 bg-background/50 px-4 font-mono text-sm font-bold outline-none transition-all focus:border-primary" placeholder="{{apiKey}}">
          </label>
        </div>
        <label class="grid gap-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Transport Placement</span>
          <div class="grid grid-cols-2 gap-2">
            <button v-for="p in ['header', 'query']" :key="p" @click="apiKeyPlacement = p as any" 
              class="h-10 border-2 font-mono text-[10px] font-black uppercase tracking-widest transition-all"
              :class="apiKeyPlacement === p ? 'border-primary bg-primary/10 text-primary' : 'border-border/40 bg-muted/5 text-muted-foreground hover:border-primary/20'"
            >
              {{ p }}
            </button>
          </div>
        </label>
      </div>

      <div v-else-if="authMode === 'basic'" class="grid grid-cols-2 gap-4">
        <label class="grid gap-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Username</span>
          <input v-model="basicUsername" class="h-11 border-2 border-primary/10 bg-background/50 px-4 font-mono text-sm font-bold outline-none transition-all focus:border-primary" placeholder="{{username}}">
        </label>
        <label class="grid gap-2">
          <span class="font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">Password</span>
          <input v-model="basicPassword" type="password" class="h-11 border-2 border-primary/10 bg-background/50 px-4 font-mono text-sm font-bold outline-none transition-all focus:border-primary" placeholder="{{password}}">
        </label>
      </div>

      <div v-else class="py-12 text-center">
        <p class="font-mono text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground/20 italic">
          {{ authMode === 'inherit' ? 'Inheriting from parent authority' : 'Protocol selection pending' }}
        </p>
      </div>
    </div>

    <!-- Live Preview Strip -->
    <div class="mt-4 border-2 border-primary/10 bg-muted/5 overflow-hidden">
      <div class="flex h-8 items-center border-b border-primary/5 bg-primary/5 px-4 font-mono text-[9px] font-black uppercase tracking-widest text-primary/60">
        Injected Manifest Preview
      </div>
      <div class="divide-y divide-border/10">
        <div v-for="row in previewRows" :key="row.key" class="grid grid-cols-[140px_1fr_100px] h-12 items-center px-4 font-mono text-[10px]">
          <span class="font-black text-foreground/80 border-r border-border/5 h-full flex items-center">{{ row.key }}</span>
          <span class="truncate px-4 font-bold text-muted-foreground">{{ row.value }}</span>
          <span class="text-right border-l border-border/5 h-full flex items-center justify-end">
            <span class="px-2 py-0.5 border text-[8px] font-black uppercase tracking-tighter" :class="row.source === 'off' ? 'border-border text-muted-foreground/20' : 'border-primary/30 text-primary'">
              {{ row.source }}
            </span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
