import type { Contact, Conversation, Note } from '~/types'

export function useContacts() {
  const api = useApi()

  async function list(params?: { page?: number, limit?: number, search?: string }) {
    return api.paginated<Contact>('/contacts', params as Record<string, unknown>)
  }

  async function get(id: string): Promise<Contact> {
    return api.get<Contact>(`/contacts/${id}`)
  }

  async function listConversations(id: string): Promise<Conversation[]> {
    return api.get<Conversation[]>(`/contacts/${id}/conversations`)
  }

  async function listNotes(id: string): Promise<Note[]> {
    return api.get<Note[]>(`/contacts/${id}/notes`)
  }

  async function createNote(id: string, content: string): Promise<Note> {
    return api.post<Note>(`/contacts/${id}/notes`, { content })
  }

  async function deleteNote(contactId: string, noteId: string): Promise<void> {
    await api.del(`/contacts/${contactId}/notes/${noteId}`)
  }

  return { list, get, listConversations, listNotes, createNote, deleteNote }
}
