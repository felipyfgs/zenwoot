<script setup lang="ts">
import type { Contact, Conversation, Note } from '~/types'

const props = defineProps<{
  contact: Contact
  open: boolean
}>()

const emit = defineEmits<{
  'update:open': [val: boolean]
}>()

const { listConversations, listNotes, createNote, deleteNote } = useContacts()
const toast = useToast()

const tab = ref('conversations')
const conversations = ref<Conversation[]>([])
const notes = ref<Note[]>([])
const newNote = ref('')
const loadingConvs = ref(false)
const loadingNotes = ref(false)
const savingNote = ref(false)

const tabs = [
  { label: 'Conversations', value: 'conversations' },
  { label: 'Notes', value: 'notes' }
]

const displayName = computed(() =>
  props.contact.pushName || props.contact.name || props.contact.identifier
)

async function loadConversations() {
  loadingConvs.value = true
  try {
    const res = await listConversations(props.contact.id)
    conversations.value = Array.isArray(res) ? res : []
  } catch {
    toast.add({ title: 'Failed to load conversations', color: 'error' })
  } finally {
    loadingConvs.value = false
  }
}

async function loadNotes() {
  loadingNotes.value = true
  try {
    const res = await listNotes(props.contact.id)
    notes.value = Array.isArray(res) ? res : []
  } catch {
    toast.add({ title: 'Failed to load notes', color: 'error' })
  } finally {
    loadingNotes.value = false
  }
}

async function addNote() {
  if (!newNote.value.trim()) return
  savingNote.value = true
  try {
    const n = await createNote(props.contact.id, newNote.value.trim())
    notes.value.unshift(n)
    newNote.value = ''
  } catch {
    toast.add({ title: 'Failed to add note', color: 'error' })
  } finally {
    savingNote.value = false
  }
}

async function removeNote(noteId: string) {
  try {
    await deleteNote(props.contact.id, noteId)
    notes.value = notes.value.filter(n => n.id !== noteId)
  } catch {
    toast.add({ title: 'Failed to delete note', color: 'error' })
  }
}

function convStatusColor(status: string) {
  if (status === 'open') return 'success'
  if (status === 'resolved') return 'neutral'
  if (status === 'pending') return 'warning'
  return 'neutral'
}

function formatDate(str: string) {
  return new Date(str).toLocaleDateString()
}

watch(() => props.open, (val) => {
  if (val) {
    loadConversations()
    loadNotes()
  }
}, { immediate: true })

watch(tab, (val) => {
  if (val === 'conversations' && conversations.value.length === 0) loadConversations()
  if (val === 'notes' && notes.value.length === 0) loadNotes()
})
</script>

<template>
  <USlideover
    :open="open"
    side="right"
    class="max-w-md"
    @update:open="emit('update:open', $event)"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <UAvatar :alt="displayName" :src="contact.avatarUrl" size="md" />
        <div>
          <p class="font-semibold text-highlighted">
            {{ displayName }}
          </p>
          <p class="text-xs text-muted font-mono">
            {{ contact.identifier }}
          </p>
        </div>
        <UBadge
          :color="contact.blocked ? 'error' : 'success'"
          variant="subtle"
          size="xs"
          class="ml-auto"
        >
          {{ contact.blocked ? 'Blocked' : 'Active' }}
        </UBadge>
      </div>
    </template>

    <template #body>
      <div class="flex gap-2 mb-4">
        <UButton
          v-for="t in tabs"
          :key="t.value"
          :variant="tab === t.value ? 'soft' : 'ghost'"
          :color="tab === t.value ? 'primary' : 'neutral'"
          size="sm"
          @click="tab = t.value"
        >
          {{ t.label }}
        </UButton>
      </div>

      <!-- Conversations tab -->
      <div v-if="tab === 'conversations'">
        <div v-if="loadingConvs" class="flex justify-center py-8">
          <UIcon name="i-lucide-loader-circle" class="size-5 animate-spin text-muted" />
        </div>
        <p v-else-if="!conversations.length" class="text-sm text-muted text-center py-8">
          No conversations
        </p>
        <div v-else class="space-y-2">
          <div
            v-for="conv in conversations"
            :key="conv.id"
            class="rounded-lg border border-default p-3 space-y-1"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-highlighted">
                {{ conv.inbox?.name || 'Inbox' }}
              </span>
              <UBadge :color="convStatusColor(conv.status)" variant="subtle" size="xs">
                {{ conv.status }}
              </UBadge>
            </div>
            <p v-if="conv.lastMessage" class="text-xs text-muted truncate">
              {{ conv.lastMessage.content }}
            </p>
            <p class="text-xs text-muted">
              {{ formatDate(conv.createdAt) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Notes tab -->
      <div v-if="tab === 'notes'" class="space-y-4">
        <div class="flex gap-2">
          <UTextarea
            v-model="newNote"
            placeholder="Add a note…"
            :rows="2"
            autoresize
            class="flex-1"
            :disabled="savingNote"
          />
          <UButton
            icon="i-lucide-plus"
            :loading="savingNote"
            :disabled="!newNote.trim()"
            color="primary"
            class="self-end"
            @click="addNote"
          />
        </div>

        <div v-if="loadingNotes" class="flex justify-center py-4">
          <UIcon name="i-lucide-loader-circle" class="size-5 animate-spin text-muted" />
        </div>
        <p v-else-if="!notes.length" class="text-sm text-muted text-center py-4">
          No notes yet
        </p>
        <div v-else class="space-y-2">
          <div
            v-for="note in notes"
            :key="note.id"
            class="rounded-lg border border-default p-3"
          >
            <div class="flex items-start justify-between gap-2">
              <p class="text-sm text-default whitespace-pre-wrap flex-1">
                {{ note.content }}
              </p>
              <UButton
                icon="i-lucide-trash-2"
                size="xs"
                color="error"
                variant="ghost"
                @click="removeNote(note.id)"
              />
            </div>
            <p class="text-xs text-muted mt-1">
              {{ formatDate(note.createdAt) }}
            </p>
          </div>
        </div>
      </div>
    </template>
  </USlideover>
</template>
