import type { AuthUser, ApiResponse } from '~/types'

const currentUser = ref<AuthUser | null>(null)

export function useAuth() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string
  const tokenCookie = useCookie<string | null>('zenwoot_token', {
    maxAge: 60 * 60 * 24,
    path: '/',
    sameSite: 'lax'
  })
  const router = useRouter()
  const toast = useToast()

  const isAuthenticated = computed(() => !!tokenCookie.value)

  async function login(email: string, password: string): Promise<void> {
    const res = await $fetch<ApiResponse<{ token: string, user: AuthUser }>>(`${baseURL}/api/v1/auth/login`, {
      method: 'POST',
      body: { email, password }
    })
    if (!res.success || !res.data?.token) {
      throw new Error(res.error || 'Login failed')
    }
    tokenCookie.value = res.data.token
    currentUser.value = res.data.user
    await router.push('/conversations')
  }

  async function fetchMe(): Promise<void> {
    if (!tokenCookie.value) return
    try {
      const res = await $fetch<ApiResponse<AuthUser>>(`${baseURL}/api/v1/auth/me`, {
        headers: { Authorization: `Bearer ${tokenCookie.value}` }
      })
      if (res.success) currentUser.value = res.data
    } catch {
      tokenCookie.value = null
      currentUser.value = null
    }
  }

  async function logout(): Promise<void> {
    tokenCookie.value = null
    currentUser.value = null
    toast.add({ title: 'Logged out', color: 'neutral' })
    await router.push('/login')
  }

  return {
    currentUser: readonly(currentUser),
    isAuthenticated,
    login,
    logout,
    fetchMe
  }
}
