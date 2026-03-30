<script setup lang="ts">
import type { Conversation } from '~/types'

const _props = defineProps<{
  conversations: Conversation[]
  loading?: boolean
  selectedId?: string
}>()

const emit = defineEmits<{
  select: [conv: Conversation]
}>()

function contactName(conv: Conversation): string {
  return conv.contact?.pushName || conv.contact?.name || conv.contact?.identifier || 'Unknown'
}

function lastMessagePreview(conv: Conversation): string {
  if (!conv.lastMessage) return 'No messages yet'
  const c = conv.lastMessage.content
  if (conv.lastMessage.contentType === 'image') return '📷 Image'
  if (conv.lastMessage.contentType === 'document') return '📄 Document'
  if (conv.lastMessage.contentType === 'audio') return '🎵 Audio'
  return c?.slice(0, 80) || ''
}

function timeAgo(dateStr: string): string {
  const d = new Date(dateStr)
  const diff = Date.now() - d.getTime()
  const m = Math.floor(diff / 60000)
  if (m < 1) return 'now'
  if (m < 60) return `${m}m`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h`
  const days = Math.floor(h / 24)
  return `${days}d`
}

const priorityColor = {
  urgent: 'error',
  high: 'warning',
  medium: 'info',
  low: 'neutral',
  none: undefined
} as const satisfies Record<string, 'error' | 'warning' | 'info' | 'neutral' | undefined>
</script>

<template>
  <div class="flex-1 overflow-y-auto divide-y divide-default">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <UIcon name="i-lucide-loader-circle" class="size-6 animate-spin text-muted" />
    </div>

    <p v-else-if="!conversations.length" class="text-center text-muted text-sm py-12">
      No conversations
    </p>

    <button
      v-for="conv in conversations"
      :key="conv.id"
      class="w-full text-left px-4 py-3 flex gap-3 items-start hover:bg-elevated/50 transition-colors"
      :class="{ 'bg-elevated': conv.id === selectedId }"
      @click="emit('select', conv)"
    >
      <UAvatar
        :alt="contactName(conv)"
        :src="conv.contact?.avatarUrl"
        size="sm"
        class="shrink-0 mt-0.5"
      />

      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-highlighted truncate">
            {{ contactName(conv) }}
          </span>
          <span class="text-xs text-muted shrink-0">
            {{ conv.lastMessage ? timeAgo(conv.lastMessage.createdAt) : timeAgo(conv.createdAt) }}
          </span>
        </div>

        <div class="flex items-center justify-between gap-2 mt-0.5">
          <p class="text-xs text-muted truncate">
            {{ lastMessagePreview(conv) }}
          </p>
          <div class="flex items-center gap-1 shrink-0">
            <UBadge
              v-if="conv.priority && conv.priority !== 'none'"
              :color="priorityColor[conv.priority] || 'neutral'"
              variant="subtle"
              size="xs"
            >
              {{ conv.priority }}
            </UBadge>
            <UBadge
              v-if="conv.unreadCount > 0"
              color="primary"
              size="xs"
            >
              {{ conv.unreadCount }}
            </UBadge>
          </div>
        </div>
      </div>
    </button>
  </div>
</template>
