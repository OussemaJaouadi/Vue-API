export function useEnvironments() {
  const environments = useState('environments:data', () => [
    {
      name: 'Local',
      visibility: 'project',
      allowedRoles: ['manager', 'developer', 'tester'],
      variables: [
        { key: 'apiBaseUrl', value: 'http://localhost:8080', secret: false },
        { key: 'wsBaseUrl', value: 'ws://localhost:8080', secret: false },
        { key: 'accessToken', value: '••••••••••••••••', secret: true },
      ],
    },
    {
      name: 'Staging',
      visibility: 'restricted',
      allowedRoles: ['manager', 'developer'],
      variables: [
        { key: 'apiBaseUrl', value: 'https://api.staging.example.com', secret: false },
        { key: 'wsBaseUrl', value: 'wss://api.staging.example.com', secret: false },
        { key: 'apiKey', value: '••••••••••••••••', secret: true },
      ],
    },
  ])

  const activeEnvironmentName = useState<string>('environments:active', () => 'Local')
  const activeEnvironment = computed(() => environments.value.find(environment => environment.name === activeEnvironmentName.value) ?? environments.value[0]!)
  const secretVariableCount = computed(() => activeEnvironment.value.variables.filter(variable => variable.secret).length)

  const handleCreate = (env: any) => {
    environments.value.push({
      ...env,
      allowedRoles: ['manager', 'developer'],
      variables: [],
    })
    activeEnvironmentName.value = env.name
  }

  const deleteEnvironment = () => {
    if (environments.value.length <= 1) return
    
    const nameToDelete = activeEnvironmentName.value
    const index = environments.value.findIndex(e => e.name === nameToDelete)
    environments.value = environments.value.filter(e => e.name !== nameToDelete)
    
    const nextEnv = environments.value[Math.max(0, index - 1)]
    if (nextEnv) {
      activeEnvironmentName.value = nextEnv.name
    }
  }

  const addVariable = () => {
    activeEnvironment.value.variables.push({
      key: 'NEW_KEY',
      value: '',
      secret: false,
    })
  }

  return {
    environments,
    activeEnvironmentName,
    activeEnvironment,
    secretVariableCount,
    handleCreate,
    deleteEnvironment,
    addVariable,
  }
}
