<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { CannedResponse } from '~/types'

const api = useApi()
const toast = useToast()

const responses = ref<CannedResponse[]>([])
const loading = ref(false)
const search = ref('')
const sortAsc = ref(true)

const createOpen = ref(false)
const editOpen = ref(false)
const deleteOpen = ref(false)
const saving = ref(false)
const editingId = ref<string | null>(null)
const deletingId = ref<string | null>(null)

const form = reactive({ shortCode: '', content: '' })

const columns: TableColumn<CannedResponse>[] = [
  { accessorKey: 'shortCode', header: 'Short Code', enableSorting: true },
  { accessorKey: 'content', header: 'Content' },
  { id: 'actions', header: '' }
]

const filtered = computed(() => {
  const q = search.value.toLowerCase()
  const list = responses.value.filter(r =>
    r.shortCode.toLowerCase().includes(q) || r.content.toLowerCase().includes(q)
  )
  return sortAsc.value
    ? list.sort((a, b) => a.shortCode.localeCompare(b.shortCode))
    : list.sort((a, b) => b.shortCode.localeCompare(a.shortCode))
})

async function load() {
  loading.value = true
  try {
    const data = await api.get<CannedResponse[]>('/api/v1/canned-responses')
    responses.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load canned responses', color: 'error' })
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { shortCode: '', content: '' })
  createOpen.value = true
}

function openEdit(r: CannedResponse) {
  editingId.value = r.id
  Object.assign(form, { shortCode: r.shortCode, content: r.content })
  editOpen.value = true
}

function confirmDelete(id: string) {
  deletingId.value = id
  deleteOpen.value = true
}

async function create() {
  saving.value = true
  try {
    await api.post('/api/v1/canned-responses', form)
    toast.add({ title: 'Canned response created', color: 'success' })
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
    await api.put(`/api/v1/canned-responses/${editingId.value}`, form)
    toast.add({ title: 'Updated', color: 'success' })
    editOpen.value = false
    await load()
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Failed to update', color: 'error' })
  } finally { saving.value = false }
}

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.del(`/api/v1/canned-responses/${deletingId.value}`)
    toast.add({ title: 'Deleted', color: 'success' })
    deleteOpen.value = false
    await load()
  } catch { toast.add({ title: 'Failed to delete', color: 'error' }) }
}

onMounted(load)
</script>

<template>
  <UDashboardPanel id="settings-canned">
    <template #header>
      <UDashboardNavbar title="Canned Responses">
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
          <UButton
            :icon="sortAsc ? 'i-lucide-arrow-up-a-z' : 'i-lucide-arrow-down-z-a'"
            size="sm"
            variant="ghost"
            color="neutral"
            @click="sortAsc = !sortAsc"
          />
          <UButton icon="i-lucide-plus" size="sm" @click="openCreate">
            New Response
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
        <template #shortCode-cell="{ row }">
          <code class="text-sm font-mono text-primary bg-primary/5 px-1.5 py-0.5 rounded">/{{ row.original.shortCode }}</code>
        </template>

        <template #content-cell="{ row }">
          <p class="text-sm text-muted truncate max-w-md">
            {{ row.original.content }}
          </p>
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
  <UModal v-model:open="createOpen" title="New Canned Response">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Short Code" name="shortCode">
          <UInput v-model="form.shortCode" placeholder="e.g. greeting" class="w-full" />
        </UFormField>
        <UFormField label="Content" name="content">
          <UTextarea
            v-model="form.content"
            placeholder="Hello! How can I help you today?"
            :rows="4"
            class="w-full"
          />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="createOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.shortCode || !form.content" @click="create">
        Create
      </UButton>
    </template>
  </UModal>

  <!-- Edit modal -->
  <UModal v-model:open="editOpen" title="Edit Canned Response">
    <template #body>
      <div class="space-y-4">
        <UFormField label="Short Code" name="shortCode">
          <UInput v-model="form.shortCode" class="w-full" />
        </UFormField>
        <UFormField label="Content" name="content">
          <UTextarea v-model="form.content" :rows="4" class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="editOpen = false">
        Cancel
      </UButton>
      <UButton :loading="saving" :disabled="!form.shortCode || !form.content" @click="update">
        Save
      </UButton>
    </template>
  </UModal>

  <!-- Delete confirm -->
  <UModal v-model:open="deleteOpen" title="Delete Canned Response">
    <template #body>
      <p class="text-sm text-default">
        Are you sure you want to delete this canned response?
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
