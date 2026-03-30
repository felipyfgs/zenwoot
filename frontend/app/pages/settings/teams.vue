<script setup lang="ts">
import type { Team, User } from '~/types'

const api = useApi()
const toast = useToast()

const teams = ref<Team[]>([])
const loading = ref(false)
const search = ref('')

const createOpen = ref(false)
const editOpen = ref(false)
const deleteOpen = ref(false)
const membersOpen = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const deletingId = ref<string | null>(null)
const selectedTeam = ref<Team | null>(null)

const form = reactive({ name: '', description: '', allowAutoAssign: false })

const teamMembers = ref<User[]>([])
const allAgents = ref<User[]>([])
const loadingMembers = ref(false)
const addingMember = ref<string | null>(null)
const removingMember = ref<string | null>(null)

const filtered = computed(() =>
  teams.value.filter(t => t.name.toLowerCase().includes(search.value.toLowerCase()))
)

const nonMembers = computed(() =>
  allAgents.value.filter(a => !teamMembers.value.some(m => m.id === a.id))
)

async function load() {
  loading.value = true
  try {
    const data = await api.get<Team[]>('/api/v1/teams')
    teams.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load teams', color: 'error' })
  } finally {
    loading.value = false
  }
}

async function loadAgents() {
  try {
    const data = await api.get<User[]>('/api/v1/users')
    allAgents.value = Array.isArray(data) ? data : []
  } catch { /* silent */ }
}

async function openMembers(team: Team) {
  selectedTeam.value = team
  membersOpen.value = true
  loadingMembers.value = true
  try {
    const data = await api.get<User[]>(`/api/v1/teams/${team.id}/members`)
    teamMembers.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load members', color: 'error' })
  } finally {
    loadingMembers.value = false
  }
}

async function addMember(userId: string) {
  if (!selectedTeam.value) return
  addingMember.value = userId
  try {
    await api.post(`/api/v1/teams/${selectedTeam.value.id}/members/${userId}`)
    const agent = allAgents.value.find(a => a.id === userId)
    if (agent) teamMembers.value.push(agent)
  } catch {
    toast.add({ title: 'Failed to add member', color: 'error' })
  } finally {
    addingMember.value = null
  }
}

async function removeMember(userId: string) {
  if (!selectedTeam.value) return
  removingMember.value = userId
  try {
    await api.del(`/api/v1/teams/${selectedTeam.value.id}/members/${userId}`)
    teamMembers.value = teamMembers.value.filter(m => m.id !== userId)
  } catch {
    toast.add({ title: 'Failed to remove member', color: 'error' })
  } finally {
    removingMember.value = null
  }
}

function openCreate() {
  Object.assign(form, { name: '', description: '', allowAutoAssign: false })
  createOpen.value = true
}

function openEdit(team: Team) {
  editingId.value = team.id
  Object.assign(form, { name: team.name, description: team.description || '', allowAutoAssign: team.allowAutoAssign })
  editOpen.value = true
}

function confirmDelete(id: string) {
  deletingId.value = id
  deleteOpen.value = true
}

async function create() {
  saving.value = true
  try {
    await api.post('/api/v1/teams', form)
    toast.add({ title: 'Team created', color: 'success' })
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
    await api.put(`/api/v1/teams/${editingId.value}`, form)
    toast.add({ title: 'Team updated', color: 'success' })
    editOpen.value = false
    await load()
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Failed to update', color: 'error' })
  } finally { saving.value = false }
}

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.del(`/api/v1/teams/${deletingId.value}`)
    toast.add({ title: 'Team deleted', color: 'success' })
    deleteOpen.value = false
    await load()
  } catch { toast.add({ title: 'Failed to delete', color: 'error' }) }
}

onMounted(async () => {
  await Promise.all([load(), loadAgents()])
})
</script>

