import { PhDatabase, PhFolderOpen, PhKey } from '@phosphor-icons/vue'
import type { AccessLevel, AccessUser, GrantTarget } from '~/types/access'

const defaultUsers: AccessUser[] = [
  {
    id: 'oussema_admin',
    username: 'oussema_admin',
    email: 'oussema@test.io',
    role: 'manager',
    status: 'active',
    inheritedFrom: 'workspace manager',
    grants: {
      collections: { Authentication: 'admin', Realtime: 'admin' },
      environments: { Local: 'admin', Staging: 'admin' },
      secrets: { Local: 'read', Staging: 'read' },
    },
  },
  {
    id: 'nora_dev',
    username: 'nora_dev',
    email: 'nora@example.com',
    role: 'developer',
    status: 'active',
    inheritedFrom: 'project developer',
    grants: {
      collections: { Authentication: 'write', Realtime: 'write' },
      environments: { Local: 'write', Staging: 'read' },
      secrets: { Local: 'none', Staging: 'none' },
    },
  },
  {
    id: 'qa_tester',
    username: 'qa_tester',
    email: 'qa@example.com',
    role: 'tester',
    status: 'pending',
    inheritedFrom: 'manual invite',
    grants: {
      collections: { Authentication: 'read', Realtime: 'none' },
      environments: { Local: 'read', Staging: 'none' },
      secrets: { Local: 'none', Staging: 'none' },
    },
  },
]

const collections = [
  { name: 'Authentication', requests: 3, defaultEnvironment: 'Local' },
  { name: 'Realtime', requests: 2, defaultEnvironment: 'Local' },
]

const environments = [
  { name: 'Local', visibility: 'project', variables: 3, secrets: 1 },
  { name: 'Staging', visibility: 'restricted', variables: 3, secrets: 1 },
]

const accessWeight: Record<AccessLevel, number> = {
  none: 0,
  read: 1,
  write: 2,
  admin: 3,
}

export function useAccess() {
  const users = useState<AccessUser[]>('access:users', () => defaultUsers)
  const selectedUserId = useState<string>('access:selected-user', () => defaultUsers[0]!.id)
  
  const selectedUser = computed(() => users.value.find(user => user.id === selectedUserId.value) ?? users.value[0]!)

  const updateRole = (role: string) => {
    if (!selectedUser.value) return
    selectedUser.value.role = role
    selectedUser.value.inheritedFrom = `project ${role}`
  }

  const kickUser = (id: string) => {
    users.value = users.value.filter(user => user.id !== id)
    if (selectedUserId.value === id) {
      selectedUserId.value = users.value[0]?.id ?? ''
    }
  }

  const updateGrant = (target: GrantTarget, name: string, level: AccessLevel) => {
    if (!selectedUser.value) return
    if (target === 'collection') {
      selectedUser.value.grants.collections[name] = level
    }
    else if (target === 'environment') {
      selectedUser.value.grants.environments[name] = level
    }
    else if (target === 'secret') {
      selectedUser.value.grants.secrets[name] = level
    }
  }

  const resolveDenied = (section: string, name: string) => {
    const targetMap: Record<string, GrantTarget> = {
      Collections: 'collection',
      Environments: 'environment',
      Secrets: 'secret',
    }
    const target = targetMap[section]
    if (target) {
      updateGrant(target, name, 'read')
    }
  }

  const canExecute = (collectionName: string, environmentName: string) => {
    if (!selectedUser.value) return false
    const collectionLevel = selectedUser.value.grants.collections[collectionName] ?? 'none'
    const environmentLevel = selectedUser.value.grants.environments[environmentName] ?? 'none'

    return accessWeight[collectionLevel] > 0 && accessWeight[environmentLevel] > 0
  }

  const grantSections = computed(() => {
    if (!selectedUser.value) return []
    return [
      {
        key: 'collection' as GrantTarget,
        label: 'Collections',
        icon: PhFolderOpen,
        rows: collections.map(collection => ({
          name: collection.name,
          meta: `${collection.requests} requests / default ${collection.defaultEnvironment}`,
          level: selectedUser.value!.grants.collections[collection.name] ?? 'none',
        })),
      },
      {
        key: 'environment' as GrantTarget,
        label: 'Environments',
        icon: PhDatabase,
        rows: environments.map(environment => ({
          name: environment.name,
          meta: `${environment.visibility} / ${environment.variables} variables`,
          level: selectedUser.value!.grants.environments[environment.name] ?? 'none',
        })),
      },
      {
        key: 'secret' as GrantTarget,
        label: 'Secrets',
        icon: PhKey,
        rows: environments.map(environment => ({
          name: environment.name,
          meta: `${environment.secrets} masked value${environment.secrets === 1 ? '' : 's'}`,
          level: selectedUser.value!.grants.secrets[environment.name] ?? 'none',
        })),
      },
    ]
  })

  const executionRows = computed(() =>
    collections.flatMap(collection => environments.map(environment => ({
      name: `${collection.name} -> ${environment.name}`,
      meta: canExecute(collection.name, environment.name) ? 'Access match' : 'Missing grant',
      level: (canExecute(collection.name, environment.name) ? 'read' : 'none') as AccessLevel,
    }))),
  )

  const deniedTargets = computed(() => {
    const targets: Array<{ section: string, name: string, level: AccessLevel }> = []

    grantSections.value.forEach((section) => {
      section.rows.forEach((row) => {
        if (row.level === 'none') {
          targets.push({ section: section.label, name: row.name, level: row.level })
        }
      })
    })

    return targets
  })

  return {
    users,
    selectedUserId,
    selectedUser,
    grantSections,
    executionRows,
    deniedTargets,
    updateRole,
    kickUser,
    updateGrant,
    resolveDenied,
    canExecute,
  }
}
