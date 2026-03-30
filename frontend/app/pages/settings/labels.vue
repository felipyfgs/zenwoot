<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Label } from '~/types'

const api = useApi()
const toast = useToast()

const labels = ref<Label[]>([])
const loading = ref(false)
const search = ref('')

const createOpen = ref(false)
const editOpen = ref(false)
const deleteOpen = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const deletingId = ref<string | null>(null)

const form = reactive({ title: '', description: '', color: '#6366f1', showOnSidebar: true })

const columns: TableColumn<Label>[] = [
  { accessorKey: 'title', header: 'Title', enableSorting: true },
  { accessorKey: 'description', header: 'Description' },
  { accessorKey: 'color', header: 'Color' },
  { id: 'actions', header: '' }
]

const filtered = computed(() =>
  labels.value.filter(l => l.title.toLowerCase().includes(search.value.toLowerCase()))
)

async function load() {
  loading.value = true
  try {
    const data = await api.get<Label[]>('/api/v1/labels')
    labels.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load labels', color: 'error' })
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { title: '', description: '', color: '#6366f1', showOnSidebar: true })
  createOpen.value = true
}

function openEdit(label: Label) {
  editingId.value = label.id
  Object.assign(form, { title: label.title, description: label.description || '', color: label.color, showOnSidebar: label.showOnSidebar })
  editOpen.value = true
}

function confirmDelete(id: string) {
  deletingId.value = id
  deleteOpen.value = true
}

async function create() {
  saving.value = true
  try {
    await api.post('/api/v1/labels', form)
    toast.add({ title: 'Label created', color: 'success' })
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
    await api.put(`/api/v1/labels/${editingId.value}`, form)
    toast.add({ title: 'Label updated', color: 'success' })
    editOpen.value = false
    await load()
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Failed to update', color: 'error' })
  } finally { saving.value = false }
}

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.del(`/api/v1/labels/${deletingId.value}`)
    toast.add({ title: 'Label deleted', color: 'success' })
    deleteOpen.value = false
    await load()
  } catch { toast.add({ title: 'Failed to delete', color: 'error' }) }
}

onMounted(load)
</script>

<template>
  <UDashboardPanel id="settings-labels">
    <template #header>
      <UDashboardNavbar title="Labels">
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
            New Label
          </UButton>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <UTable
        :data="filtered"
        :columns="columns"
        :loading="loading"
        class="w-full"
      >
        <template #title-cell="{ row }">
          <div class="flex items-center gap-2">
            <span
              class="size-3 rounded-full shrink-0 ring-1 ring-default"
              :style="{ background: row.original.color }"
            />
            <span class="font-medium text-highlighted">{{ row.original.title }}</span>
          </div>
        </template>

        <template #description-cell="{ row }">
          <span class="text-muted text-sm">{{ row.original.description || '—' }}</span>
        </template>

        <template #color-cell="{ row }">
          <span class="font-mono text-xs text-muted">{{ row.original.color }}</span>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center justify-end gap-1">
            <UButton
              size="xs"
              variant="ghost"
              icon="i-lucide-pencil"
              @click="openEdit(row.original)"
            />
            <UButton
              size="xs"
              variant="ghost"
              color="error"
              icon="i-lucide-trash-2"
              @click="confirmDelete(row.original.id)"
            />
          </div>
        </template>
      </UTable>
    </template>
  </UDashboardPanel>

  <!-- Create modal -->
  <UModal v-model:open="createOpen" title="New Label">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Title" name="title">
          <UInput v-model="form.title" placeholder="e.g. bug, feature-request" class="w-full" />
        </UFormField>
        <UFormField label="Description" name="description">
          <UInput v-model="form.description" placeholder="Optional" class="w-full" />
        </UFormField>
        <UFormField label="Color" name="color">
          <div class="flex items-center gap-2">
            <input v-model="form.color" type="color" class="size-8 rounded cursor-pointer border border-default">
            <UInput v-model="form.color" placeholder="#6366f1" class="flex-1 font-mono" />
          </div>
        </UFormField>
        <UFormField label="Show on sidebar" name="showOnSidebar">
          <USwitch v-model="form.showOnSidebar" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="createOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.title" @click="create">
        Create
      </UButton>
    </template>
  </UModal>

  <!-- Edit modal -->
  <UModal v-model:open="editOpen" title="Edit Label">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Title" name="title">
          <UInput v-model="form.title" class="w-full" />
        </UFormField>
        <UFormField label="Description" name="description">
          <UInput v-model="form.description" class="w-full" />
        </UFormField>
        <UFormField label="Color" name="color">
          <div class="flex items-center gap-2">
            <input v-model="form.color" type="color" class="size-8 rounded cursor-pointer border border-default">
            <UInput v-model="form.color" class="flex-1 font-mono" />
          </div>
        </UFormField>
        <UFormField label="Show on sidebar" name="showOnSidebar">
          <USwitch v-model="form.showOnSidebar" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="editOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.title" @click="update">
        Save
      </UButton>
    </template>
  </UModal>

  <!-- Delete confirm -->
  <UModal v-model:open="deleteOpen" title="Delete Label">
    <template #body>
      <p class="text-sm text-default">
        Are you sure you want to delete this label?
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
