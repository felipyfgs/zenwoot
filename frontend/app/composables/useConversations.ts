import type { Conversation, Message, ConversationStatus, Label } from '~/types'

export function useConversations() {
  const api = useApi()

  async function listByInbox(inboxId: string, params?: { status?: ConversationStatus, page?: number, limit?: number }) {
    return api.paginated<Conversation>(`/inboxes/${inboxId}/conversations`, params as Record<string, unknown>)
  }

  async function get(id: string): Promise<Conversation> {
    return api.get<Conversation>(`/conversations/${id}`)
  }

  async function toggleStatus(id: string, status: ConversationStatus): Promise<Conversation> {
    return api.post<Conversation>(`/conversations/${id}/toggle_status`, { status })
  }

  async function markRead(id: string): Promise<void> {
    await api.post(`/conversations/${id}/read`)
  }

  async function listMessages(id: string): Promise<Message[]> {
    return api.get<Message[]>(`/conversations/${id}/messages`)
  }

  async function assign(id: string, assigneeId: string | null, teamId?: string | null): Promise<Conversation> {
    return api.post<Conversation>(`/conversations/${id}/assign`, { assigneeId, teamId })
  }

  async function updatePriority(id: string, priority: string): Promise<Conversation> {
    return api.post<Conversation>(`/conversations/${id}/priority`, { priority })
  }

  async function snooze(id: string, snoozedUntil: string): Promise<Conversation> {
    return api.post<Conversation>(`/conversations/${id}/snooze`, { snoozedUntil })
  }

  async function unsnooze(id: string): Promise<Conversation> {
    return api.post<Conversation>(`/conversations/${id}/unsnooze`)
  }

  async function getLabels(id: string): Promise<Label[]> {
    return api.get<Label[]>(`/conversations/${id}/labels`)
  }

  async function setLabels(id: string, labelIds: string[]): Promise<void> {
    await api.put(`/conversations/${id}/labels`, { labelIds })
  }

  async function addLabel(id: string, labelId: string): Promise<void> {
    await api.post(`/conversations/${id}/labels/${labelId}`)
  }

  async function removeLabel(id: string, labelId: string): Promise<void> {
    await api.del(`/conversations/${id}/labels/${labelId}`)
  }

  return {
    listByInbox,
    get,
    toggleStatus,
    markRead,
    listMessages,
    assign,
    updatePriority,
    snooze,
    unsnooze,
    getLabels,
    setLabels,
    addLabel,
    removeLabel
  }
}
