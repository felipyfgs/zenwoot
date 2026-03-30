<script setup lang="ts">
import type { Conversation, Message, ConversationStatus } from '~/types'

const props = defineProps<{
  conversation: Conversation
}>()

const emit = defineEmits<{
  updated: []
  close: []
}>()

const { toggleStatus, markRead, listMessages } = useConversations()
const toast = useToast()

const messages = ref<Message[]>([])
const loadingMessages = ref(false)
const replyText = ref('')
const sendingReply = ref(false)

const statusLabel = computed(() =>
  props.conversation.status === 'open' ? 'Resolve' : 'Reopen'
)
const statusIcon = computed(() =>
  props.conversation.status === 'open' ? 'i-lucide-check-circle' : 'i-lucide-refresh-ccw'
)

function contactName(conv: Conversation): string {
  return conv.contact?.pushName || conv.contact?.name || conv.contact?.identifier || 'Unknown'
}

async function fetchMessages() {
  loadingMessages.value = true
  try {
    const res = await listMessages(props.conversation.id)
    messages.value = Array.isArray(res) ? res : []
    await markRead(props.conversation.id)
  } catch {
    toast.add({ title: 'Failed to load messages', color: 'error' })
  } finally {
    loadingMessages.value = false
  }
}

async function handleToggleStatus() {
  const next: ConversationStatus = props.conversation.status === 'open' ? 'resolved' : 'open'
  try {
    await toggleStatus(props.conversation.id, next)
    toast.add({ title: `Conversation ${next}`, color: 'success' })
    emit('updated')
  } catch {
    toast.add({ title: 'Failed to update status', color: 'error' })
  }
}

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

watch(() => props.conversation.id, fetchMessages, { immediate: true })
</script>

<template>
  <UDashboardPanel id="conversation-detail">
    <UDashboardNavbar :toggle="false">
      <template #leading>
        <UButton
          icon="i-lucide-x"
          color="neutral"
          variant="ghost"
          class="-ms-1.5 lg:hidden"
          @click="emit('close')"
        />
        <div class="flex items-center gap-2">
          <UAvatar
            :alt="contactName(conversation)"
            :src="conversation.contact?.avatarUrl"
            size="sm"
          />
          <div>
            <p class="text-sm font-semibold text-highlighted leading-none">
              {{ contactName(conversation) }}
            </p>
            <p class="text-xs text-muted mt-0.5">
              {{ conversation.inbox?.name || 'Inbox' }}
            </p>
          </div>
        </div>
      </template>

      <template #right>
        <UBadge
          v-if="conversation.priority && conversation.priority !== 'none' && conversation.priority !== 'low'"
          variant="subtle"
          size="sm"
        >
          {{ conversation.priority }}
        </UBadge>
        <UButton
          :icon="statusIcon"
          :label="statusLabel"
          :color="conversation.status === 'open' ? 'success' : 'neutral'"
          variant="soft"
          size="sm"
          @click="handleToggleStatus"
        />
      </template>
    </UDashboardNavbar>

    <!-- Messages body -->
    <div class="flex-1 flex flex-col gap-3 p-4 sm:p-6 overflow-y-auto">
      <div v-if="loadingMessages" class="flex justify-center py-8">
        <UIcon name="i-lucide-loader-circle" class="size-6 animate-spin text-muted" />
      </div>

      <template v-else>
        <p v-if="!messages.length" class="text-center text-muted text-sm py-8">
          No messages yet
        </p>
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex"
          :class="msg.direction === 'outbound' ? 'justify-end' : 'justify-start'"
        >
          <div
            class="max-w-[70%] rounded-2xl px-4 py-2 text-sm"
            :class="msg.direction === 'outbound'
              ? 'bg-primary text-white rounded-br-sm'
              : 'bg-elevated text-default rounded-bl-sm'"
          >
            <img
              v-if="msg.contentType === 'image' && msg.mediaUrl"
              :src="msg.mediaUrl"
              class="rounded-lg max-w-xs mb-1"
              alt="image"
            >
            <p v-if="msg.content" class="whitespace-pre-wrap break-words">
              {{ msg.content }}
            </p>
            <p v-else-if="msg.contentType === 'document'" class="flex items-center gap-1">
              <UIcon name="i-lucide-file" class="size-4" /> Document
            </p>
            <p v-else-if="msg.contentType === 'audio'" class="flex items-center gap-1">
              <UIcon name="i-lucide-mic" class="size-4" /> Audio
            </p>
            <p class="text-xs opacity-60 mt-1 text-right">
              {{ formatTime(msg.createdAt) }}
            </p>
          </div>
        </div>
      </template>
    </div>

    <!-- Reply footer -->
    <div class="px-4 pb-4 sm:px-6 shrink-0">
      <UCard
        variant="subtle"
        :ui="{ header: 'flex items-center gap-1.5 text-dimmed' }"
      >
        <template #header>
          <UIcon name="i-lucide-reply" class="size-4" />
          <span class="text-sm truncate">Reply to {{ contactName(conversation) }}</span>
        </template>

        <UTextarea
          v-model="replyText"
          color="neutral"
          variant="none"
          autoresize
          placeholder="Write your reply…"
          :rows="3"
          :disabled="sendingReply"
          class="w-full"
          :ui="{ base: 'p-0 resize-none' }"
        />

        <div class="flex items-center justify-end gap-2 mt-2">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-paperclip"
          />
          <UButton
            icon="i-lucide-send"
            label="Send"
            :loading="sendingReply"
            :disabled="!replyText.trim()"
            @click="() => {}"
          />
        </div>
      </UCard>
    </div>
  </UDashboardPanel>
</template>
