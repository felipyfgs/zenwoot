<script setup lang="ts">
import type { User, UserRole } from '~/types'

const api = useApi()
const toast = useToast()

const agents = ref<User[]>([])
const loading = ref(false)
const search = ref('')

const createOpen = ref(false)
const editOpen = ref(false)
const deleteOpen = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const deletingId = ref<string | null>(null)

const form = reactive({ email: '', name: '', displayName: '', role: 'agent' as UserRole, password: '' })

const roleOptions = [
  { label: 'Agent', value: 'agent' },
  { label: 'Admin', value: 'admin' }
]

const filtered = computed(() =>
  agents.value.filter(a =>
    a.name.toLowerCase().includes(search.value.toLowerCase())
    || a.email.toLowerCase().includes(search.value.toLowerCase())
  )
)

async function load() {
  loading.value = true
  try {
    const data = await api.get<User[]>('/api/v1/users')
    agents.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load agents', color: 'error' })
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { email: '', name: '', displayName: '', role: 'agent', password: '' })
  createOpen.value = true
}

function openEdit(agent: User) {
  editingId.value = agent.id
  Object.assign(form, { email: agent.email, name: agent.name, displayName: agent.displayName || '', role: agent.role, password: '' })
  editOpen.value = true
}

function confirmDelete(id: string) {
  deletingId.value = id
  deleteOpen.value = true
}

async function create() {
  saving.value = true
  try {
    await api.post('/api/v1/users', form)
    toast.add({ title: 'Agent created', color: 'success' })
    createOpen.value = false
    await load()
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Failed to create', color: 'error' })
  } finally { saving.value = false }
}

async function update() {
  if (!editingId.value) return
  saving.value = true
  try {
    const payload: Record<string, string> = { name: form.name, displayName: form.displayName, role: form.role }
    if (form.password) payload.password = form.password
    await api.put(`/api/v1/users/${editingId.value}`, payload)
    toast.add({ title: 'Agent updated', color: 'success' })
    editOpen.value = false
    await load()
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Failed to update', color: 'error' })
  } finally { saving.value = false }
}

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.del(`/api/v1/users/${deletingId.value}`)
    toast.add({ title: 'Agent deleted', color: 'success' })
    deleteOpen.value = false
    await load()
  } catch { toast.add({ title: 'Failed to delete', color: 'error' }) }
}

function roleColor(role: string) {
  if (role === 'admin' || role === 'super_admin') return 'primary'
  return 'neutral'
}

onMounted(load)
</script>

<template>
  <UDashboardPanel id="settings-agents">
    <template #header>
      <UDashboardNavbar title="Agents">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
        <template #right>
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Search…"
            size="sm"
            class="w-40"
          />
          <UButton icon="i-lucide-plus" size="sm" @click="openCreate">
            New Agent
          </UButton>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div v-if="loading" class="flex justify-center py-12">
        <UIcon name="i-lucide-loader-circle" class="size-6 animate-spin text-muted" />
      </div>
      <p v-else-if="!filtered.length" class="text-center text-muted text-sm py-12">
        No agents found
      </p>
      <ul v-else class="divide-y divide-default">
        <li v-for="agent in filtered" :key="agent.id" class="flex items-center gap-4 px-4 py-3">
          <UAvatar
            :alt="agent.displayName || agent.name"
            :src="agent.avatarUrl"
            size="sm"
            class="shrink-0"
          />
          <div class="flex-1 min-w-0">
            <p class="font-medium text-highlighted truncate">
              {{ agent.displayName || agent.name }}
            </p>
            <p class="text-xs text-muted truncate">
              {{ agent.email }}
            </p>
          </div>
          <UBadge
            :color="roleColor(agent.role)"
            variant="subtle"
            size="xs"
            class="capitalize"
          >
            {{ agent.role }}
          </UBadge>
          <div class="flex items-center gap-1">
            <UButton
              size="xs"
              variant="ghost"
              icon="i-lucide-pencil"
              @click="openEdit(agent)"
            />
            <UButton
              size="xs"
              variant="ghost"
              color="error"
              icon="i-lucide-trash-2"
              @click="confirmDelete(agent.id)"
            />
          </div>
        </li>
      </ul>
    </template>
  </UDashboardPanel>

  <!-- Create modal -->
  <UModal v-model:open="createOpen" title="New Agent">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Email" name="email">
          <UInput
            v-model="form.email"
            type="email"
            placeholder="agent@company.com"
            class="w-full"
          />
        </UFormField>
        <UFormField label="Name" name="name">
          <UInput v-model="form.name" placeholder="Full name" class="w-full" />
        </UFormField>
        <UFormField label="Display Name" name="displayName">
          <UInput v-model="form.displayName" placeholder="Optional" class="w-full" />
        </UFormField>
        <UFormField label="Role" name="role">
          <USelect
            v-model="form.role"
            :items="roleOptions"
            value-key="value"
            class="w-full"
          />
        </UFormField>
        <UFormField label="Password" name="password">
          <UInput
            v-model="form.password"
            type="password"
            placeholder="••••••••"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="createOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.email || !form.name || !form.password" @click="create">
        Create
      </UButton>
    </template>
  </UModal>

  <!-- Edit modal -->
  <UModal v-model:open="editOpen" title="Edit Agent">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Name" name="name">
          <UInput v-model="form.name" placeholder="Full name" class="w-full" />
        </UFormField>
        <UFormField label="Display Name" name="displayName">
          <UInput v-model="form.displayName" placeholder="Optional" class="w-full" />
        </UFormField>
        <UFormField label="Role" name="role">
          <USelect
            v-model="form.role"
            :items="roleOptions"
            value-key="value"
            class="w-full"
          />
        </UFormField>
        <UFormField label="New Password" name="password">
          <UInput
            v-model="form.password"
            type="password"
            placeholder="Leave blank to keep current"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="editOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.name" @click="update">
        Save
      </UButton>
    </template>
  </UModal>

  <!-- Delete confirm -->
  <UModal v-model:open="deleteOpen" title="Delete Agent">
    <template #body>
      <p class="text-sm text-default">
        Are you sure you want to delete this agent?
      </p>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="deleteOpen = false">
        Cancel
      </UButton>
      <UButton color="error" @click="doDelete">
        Delete
      </UButton>
    </template>
  </UModal>
</template>
