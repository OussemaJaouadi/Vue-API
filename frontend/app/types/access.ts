export type AccessLevel = 'none' | 'read' | 'write' | 'admin'
export type GrantTarget = 'collection' | 'environment' | 'secret' | 'execution'

export interface AccessUser {
  id: string
  username: string
  email: string
  role: string
  status: string
  inheritedFrom: string
  grants: {
    collections: Record<string, AccessLevel>
    environments: Record<string, AccessLevel>
    secrets: Record<string, AccessLevel>
  }
}

export const accessWeight: Record<AccessLevel, number> = {
  none: 0,
  read: 1,
  write: 2,
  admin: 3,
}

export const accessTone = (level: AccessLevel) => {
  if (level === 'admin') return 'border-amber-500/50 bg-amber-500/10 text-amber-600 dark:text-amber-400 shadow-[0_0_15px_rgba(245,158,11,0.12)]'
  if (level === 'write') return 'border-blue-500/50 bg-blue-500/10 text-blue-600 dark:text-blue-400 shadow-[0_0_15px_rgba(59,130,246,0.12)]'
  if (level === 'read') return 'border-emerald-500/50 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 shadow-[0_0_15px_rgba(16,185,129,0.12)]'
  return 'border-dashed border-border/60 bg-muted/5 text-muted-foreground/65'
}

export const roleTone = (role: string) => {
  const r = role.toLowerCase()
  if (r.includes('manager') || r.includes('admin')) return 'border-indigo-500/50 bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 shadow-[0_0_15px_rgba(99,102,241,0.12)]'
  if (r.includes('developer')) return 'border-sky-500/50 bg-sky-500/10 text-sky-600 dark:text-sky-400 shadow-[0_0_15px_rgba(14,165,233,0.12)]'
  return 'border-teal-500/50 bg-teal-500/10 text-teal-600 dark:text-teal-400 shadow-[0_0_15px_rgba(20,184,166,0.12)]'
}

