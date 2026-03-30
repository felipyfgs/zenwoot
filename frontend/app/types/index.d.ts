// ─── Auth ─────────────────────────────────────────────────────────────────────

export interface AuthUser {
  id: string
  email: string
  name: string
  displayName?: string
  avatarUrl?: string
  role: UserRole
  status: UserStatus
  provider?: string
  settings?: unknown
  createdAt: string
  updatedAt: string
}

// ─── User / Agent ─────────────────────────────────────────────────────────────

export type UserRole = 'agent' | 'admin' | 'super_admin'
export type UserStatus = 'active' | 'inactive' | 'pending'

export interface User {
  id: string
  email: string
  name: string
  displayName?: string
  avatarUrl?: string
  role: UserRole
  status: UserStatus
  provider?: string
  settings?: unknown
  createdAt: string
  updatedAt: string
}

// ─── Inbox ────────────────────────────────────────────────────────────────────

export type ChannelType = 'whatsapp'
export type InboxStatus = 'active' | 'inactive' | 'connecting' | 'disconnected'

export interface Inbox {
  id: string
  accountId: string
  name: string
  channelType: ChannelType
  status: InboxStatus
  settings?: unknown
  createdAt: string
  updatedAt: string
}

// ─── Contact ─────────────────────────────────────────────────────────────────

export interface Contact {
  id: string
  accountId: string
  name: string
  pushName?: string
  identifier: string
  phone?: string
  email?: string
  avatarUrl?: string
  blocked: boolean
  additionalAttributes?: unknown
  createdAt: string
  updatedAt: string
}

// ─── Conversation ─────────────────────────────────────────────────────────────

export type ConversationStatus = 'open' | 'resolved' | 'pending' | 'snoozed'
export type ConversationPriority = 'none' | 'low' | 'medium' | 'high' | 'urgent'

export interface Conversation {
  id: string
  accountId: string
  inboxId: string
  contactId: string
  status: ConversationStatus
  priority: ConversationPriority
  assigneeId?: string
  teamId?: string
  unreadCount: number
  snoozedUntil?: string
  additionalAttributes?: unknown
  createdAt: string
  updatedAt: string
  contact?: Contact
  inbox?: Inbox
  lastMessage?: Message
  labels?: Label[]
  assignee?: User
}

// ─── Message ─────────────────────────────────────────────────────────────────

export type MessageDirection = 'inbound' | 'outbound'
export type ContentType = 'text' | 'image' | 'video' | 'audio' | 'document' | 'sticker' | 'location'
export type MessageStatus = 'sent' | 'delivered' | 'read' | 'failed'

export interface Message {
  id: string
  conversationId: string
  inboxId: string
  direction: MessageDirection
  contentType: ContentType
  content: string
  mediaUrl?: string
  mediaCaption?: string
  status: MessageStatus
  authorId?: string
  externalId?: string
  createdAt: string
  updatedAt: string
}

// ─── Team ─────────────────────────────────────────────────────────────────────

export interface Team {
  id: string
  accountId: string
  name: string
  description?: string
  allowAutoAssign: boolean
  createdAt: string
  updatedAt: string
}

// ─── Label ───────────────────────────────────────────────────────────────────

export interface Label {
  id: string
  accountId: string
  title: string
  description?: string
  color: string
  showOnSidebar: boolean
  createdAt: string
  updatedAt: string
}

// ─── Canned Response ─────────────────────────────────────────────────────────

export interface CannedResponse {
  id: string
  accountId: string
  shortCode: string
  content: string
  createdAt: string
  updatedAt: string
}

// ─── Note ────────────────────────────────────────────────────────────────────

export interface Note {
  id: string
  contactId: string
  userId: string
  content: string
  createdAt: string
  updatedAt: string
}

// ─── API Response envelope ───────────────────────────────────────────────────

export interface ApiResponse<T = unknown> {
  success: boolean
  data: T
  error: string
  message: string
}

export interface PaginatedResponse<T = unknown> {
  success: boolean
  data: T[]
  total: number
  page: number
  limit: number
  error: string
  message: string
}
