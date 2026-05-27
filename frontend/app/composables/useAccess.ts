import { PhDatabase, PhFolderOpen, PhKey } from '@phosphor-icons/vue'
import type { AccessLevel, AccessUser, GrantTarget } from '~/types/access'

const accessWeight: Record<AccessLevel, number> = {
  none: 0,
  read: 1,
  write: 2,
  admin: 3,
}

export function useAccess() {
  const { get } = useApiClient()
  const workbench = useWorkbench()
  const { environments } = useEnvironments()

  const users = useState<AccessUser[]>('access:users', () => [])
  const usersLoading = useState<boolean>('access:loading', () => true)
  const selectedUserId = useState<string>('access:selected-user', () => '')

  const selectedUser = computed(() => users.value.find(user => user.id === selectedUserId.value) ?? users.value[0])

  const loadUsers = async () => {
    usersLoading.value = true
    try {
      const data = await get<any[]>('/v1/users')
      users.value = data.map((u: any) => ({
        id: u.id,
        username: u.username,
        email: u.email,
        role: u.role,
        status: u.active ? 'active' : 'inactive',
        inheritedFrom: '',
        grants: { collections: {}, environments: {}, secrets: {} },
      }))
      if (data.length > 0 && !selectedUserId.value) {
        selectedUserId.value = data[0].id
      }
    } catch (err) {
      console.error('Failed to load users', err)
    } finally {
      usersLoading.value = false
    }
  }

  if (usersLoading.value) {
    loadUsers()
  }

  const collectionEntries = computed(() =>
    workbench.treeItems.value.map(group => ({
      name: group.name,
      requests: group.requests.length,
      defaultEnvironment: environments.value[0]?.name ?? '—',
    }))
  )

  const environmentEntries = computed(() =>
    environments.value.map((env: any) => ({
      name: env.name,
      visibility: env.visibility,
      variables: env.variables?.length ?? 0,
      secrets: env.variables?.filter((v: any) => v.secret).length ?? 0,
    }))
  )

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
        rows: collectionEntries.value.map(collection => ({
          name: collection.name,
          meta: `${collection.requests} requests / default ${collection.defaultEnvironment}`,
          level: selectedUser.value!.grants.collections[collection.name] ?? 'none',
        })),
      },
      {
        key: 'environment' as GrantTarget,
        label: 'Environments',
        icon: PhDatabase,
        rows: environmentEntries.value.map(environment => ({
          name: environment.name,
          meta: `${environment.visibility} / ${environment.variables} variables`,
          level: selectedUser.value!.grants.environments[environment.name] ?? 'none',
        })),
      },
      {
        key: 'secret' as GrantTarget,
        label: 'Secrets',
        icon: PhKey,
        rows: environmentEntries.value.map(environment => ({
          name: environment.name,
          meta: `${environment.secrets} masked value${environment.secrets === 1 ? '' : 's'}`,
          level: selectedUser.value!.grants.secrets[environment.name] ?? 'none',
        })),
      },
    ]
  })

  const executionRows = computed(() =>
    collectionEntries.value.flatMap(collection => environmentEntries.value.map(environment => ({
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
    usersLoading,
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
    loadUsers,
  }
}
