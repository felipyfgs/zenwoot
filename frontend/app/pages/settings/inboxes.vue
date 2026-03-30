<script setup lang="ts">
import type { Inbox } from '~/types'

const api = useApi()
const toast = useToast()

const inboxes = ref<Inbox[]>([])
const loading = ref(false)
const search = ref('')

const qrOpen = ref(false)
const qrData = ref('')
const qrInboxId = ref('')
const qrLoading = ref(false)
const deleteOpen = ref(false)
const deletingId = ref<string | null>(null)


const filtered = computed(() =>
  inboxes.value.filter(i => i.name.toLowerCase().includes(search.value.toLowerCase()))
)

async function load() {
  loading.value = true
  try {
    const data = await api.get<Inbox[]>('/inboxes')
    inboxes.value = Array.isArray(data) ? data : []
  } catch {
    toast.add({ title: 'Failed to load inboxes', color: 'error' })
  } finally {
    loading.value = false
  }
}

async function connect(id: string) {
  try {
    await api.post(`/inboxes/${id}/connect`)
    toast.add({ title: 'Connecting…', color: 'info' })
    await load()
  } catch { toast.add({ title: 'Failed to connect', color: 'error' }) }
}

async function disconnect(id: string) {
  try {
    await api.post(`/inboxes/${id}/disconnect`)
    toast.add({ title: 'Disconnected', color: 'neutral' })
    await load()
  } catch { toast.add({ title: 'Failed to disconnect', color: 'error' }) }
}

async function showQR(id: string) {
  qrInboxId.value = id
  qrData.value = ''
  qrOpen.value = true
  qrLoading.value = true
  try {
    const res = await api.get<{ qr: string }>(`/inboxes/${id}/qr`)
    qrData.value = res?.qr || ''
  } catch {
    toast.add({ title: 'Failed to fetch QR', color: 'error' })
  } finally {
    qrLoading.value = false
  }
}

function confirmDelete(id: string) {
  deletingId.value = id
  deleteOpen.value = true
}

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.del(`/inboxes/${deletingId.value}`)
    toast.add({ title: 'Inbox deleted', color: 'success' })
    deleteOpen.value = false
    await load()
  } catch { toast.add({ title: 'Failed to delete', color: 'error' }) }
}

function statusColor(status: string) {
  if (status === 'active') return 'success'
  if (status === 'connecting') return 'warning'
  if (status === 'disconnected') return 'error'
  return 'neutral'
}

onMounted(load)
</script>

<template>
  <UDashboardPanel id="settings-inboxes">
    <template #header>
      <UDashboardNavbar title="Inboxes">
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
          <NuxtLink to="/settings/inbox-new">
            <UButton icon="i-lucide-plus" size="sm">
              Nova Caixa
            </UButton>
          </NuxtLink>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div v-if="loading" class="flex justify-center py-12">
        <UIcon name="i-lucide-loader-circle" class="size-6 animate-spin text-muted" />
      </div>
      <p v-else-if="!filtered.length" class="text-center text-muted text-sm py-12">
        No inboxes found
      </p>
      <ul v-else class="divide-y divide-default">
        <li v-for="inbox in filtered" :key="inbox.id" class="flex items-center gap-4 px-4 py-3">
          <div class="size-8 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
            <UIcon name="i-lucide-smartphone" class="size-4 text-primary" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="font-medium text-highlighted truncate">
              {{ inbox.name }}
            </p>
            <p class="text-xs text-muted capitalize">
              {{ inbox.channelType }}
            </p>
          </div>
          <UBadge :color="statusColor(inbox.status)" variant="subtle" size="xs">
            {{ inbox.status }}
          </UBadge>
          <div class="flex items-center gap-1">
            <UButton
              v-if="inbox.status === 'disconnected' || inbox.status === 'inactive'"
              size="xs"
              variant="ghost"
              icon="i-lucide-plug"
              @click="connect(inbox.id)"
            >
              Connect
            </UButton>
            <UButton
              v-if="inbox.status === 'active'"
              size="xs"
              variant="ghost"
              icon="i-lucide-plug-zap"
              @click="disconnect(inbox.id)"
            >
              Disconnect
            </UButton>
            <UButton
              size="xs"
              variant="ghost"
              icon="i-lucide-qr-code"
              @click="showQR(inbox.id)"
            >
              QR
            </UButton>
            <UButton
              size="xs"
              variant="ghost"
              color="error"
              icon="i-lucide-trash-2"
              @click="confirmDelete(inbox.id)"
            />
          </div>
        </li>
      </ul>
    </template>
  </UDashboardPanel>

  <!-- QR modal -->
  <UModal v-model:open="qrOpen" title="Scan QR Code">
    <template #body>
      <div class="flex flex-col items-center gap-4 py-4">
        <UIcon v-if="qrLoading" name="i-lucide-loader-circle" class="size-12 animate-spin text-muted" />
        <img
          v-else-if="qrData"
          :src="qrData"
          class="size-64 object-contain"
          alt="QR Code"
        >
        <p v-else class="text-sm text-muted">
          No QR available — inbox may already be connected.
        </p>
      </div>
    </template>
    <template #footer>
      <UButton variant="ghost" color="neutral" @click="qrOpen = false">
        Close
      </UButton>
    </template>
  </UModal>

  <!-- Delete confirm modal -->
  <UModal v-model:open="deleteOpen" title="Delete Inbox">
    <template #body>
      <p class="text-sm text-default">
        Are you sure you want to delete this inbox? This action cannot be undone.
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
