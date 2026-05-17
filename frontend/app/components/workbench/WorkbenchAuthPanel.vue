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
  <div class="bg-card font-sans text-sm">
    <div class="border-b bg-muted/20 p-2 text-xs">
      <UiDropdownMenu>
        <UiDropdownMenuTrigger as-child>
          <button
            class="grid h-11 w-full grid-cols-[20px_minmax(0,1fr)_16px] items-center gap-3 border border-primary/25 bg-background px-3 text-left text-foreground outline-none transition-colors hover:border-primary/50 focus:border-primary/60"
            type="button"
            :title="selectedMode.hint"
          >
            <component :is="selectedMode.icon" class="size-4 text-primary" />
            <span class="min-w-0">
              <span class="block truncate text-[11px] font-black uppercase tracking-widest leading-none">{{ selectedMode.label }}</span>
              <span class="block truncate text-[9px] font-bold uppercase tracking-wider text-muted-foreground leading-none mt-1.5">{{ selectedMode.hint }}</span>
            </span>
            <PhCaretDown class="size-4 text-muted-foreground" />
          </button>
        </UiDropdownMenuTrigger>

        <UiDropdownMenuContent align="start" class="w-72 rounded-none border-2 p-1">
          <UiDropdownMenuItem
            v-for="mode in modes"
            :key="mode.value"
            class="grid grid-cols-[20px_minmax(0,1fr)_16px] gap-3 rounded-none px-3 py-2"
            :class="authMode === mode.value ? 'bg-primary/8 text-foreground' : 'text-muted-foreground'"
            @click="authMode = mode.value"
          >
            <component
              :is="mode.icon"
              class="size-4 self-center"
              :class="authMode === mode.value ? 'text-primary' : 'text-muted-foreground'"
            />
            <span class="min-w-0">
              <span class="block truncate text-[11px] font-black uppercase tracking-widest leading-none">{{ mode.label }}</span>
              <span class="block truncate text-[9px] font-bold uppercase tracking-wider text-muted-foreground leading-none mt-1.5">{{ mode.hint }}</span>
            </span>
            <PhCheck
              v-if="authMode === mode.value"
              class="size-4 self-center text-primary"
            />
          </UiDropdownMenuItem>
        </UiDropdownMenuContent>
      </UiDropdownMenu>
    </div>

    <div v-if="authMode === 'bearer'" class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        Token
      </div>
      <div class="flex min-w-0 items-center gap-3 p-3">
        <input
          v-model="bearerToken"
          :disabled="authMode !== 'bearer'"
          class="h-10 min-w-0 flex-1 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50 disabled:cursor-not-allowed disabled:opacity-45"
          placeholder="{{accessToken}}"
          spellcheck="false"
        >
        <UiBadge
          variant="outline"
          class="h-7 rounded-none border-primary/25 bg-primary/8 px-2 font-mono text-[10px] uppercase text-primary"
        >
          private
        </UiBadge>
      </div>
    </div>

    <div v-else-if="authMode === 'api-key'" class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        API key
      </div>
      <div class="grid min-w-0 grid-cols-[140px_minmax(0,1fr)_110px] gap-3 p-3">
        <input
          v-model="apiKeyName"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="x-api-key"
          spellcheck="false"
        >
        <input
          v-model="apiKeyValue"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="{{apiKey}}"
          spellcheck="false"
        >
        <select
          v-model="apiKeyPlacement"
          class="h-10 border bg-background px-3 font-bold uppercase tracking-widest text-muted-foreground outline-none focus:border-primary/50"
        >
          <option value="header">Header</option>
          <option value="query">Query</option>
        </select>
      </div>
    </div>

    <div v-else-if="authMode === 'basic'" class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        Basic
      </div>
      <div class="grid min-w-0 grid-cols-2 gap-3 p-3">
        <input
          v-model="basicUsername"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="{{username}}"
          spellcheck="false"
        >
        <input
          v-model="basicPassword"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="{{password}}"
          spellcheck="false"
          type="password"
        >
      </div>
    </div>

    <div v-else-if="authMode === 'oauth2'" class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        OAuth2
      </div>
      <div class="grid min-w-0 grid-cols-1 gap-3 p-3 sm:grid-cols-2">
        <select
          v-model="oauthGrant"
          class="h-10 min-w-0 border bg-background px-3 font-bold uppercase tracking-widest text-muted-foreground outline-none focus:border-primary/50"
        >
          <option value="authorization-code-pkce">Auth Code + PKCE</option>
          <option value="client-credentials">Client Credentials</option>
          <option value="refresh-token">Refresh Token</option>
        </select>
        <input
          v-model="oauthClientId"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Client ID (e.g. {{clientId}})"
          spellcheck="false"
        >
        <input
          v-model="oauthAccessToken"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Access Token (e.g. {{oauthAccessToken}})"
          spellcheck="false"
        >
        <input
          v-model="oauthTokenUrl"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50 sm:col-span-2"
          placeholder="Token URL (e.g. {{tokenUrl}})"
          spellcheck="false"
        >
        <input
          v-model="oauthScopes"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Scopes (e.g. openid profile email)"
          spellcheck="false"
        >
      </div>
    </div>

    <div v-else-if="authMode === 'oidc'" class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        OIDC
      </div>
      <div class="grid min-w-0 grid-cols-1 gap-3 p-3">
        <input
          v-model="oidcIssuerUrl"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Issuer URL (e.g. {{issuerUrl}})"
          spellcheck="false"
        >
        <input
          v-model="oidcAudience"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Audience (e.g. {{audience}})"
          spellcheck="false"
        >
        <input
          v-model="oauthScopes"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Scopes (e.g. openid profile email)"
          spellcheck="false"
        >
        <input
          v-model="oauthAccessToken"
          class="h-10 min-w-0 border bg-background px-4 font-bold text-foreground outline-none transition-colors placeholder:text-muted-foreground/35 focus:border-primary/50"
          placeholder="Access Token (e.g. {{oidcAccessToken}})"
          spellcheck="false"
        >
      </div>
    </div>

    <div v-else class="grid grid-cols-[110px_minmax(0,1fr)] border-b text-xs">
      <div class="flex items-center border-r bg-muted/5 px-4 font-black uppercase tracking-widest text-muted-foreground/60">
        Config
      </div>
      <div class="px-4 py-3 text-[11px] font-bold uppercase tracking-wider text-muted-foreground/70">
        <span v-if="authMode === 'inherit'">Use parent auth profile.</span>
        <span v-else>No request auth applied.</span>
      </div>
    </div>

    <div class="p-3">
      <div class="border bg-background">
        <div class="grid grid-cols-[120px_minmax(0,1fr)_150px] border-b bg-muted/10 py-2 text-[10px] font-black uppercase tracking-widest text-muted-foreground/40">
          <span class="px-4">Header</span>
          <span class="px-4">Value</span>
          <span class="px-4 text-right">Source</span>
        </div>

        <div
          v-for="row in previewRows"
          :key="`${row.key}-${row.source}`"
          class="grid min-h-12 grid-cols-[120px_minmax(0,1fr)_150px] items-center border-b text-xs last:border-b-0"
        >
          <div class="border-r px-4 font-black text-foreground/80">{{ row.key }}</div>
          <div class="min-w-0 truncate border-r px-4 text-muted-foreground font-bold font-mono">
            {{ row.value }}
          </div>
          <div class="px-4 text-right">
            <UiBadge
              variant="outline"
              class="h-6 max-w-full rounded-none px-2 font-mono text-[9px] uppercase"
              :class="row.source === 'off' ? 'border-muted-foreground/20 text-muted-foreground' : 'border-primary/25 bg-primary/8 text-primary'"
            >
              <span class="block truncate">{{ row.source }}</span>
            </UiBadge>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