<template>
  <UDashboardPanel id="settings-teams">
    <template #header>
      <UDashboardNavbar title="Teams">
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
            New Team
          </UButton>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div v-if="loading" class="flex justify-center py-12">
        <UIcon name="i-lucide-loader-circle" class="size-6 animate-spin text-muted" />
      </div>
      <p v-else-if="!filtered.length" class="text-center text-muted text-sm py-12">
        No teams found
      </p>
      <ul v-else class="divide-y divide-default">
        <li v-for="team in filtered" :key="team.id" class="flex items-center gap-4 px-4 py-3">
          <div class="size-8 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
            <UIcon name="i-lucide-users-2" class="size-4 text-primary" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="font-medium text-highlighted truncate">
              {{ team.name }}
            </p>
            <p v-if="team.description" class="text-xs text-muted truncate">
              {{ team.description }}
            </p>
          </div>
          <UBadge
            v-if="team.allowAutoAssign"
            color="info"
            variant="subtle"
            size="xs"
          >
            Auto-assign
          </UBadge>
          <div class="flex items-center gap-1">
            <UButton
              size="xs"
              variant="ghost"
              icon="i-lucide-users"
              @click="openMembers(team)"
            >
              Members
            </UButton>
            <UButton
              size="xs"
              variant="ghost"
              icon="i-lucide-pencil"
              @click="openEdit(team)"
            />
            <UButton
              size="xs"
              variant="ghost"
              color="error"
              icon="i-lucide-trash-2"
              @click="confirmDelete(team.id)"
            />
          </div>
        </li>
      </ul>
    </template>
  </UDashboardPanel>

  <!-- Create modal -->
  <UModal v-model:open="createOpen" title="New Team">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Name" name="name">
          <UInput v-model="form.name" placeholder="Support Team" class="w-full" />
        </UFormField>
        <UFormField label="Description" name="description">
          <UInput v-model="form.description" placeholder="Optional" class="w-full" />
        </UFormField>
        <UFormField label="Auto-assign conversations" name="allowAutoAssign">
          <USwitch v-model="form.allowAutoAssign" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="createOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.name" @click="create">
        Create
      </UButton>
    </template>
  </UModal>

  <!-- Edit modal -->
  <UModal v-model:open="editOpen" title="Edit Team">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Name" name="name">
          <UInput v-model="form.name" class="w-full" />
        </UFormField>
        <UFormField label="Description" name="description">
          <UInput v-model="form.description" class="w-full" />
        </UFormField>
        <UFormField label="Auto-assign conversations" name="allowAutoAssign">
          <USwitch v-model="form.allowAutoAssign" />
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
  <UModal v-model:open="deleteOpen" title="Delete Team">
    <template #body>
      <p class="text-sm text-default">
        Are you sure you want to delete this team?
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

  <!-- Members slideover -->
  <USlideover v-model:open="membersOpen" title="Team Members" class="max-w-sm">
    <template #body>
      <div class="space-y-4">
        <div v-if="loadingMembers" class="flex justify-center py-6">
          <UIcon name="i-lucide-loader-circle" class="size-5 animate-spin text-muted" />
        </div>
        <template v-else>
          <div>
            <p class="text-xs font-semibold text-muted uppercase tracking-wider mb-2">
              Current Members
            </p>
            <p v-if="!teamMembers.length" class="text-sm text-muted">
              No members yet
            </p>
            <ul v-else class="divide-y divide-default">
              <li v-for="member in teamMembers" :key="member.id" class="flex items-center gap-3 py-2">
                <UAvatar :alt="member.displayName || member.name" size="xs" />
                <span class="flex-1 text-sm text-default truncate">{{ member.displayName || member.name }}</span>
                <UButton
                  size="xs"
                  variant="ghost"
                  color="error"
                  icon="i-lucide-x"
                  :loading="removingMember === member.id"
                  @click="removeMember(member.id)"
                />
              </li>
            </ul>
          </div>
          <div v-if="nonMembers.length">
            <p class="text-xs font-semibold text-muted uppercase tracking-wider mb-2">
              Add Members
            </p>
            <ul class="divide-y divide-default">
              <li v-for="agent in nonMembers" :key="agent.id" class="flex items-center gap-3 py-2">
                <UAvatar :alt="agent.displayName || agent.name" size="xs" />
                <span class="flex-1 text-sm text-default truncate">{{ agent.displayName || agent.name }}</span>
                <UButton
                  size="xs"
                  variant="ghost"
                  color="primary"
                  icon="i-lucide-plus"
                  :loading="addingMember === agent.id"
                  @click="addMember(agent.id)"
                />
              </li>
            </ul>
          </div>
        </template>
      </div>
    </template>
  </USlideover>
</template>
