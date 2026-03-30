<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { Contact } from '~/types'

const { list } = useContacts()
const toast = useToast()

const search = ref('')
const page = ref(1)
const limit = 25
const contacts = ref<Contact[]>([])
const total = ref(0)
const loading = ref(false)
const selectedContact = ref<Contact | null>(null)
const slideoverOpen = ref(false)

const columns: TableColumn<Contact>[] = [{
  accessorKey: 'name',
  header: 'Name',
  enableSorting: true
}, {
  accessorKey: 'identifier',
  header: 'Identifier'
}, {
  accessorKey: 'blocked',
  header: 'Status'
}, {
  accessorKey: 'createdAt',
  header: 'Created',
  enableSorting: true
}]

const pageCount = computed(() => Math.ceil(total.value / limit))

async function loadContacts() {
  loading.value = true
  try {
    const res = await list({ page: page.value, limit, search: search.value || undefined })
    contacts.value = res.data ?? []
    total.value = res.total ?? 0
  } catch {
    toast.add({ title: 'Failed to load contacts', color: 'error' })
  } finally {
    loading.value = false
  }
}

function openContact(_e: Event, row: { original: Contact }) {
  selectedContact.value = row.original
  slideoverOpen.value = true
}

function formatDate(str: string): string {
  return new Date(str).toLocaleDateString()
}

let searchTimer: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    loadContacts()
  }, 300)
})
watch(page, loadContacts)
onMounted(loadContacts)
</script>

<template>
  <UDashboardPanel id="contacts">
    <template #header>
      <UDashboardNavbar title="Contacts">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
        <template #right>
          <UInput
            v-model="search"
            icon="i-lucide-search"
            placeholder="Search contacts…"
            size="sm"
            class="w-48"
          />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <UTable
        :data="contacts"
        :columns="columns"
        :loading="loading"
        class="w-full"
        @select="openContact"
      >
        <template #name-cell="{ row }">
          <div class="flex items-center gap-3">
            <UAvatar :alt="row.original.pushName || row.original.name" :src="row.original.avatarUrl" size="sm" />
            <div>
              <p class="font-medium text-highlighted">
                {{ row.original.pushName || row.original.name }}
              </p>
              <p v-if="row.original.pushName && row.original.name !== row.original.pushName" class="text-xs text-muted">
                {{ row.original.name }}
              </p>
            </div>
          </div>
        </template>

        <template #identifier-cell="{ row }">
          <span class="font-mono text-sm">{{ row.original.identifier }}</span>
        </template>

        <template #blocked-cell="{ row }">
          <UBadge
            :color="row.original.blocked ? 'error' : 'success'"
            variant="subtle"
            size="xs"
          >
            {{ row.original.blocked ? 'Blocked' : 'Active' }}
          </UBadge>
        </template>

        <template #createdAt-cell="{ row }">
          <span class="text-muted text-sm">{{ formatDate(row.original.createdAt) }}</span>
        </template>
      </UTable>

      <div v-if="pageCount > 1" class="flex justify-center py-4">
        <UPagination v-model:page="page" :total="total" :page-size="limit" />
      </div>
    </template>
  </UDashboardPanel>

  <ContactsContactSlideover
    v-if="selectedContact"
    v-model:open="slideoverOpen"
    :contact="selectedContact"
  />
</template>
