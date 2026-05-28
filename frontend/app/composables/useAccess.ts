import { PhDatabase, PhFolderOpen, PhKey } from '@phosphor-icons/vue'
import type { AccessLevel, AccessUser, GrantTarget } from '~/types/access'

const accessWeight: Record<AccessLevel, number> = {
  none: 0,
  read: 1,
  write: 2,
  admin: 3,
}

export function useAccess() {
  const { get, put, delete: del } = useApiClient()
  const workbench = useWorkbench()
  const { environments } = useEnvironments()

  const users = useState<AccessUser[]>('access:users', () => [])
  const usersLoading = useState<boolean>('access:loading', () => true)
  const usersError = useState<string | null>('access:error', () => null)
  const selectedUserId = useState<string>('access:selected-user', () => '')
  const usersWorkspaceId = useState<string>('access:workspace-id', () => '')

  const selectedUser = computed(() => users.value.find(user => user.id === selectedUserId.value) ?? users.value[0])

  const workspaceId = useState<string>('workspace:id', () => '')
  const workspacesLoading = useState<boolean>('workspaces:loading', () => false)

  const serializeGrants = (user: AccessUser) => {
    const entries: Array<{ resourceType: GrantTarget, grants: Record<string, AccessLevel> }> = [
      { resourceType: 'collection', grants: user.grants.collections },
      { resourceType: 'environment', grants: user.grants.environments },
      { resourceType: 'secret', grants: user.grants.secrets },
    ]

    return entries.flatMap(entry =>
      Object.entries(entry.grants)
        .filter(([, accessLevel]) => accessLevel !== 'none')
        .map(([resourceId, accessLevel]) => ({
          resourceType: entry.resourceType,
          resourceId,
          accessLevel,
        })),
    )
  }

  const persistGrants = async (user: AccessUser) => {
    const wid = workspaceId.value
    if (!wid) return

    await put(`/v1/workspaces/${wid}/members/${user.id}/grants`, {
      grants: serializeGrants(user),
    })
  }

  const loadGrantsForUser = async (userId: string) => {
    const wid = workspaceId.value
    if (!wid) return { collections: {}, environments: {}, secrets: {} }
    try {
      return await get<any>(`/v1/workspaces/${wid}/members/${userId}/grants`)
    } catch {
      return { collections: {}, environments: {}, secrets: {} }
    }
  }

  const loadUsers = async () => {
    const wid = workspaceId.value
    if (!wid) {
      if (!workspacesLoading.value) usersLoading.value = false
      return
    }
    usersLoading.value = true
    usersError.value = null
    try {
      const data = await get<any[]>(`/v1/workspaces/${wid}/members`)
      users.value = await Promise.all(data.map(async (m: any) => {
        const grants = await loadGrantsForUser(m.userId)
        return {
          id: m.userId,
          username: m.username,
          email: m.email,
          role: m.role,
          status: 'active',
          inheritedFrom: '',
          grants,
        }
      }))
      if (data.length > 0 && !selectedUserId.value) {
        selectedUserId.value = data[0].userId
      }
      usersWorkspaceId.value = wid
    } catch (err: any) {
      users.value = []
      selectedUserId.value = ''
      usersWorkspaceId.value = ''
      usersError.value = err?.data?.error || err?.message || 'Failed to load workspace members'
    } finally {
      usersLoading.value = false
    }
  }

  const settleUsersState = () => {
    if (workspacesLoading.value) {
      usersLoading.value = true
      usersError.value = null
      return
    }

    if (!workspaceId.value) {
      users.value = []
      selectedUserId.value = ''
      usersWorkspaceId.value = ''
      usersLoading.value = false
      return
    }

    if (usersWorkspaceId.value !== workspaceId.value) {
      loadUsers()
      return
    }

    usersLoading.value = false
  }

  watch([workspaceId, workspacesLoading], settleUsersState, { immediate: true })

  const collectionEntries = computed(() =>
    workbench.treeItems.value.map(group => ({
      id: group.id ?? group.name,
      name: group.name,
      requests: group.requests.length,
      defaultEnvironment: environments.value[0]?.name ?? '—',
    }))
  )

  const environmentEntries = computed(() =>
    environments.value.map((env: any) => ({
      id: env.id,
      name: env.name,
      visibility: env.visibility,
      variables: env.variables?.length ?? 0,
      secrets: env.variables?.filter((v: any) => v.secret).length ?? 0,
    }))
  )

  const updateRole = async (role: string) => {
    if (!selectedUser.value) return
    const wid = workspaceId.value
    if (!wid) return

    const previousRole = selectedUser.value.role
    selectedUser.value.role = role
    selectedUser.value.inheritedFrom = `project ${role}`

    try {
      await put(`/v1/workspaces/${wid}/members/${selectedUser.value.id}`, { role })
    } catch (err) {
      selectedUser.value.role = previousRole
      selectedUser.value.inheritedFrom = ''
      console.error('Failed to update member role', err)
    }
  }

  const kickUser = async (id: string) => {
    const wid = workspaceId.value
    if (!wid) return

    try {
      await del(`/v1/workspaces/${wid}/members/${id}`)
    } catch (err) {
      console.error('Failed to remove workspace member', err)
      return
    }

    users.value = users.value.filter(user => user.id !== id)
    if (selectedUserId.value === id) {
      selectedUserId.value = users.value[0]?.id ?? ''
    }
  }

  const updateGrant = async (target: GrantTarget, id: string, level: AccessLevel) => {
    if (!selectedUser.value) return
    const user = selectedUser.value
    let grants: Record<string, AccessLevel>
    if (target === 'collection') {
      grants = user.grants.collections
    }
    else if (target === 'environment') {
      grants = user.grants.environments
    }
    else if (target === 'secret') {
      grants = user.grants.secrets
    }
    else {
      return
    }

    const previousLevel = grants[id]
    grants[id] = level

    try {
      await persistGrants(user)
    } catch (err) {
      if (previousLevel) {
        grants[id] = previousLevel
      } else {
        delete grants[id]
      }
      console.error('Failed to update grants', err)
    }
  }

  const resolveDenied = (target: GrantTarget, id: string) => {
    updateGrant(target, id, 'read')
  }

  const canExecute = (collectionId: string, environmentId: string) => {
    if (!selectedUser.value) return false
    const collectionLevel = selectedUser.value.grants.collections[collectionId] ?? 'none'
    const environmentLevel = selectedUser.value.grants.environments[environmentId] ?? 'none'
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
          id: collection.id,
          name: collection.name,
          meta: `${collection.requests} requests / default ${collection.defaultEnvironment}`,
          level: selectedUser.value!.grants.collections[collection.id] ?? 'none',
        })),
      },
      {
        key: 'environment' as GrantTarget,
        label: 'Environments',
        icon: PhDatabase,
        rows: environmentEntries.value.map(environment => ({
          id: environment.id,
          name: environment.name,
          meta: `${environment.visibility} / ${environment.variables} variables`,
          level: selectedUser.value!.grants.environments[environment.id] ?? 'none',
        })),
      },
      {
        key: 'secret' as GrantTarget,
        label: 'Secrets',
        icon: PhKey,
        rows: environmentEntries.value.map(environment => ({
          id: environment.id,
          name: environment.name,
          meta: `${environment.secrets} masked value${environment.secrets === 1 ? '' : 's'}`,
          level: selectedUser.value!.grants.secrets[environment.id] ?? 'none',
        })),
      },
    ]
  })

  const executionRows = computed(() =>
    collectionEntries.value.flatMap(collection => environmentEntries.value.map(environment => ({
      name: `${collection.name} -> ${environment.name}`,
      meta: canExecute(collection.id, environment.id) ? 'Access match' : 'Missing grant',
      level: (canExecute(collection.id, environment.id) ? 'read' : 'none') as AccessLevel,
    }))),
  )

  const deniedTargets = computed(() => {
    const targets: Array<{ target: GrantTarget, id: string, section: string, name: string, level: AccessLevel }> = []
    grantSections.value.forEach((section) => {
      section.rows.forEach((row) => {
        if (row.level === 'none') {
          targets.push({ target: section.key, id: row.id, section: section.label, name: row.name, level: row.level })
        }
      })
    })
    return targets
  })

  return {
    users,
    usersLoading,
    usersError,
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
