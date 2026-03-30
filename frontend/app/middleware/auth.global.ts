export default defineNuxtRouteMiddleware(async (to) => {
  const publicRoutes = ['/login', '/onboarding']
  const setupChecked = useState('setup:checked', () => false)
  const needsSetup = useState('setup:needs', () => false)

  if (!setupChecked.value) {
    const config = useRuntimeConfig()
    const baseURL = config.public.apiBase as string
    try {
      const res = await $fetch<{ success: boolean, data: { needsSetup: boolean } }>(`${baseURL}/setup`)
      needsSetup.value = res?.data?.needsSetup ?? false
    } catch { /* backend unreachable — skip */ } finally {
      setupChecked.value = true
    }
  }

  if (needsSetup.value && to.path !== '/onboarding') {
    return navigateTo('/onboarding')
  }
  if (!needsSetup.value && to.path === '/onboarding') {
    return navigateTo('/login')
  }

  const token = useCookie('zenwoot_token')

  if (!token.value && !publicRoutes.includes(to.path)) {
    return navigateTo('/login')
  }

  if (token.value && to.path === '/login') {
    return navigateTo('/conversations')
  }
})
