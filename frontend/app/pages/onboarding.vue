<script setup lang="ts">
import type { ApiResponse, AuthUser } from '~/types'

definePageMeta({ layout: false })

const config = useRuntimeConfig()
const baseURL = config.public.apiBase as string
const tokenCookie = useCookie<string | null>('zenwoot_token', { maxAge: 60 * 60 * 24, path: '/', sameSite: 'lax' })
const router = useRouter()
const toast = useToast()

const setupChecked = useState('setup:checked', () => false)
const needsSetup = useState('setup:needs', () => false)

const step = ref<'welcome' | 'account' | 'done'>('welcome')
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const passwordMismatch = computed(
  () => form.confirmPassword.length > 0 && form.password !== form.confirmPassword
)
const canSubmit = computed(
  () => form.name.trim() && form.email.trim() && form.password.length >= 8 && !passwordMismatch.value
)

async function handleSetup() {
  errorMsg.value = ''
  loading.value = true
  try {
    const res = await $fetch<ApiResponse<{ token: string, user: AuthUser }>>(`${baseURL}/api/v1/setup`, {
      method: 'POST',
      body: { email: form.email, name: form.name, password: form.password }
    })
    if (!res.success || !res.data?.token) {
      throw new Error(res.error || 'Setup failed')
    }
    tokenCookie.value = res.data.token
    needsSetup.value = false
    setupChecked.value = true
    step.value = 'done'
    setTimeout(() => router.push('/conversations'), 1800)
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : 'Setup failed'
    toast.add({ title: errorMsg.value, color: 'error' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-default flex flex-col items-center justify-center px-4">
    <!-- Logo / Brand -->
    <div class="flex flex-col items-center mb-8 space-y-2">
      <div class="size-14 rounded-2xl bg-primary flex items-center justify-center shadow-lg">
        <UIcon name="i-lucide-zap" class="size-7 text-white" />
      </div>
      <h1 class="text-3xl font-bold text-highlighted">
        Zenwoot
      </h1>
      <p class="text-muted text-sm">
        WhatsApp customer support platform
      </p>
    </div>

    <!-- Step: Welcome -->
    <Transition name="fade" mode="out-in">
      <div v-if="step === 'welcome'" class="w-full max-w-md">
        <UCard class="shadow-lg">
          <div class="text-center space-y-3 py-4">
            <div class="size-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto">
              <UIcon name="i-lucide-party-popper" class="size-6 text-primary" />
            </div>
            <h2 class="text-xl font-semibold text-highlighted">
              Welcome to Zenwoot!
            </h2>
            <p class="text-sm text-muted leading-relaxed">
              Looks like this is a fresh installation.<br>
              Let's create your admin account to get started.
            </p>
          </div>
          <template #footer>
            <UButton block size="lg" @click="step = 'account'">
              Get started
              <template #trailing>
                <UIcon name="i-lucide-arrow-right" class="size-4" />
              </template>
            </UButton>
          </template>
        </UCard>
      </div>

      <!-- Step: Account creation -->
      <div v-else-if="step === 'account'" class="w-full max-w-md">
        <UCard class="shadow-lg">
          <template #header>
            <div class="flex items-center gap-2">
              <UButton
                icon="i-lucide-arrow-left"
                variant="ghost"
                color="neutral"
                size="xs"
                square
                @click="step = 'welcome'"
              />
              <div>
                <h2 class="text-base font-semibold text-highlighted">
                  Create admin account
                </h2>
                <p class="text-xs text-muted">
                  This will be the primary administrator
                </p>
              </div>
            </div>
          </template>

          <form class="space-y-4" @submit.prevent="handleSetup">
            <UFormField label="Your name" name="name">
              <UInput
                v-model="form.name"
                placeholder="John Doe"
                autocomplete="name"
                class="w-full"
                :disabled="loading"
              />
            </UFormField>

            <UFormField label="Email" name="email">
              <UInput
                v-model="form.email"
                type="email"
                placeholder="admin@company.com"
                autocomplete="email"
                class="w-full"
                :disabled="loading"
              />
            </UFormField>

            <UFormField label="Password" name="password" :hint="form.password.length > 0 && form.password.length < 8 ? 'Minimum 8 characters' : undefined">
              <UInput
                v-model="form.password"
                type="password"
                placeholder="••••••••"
                autocomplete="new-password"
                class="w-full"
                :disabled="loading"
              />
            </UFormField>

            <UFormField
              label="Confirm password"
              name="confirmPassword"
              :error="passwordMismatch ? 'Passwords do not match' : undefined"
            >
              <UInput
                v-model="form.confirmPassword"
                type="password"
                placeholder="••••••••"
                autocomplete="new-password"
                class="w-full"
                :disabled="loading"
                :color="passwordMismatch ? 'error' : undefined"
              />
            </UFormField>

            <p v-if="errorMsg" class="text-sm text-error text-center">
              {{ errorMsg }}
            </p>
          </form>

          <template #footer>
            <UButton
              block
              size="lg"
              :loading="loading"
              :disabled="!canSubmit"
              @click="handleSetup"
            >
              Create account & continue
            </UButton>
          </template>
        </UCard>
      </div>

      <!-- Step: Done -->
      <div v-else-if="step === 'done'" class="w-full max-w-md">
        <UCard class="shadow-lg">
          <div class="text-center space-y-3 py-6">
            <div class="size-14 rounded-full bg-success/10 flex items-center justify-center mx-auto">
              <UIcon name="i-lucide-check-circle" class="size-8 text-success" />
            </div>
            <h2 class="text-xl font-semibold text-highlighted">
              All set!
            </h2>
            <p class="text-sm text-muted">
              Your admin account has been created.<br>
              Redirecting to the dashboard…
            </p>
            <UIcon name="i-lucide-loader-circle" class="size-5 animate-spin text-muted mx-auto" />
          </div>
        </UCard>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from {
  opacity: 0;
  transform: translateX(12px);
}
.fade-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
</style>
