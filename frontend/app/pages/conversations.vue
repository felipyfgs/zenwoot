<script setup lang="ts">
import { breakpointsTailwind } from '@vueuse/core'
import type { Conversation, Inbox, ConversationStatus } from '~/types'

const { get: getApi, paginated } = useApi()
const toast = useToast()

const STATUS_TABS: { label: string, value: ConversationStatus | 'all' }[] = [
  { label: 'Open', value: 'open' },
  { label: 'Resolved', value: 'resolved' },
  { label: 'Pending', value: 'pending' },
  { label: 'Snoozed', value: 'snoozed' }
]

const inboxes = ref<Inbox[]>([])
const selectedInboxId = ref<string | undefined>(undefined)
const activeStatus = ref<ConversationStatus | 'all'>('open')
const conversations = ref<Conversation[]>([])
const selectedConversation = ref<Conversation | null>(null)
const loading = ref(false)

const isConvPanelOpen = computed({
  get() { return !!selectedConversation.value },
  set(v: boolean) { if (!v) selectedConversation.value = null }, // eslint-disable-line
})

const inboxItems = computed(() =>
  inboxes.value.map(i => ({ label: i.name, value: i.id }))
)

const breakpoints = useBreakpoints(breakpointsTailwind)
const isMobile = breakpoints.smaller('lg')

async function loadInboxes() {
  try {
    const data = await getApi<Inbox[]>('/inboxes')
    inboxes.value = Array.isArray(data) ? data : []
    if (inboxes.value.length && !selectedInboxId.value) {
      selectedInboxId.value = inboxes.value[0]?.id
    }
  } catch {
    toast.add({ title: 'Failed to load inboxes', color: 'error' })
  }
}

async function loadConversations() {
  if (!selectedInboxId.value) return
  loading.value = true
  try {
    const query: Record<string, unknown> = { page: 1, limit: 50 }
    if (activeStatus.value !== 'all') query.status = activeStatus.value
    const res = await paginated<Conversation>(`/inboxes/${selectedInboxId.value}/conversations`, query)
    conversations.value = res.data ?? []
  } catch {
    toast.add({ title: 'Failed to load conversations', color: 'error' })
  } finally {
    loading.value = false
  }
}

watch([selectedInboxId, activeStatus], loadConversations)
onMounted(async () => {
  await loadInboxes()
  await loadConversations()
})
</script>

<template>
  <!-- List panel -->
  <UDashboardPanel
    id="conversations-list"
    :default-size="22"
    :min-size="16"
    :max-size="35"
    resizable
  >
    <UDashboardNavbar title="Conversations">
      <template #leading>
        <UDashboardSidebarCollapse />
      </template>
      <template #trailing>
        <UBadge :label="String(conversations.length)" variant="subtle" />
      </template>
      <template #right>
        <USelect
          v-if="inboxItems.length"
          v-model="selectedInboxId"
          :items="inboxItems"
          value-key="value"
          size="xs"
          class="w-36"
        />
      </template>
    </UDashboardNavbar>

    <UDashboardToolbar>
      <template #left>
        <div class="flex items-center gap-1">
          <UButton
            v-for="tab in STATUS_TABS"
            :key="tab.value"
            :variant="activeStatus === tab.value ? 'soft' : 'ghost'"
            :color="activeStatus === tab.value ? 'primary' : 'neutral'"
            size="xs"
            @click="activeStatus = tab.value"
          >
            {{ tab.label }}
          </UButton>
        </div>
      </template>
    </UDashboardToolbar>

    <ConversationsConversationList
      :conversations="conversations"
      :loading="loading"
      :selected-id="selectedConversation?.id"
      @select="(conv) => { selectedConversation = conv }"
    />
  </UDashboardPanel>

  <!-- Detail panel (desktop) -->
  <ConversationsConversationDetail
    v-if="selectedConversation && !isMobile"
    :conversation="selectedConversation"
    @updated="loadConversations"
    @close="selectedConversation = null"
  />
  <div v-else-if="!isMobile" class="hidden lg:flex flex-1 flex-col items-center justify-center gap-3 text-dimmed">
    <UIcon name="i-lucide-message-circle" class="size-16" />
    <p class="text-sm">
      Select a conversation
    </p>
  </div>

  <!-- Mobile slideover -->
  <ClientOnly>
    <USlideover v-if="isMobile" v-model:open="isConvPanelOpen">
      <template #content>
        <ConversationsConversationDetail
          v-if="selectedConversation"
          :conversation="selectedConversation"
          @updated="loadConversations"
          @close="selectedConversation = null"
        />
      </template>
    </USlideover>
  </ClientOnly>
</template>
