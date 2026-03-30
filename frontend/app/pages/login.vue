<script setup lang="ts">
definePageMeta({ layout: false })

const { login } = useAuth()
const toast = useToast()

const form = reactive({ email: '', password: '' })
const loading = ref(false)
const errorMsg = ref('')

async function handleSubmit() {
  errorMsg.value = ''
  loading.value = true
  try {
    await login(form.email, form.password)
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : 'Login failed'
    toast.add({ title: errorMsg.value, color: 'error' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-default">
    <div class="w-full max-w-sm space-y-6 px-4">
      <div class="text-center space-y-1">
        <h1 class="text-2xl font-bold text-highlighted">
          Zenwoot
        </h1>
        <p class="text-sm text-muted">
          Sign in to your account
        </p>
      </div>

      <UCard class="shadow-md">
        <form class="space-y-4" @submit.prevent="handleSubmit">
          <UFormField label="Email" name="email">
            <UInput
              v-model="form.email"
              type="email"
              placeholder="you@company.com"
              autocomplete="email"
              class="w-full"
              :disabled="loading"
            />
          </UFormField>

          <UFormField label="Password" name="password">
            <UInput
              v-model="form.password"
              type="password"
              placeholder="••••••••"
              autocomplete="current-password"
              class="w-full"
              :disabled="loading"
            />
          </UFormField>

          <p v-if="errorMsg" class="text-sm text-error text-center">
            {{ errorMsg }}
          </p>

          <UButton
            type="submit"
            block
            :loading="loading"
            :disabled="!form.email || !form.password"
          >
            Sign in
          </UButton>
        </form>
      </UCard>
    </div>
  </div>
</template>
