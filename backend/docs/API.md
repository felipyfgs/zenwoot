# wzap API Reference
Multi-session WhatsApp REST API built on [wzap API](https://github.com/wzap/wzap).

---

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Media Input (URL / Base64 / Form-Data)](#media-input)
- [Health](#health)
- [Sessions](#sessions)
- [Messages](#messages)
- [Contacts](#contacts)
- [Groups](#groups)
- [Chat](#chat)
- [Labels](#labels)
- [Newsletter](#newsletter)
- [Community](#community)
- [Webhooks](#webhooks)
- [Webhook Payload](#webhook-payload)
- [Supported Event Types](#supported-event-types)

---

## Base URL

```
http://<host>:<port>
```

Default port: `8080` (configurable via `PORT` env var).

---

## Authentication

All endpoints (except `/health` and `/swagger/*`) require an `ApiKey` header:

```
ApiKey: <key>
```

| Type | Value | Access |
|---|---|---|
| **Admin** | `API_KEY` env var | Full access -- session CRUD + all session-scoped routes |
| **Session** | Session `apiKey` returned at creation | Scoped to that session only |

---

## Response Format

### Success

```json
{
  "success": true,
  "data": { ... },
  "message": "success"
}
```

### Error

```json
{
  "success": false,
  "error": "Error type",
  "message": "Details about the error"
}
```

---

## Media Input

All media endpoints (image, video, document, audio, sticker) accept **three input formats**. The system auto-detects the format from the field value:

| Format | Detection | Example |
|---|---|---|
| **URL publica** | Starts with `http://` or `https://` | `"https://example.com/photo.jpg"` |
| **Data URI** | Starts with `data:` | `"data:image/png;base64,iVBOR..."` |
| **Base64 puro** | Anything else | `"iVBORw0KGgoAAAANSUhEUg..."` |

Each media route accepts a **named field** matching the media type (`image`, `video`, `document`, `audio`, `sticker`), or the generic `url` field, or the legacy `base64` field. Priority: named field > `url` > `base64`.

**Multipart form-data** is also supported: send the file in a `file` field plus other fields (`phone`, `type`, `mimeType`, `caption`, `filename`) as form values.

### Examples

Via URL:
```json
{ "phone": "5511999999999", "image": "https://example.com/photo.jpg", "caption": "Look!" }
```

Via base64:
```json
{ "phone": "5511999999999", "image": "iVBORw0KGgoAAAANSUhEUg...", "mimeType": "image/png" }
```

Via unified route:
```json
{ "phone": "5511999999999", "type": "image", "url": "https://example.com/photo.jpg" }
```

Via form-data:
```bash
curl -X POST .../messages/image \
  -H 'ApiKey: KEY' \
  -F 'phone=5511999999999' \
  -F 'caption=Look!' \
  -F 'file=@photo.jpg'
```

---

## Health

### Health Check

```
GET /health
```

No authentication required.

```bash
curl http://localhost:8080/health
```

```json
{
  "success": true,
  "data": {
    "status": "UP",
    "services": { "database": true, "nats": true, "minio": true }
  }
}
```

---

## Sessions

All session-scoped routes use `/sessions/:sessionId/...` where `:sessionId` is the session name or UUID.

### Create Session (Admin Only)

```
POST /sessions
```

```json
{
  "name": "my-session",
  "apiKey": "optional-custom-key",
  "proxy": {
    "host": "proxy.example.com",
    "port": 3128,
    "protocol": "http",
    "username": "",
    "password": ""
  },
  "webhook": {
    "url": "https://my-server.com/hook",
    "events": ["Message", "Connected"]
  },
  "settings": {
    "alwaysOnline": false,
    "rejectCall": false,
    "msgRejectCall": "",
    "readMessages": false,
    "ignoreGroups": false,
    "ignoreStatus": false
  }
}
```

All fields except `name` are optional. `apiKey` is auto-generated if omitted.

**Response** (201):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "my-session",
    "apiKey": "generated-or-custom-key",
    "status": "disconnected",
    "connected": 0,
    "proxy": {},
    "settings": { ... },
    "webhook": { "id": "uuid", "url": "...", "events": [...] },
    "createdAt": "...",
    "updatedAt": "..."
  }
}
```

### List Sessions (Admin Only)

```
GET /sessions
```

### Get Session

```
GET /sessions/:sessionId
```

### Delete Session

```
DELETE /sessions/:sessionId
```

### Connect Session

```
POST /sessions/:sessionId/connect
```

Connects the WhatsApp session. If not yet paired, starts QR generation (poll `/qr`).

**Response:**
```json
{ "success": true, "data": { "status": "CONNECTED" } }
```

Status values: `PAIRING` | `CONNECTING` | `CONNECTED`

### Disconnect Session

```
POST /sessions/:sessionId/disconnect
```

### Get QR Code

```
GET /sessions/:sessionId/qr
```

Call `/connect` first, then poll this endpoint. Returns 404 if no QR is available (already paired or not yet connecting).

```json
{
  "success": true,
  "data": {
    "qr": "2@raw-qr-string...",
    "image": "data:image/png;base64,iVBOR..."
  }
}
```

---

## Messages

All message endpoints use the `phone` field, which accepts a phone number (`559981769536`) or full JID (`559981769536@s.whatsapp.net`, `120362023605733675@g.us`).

### Send Text

```
POST /sessions/:sessionId/messages/text
```

```json
{
  "phone": "559981769536",
  "text": "Hello!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "messageId": "3EB0FFFB732591EF45FD8B",
    "timestamp": "2026-03-29T20:24:33+02:00",
    "to": "559981769536@s.whatsapp.net",
    "from": "559981769536:54@s.whatsapp.net"
  }
}
```

### Send Image

```
POST /sessions/:sessionId/messages/image
```

```json
{
  "phone": "559981769536",
  "image": "https://example.com/photo.jpg",
  "caption": "Look at this",
  "mimeType": "image/jpeg"
}
```

`mimeType` is optional (auto-detected from content).

### Send Video

```
POST /sessions/:sessionId/messages/video
```

```json
{
  "phone": "559981769536",
  "video": "https://example.com/video.mp4",
  "caption": "Watch this"
}
```

### Send Document

```
POST /sessions/:sessionId/messages/document
```

```json
{
  "phone": "559981769536",
  "document": "https://example.com/file.pdf",
  "filename": "report.pdf",
  "caption": "Here is the report"
}
```

### Send Audio

```
POST /sessions/:sessionId/messages/audio
```

```json
{
  "phone": "559981769536",
  "audio": "https://example.com/audio.mp3"
}
```

Sent as PTT (push-to-talk voice message).

### Send Sticker

```
POST /sessions/:sessionId/messages/sticker
```

```json
{
  "phone": "559981769536",
  "sticker": "https://example.com/sticker.webp",
  "mimeType": "image/webp"
}
```

### Send Media (Unified)

```
POST /sessions/:sessionId/messages/media
```

Single route for any media type. The `type` field is required.

```json
{
  "phone": "559981769536",
  "type": "image",
  "url": "https://example.com/photo.jpg",
  "caption": "Sent via unified route"
}
```

`type` values: `image` | `video` | `document` | `audio` | `sticker`

### Send Contact

```
POST /sessions/:sessionId/messages/contact
```

```json
{
  "phone": "559981769536",
  "name": "John Doe",
  "vcard": "BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"
}
```

### Send Location

```
POST /sessions/:sessionId/messages/location
```

```json
{
  "phone": "559981769536",
  "lat": -23.5505,
  "lng": -46.6333,
  "name": "Sao Paulo",
  "address": "Av. Paulista, 1000"
}
```

### Send Poll

```
POST /sessions/:sessionId/messages/poll
```

```json
{
  "phone": "559981769536",
  "name": "Favorite color?",
  "options": ["Red", "Green", "Blue"],
  "selectableCount": 1
}
```

`selectableCount`: how many options can be selected. `0` = unlimited.

### Send Link Preview

```
POST /sessions/:sessionId/messages/link
```

```json
{
  "phone": "559981769536",
  "url": "https://github.com",
  "title": "GitHub",
  "description": "Dev platform"
}
```

### Edit Message

```
POST /sessions/:sessionId/messages/edit
```

```json
{
  "phone": "559981769536",
  "messageId": "3EB0FFFB732591EF45FD8B",
  "text": "Updated text"
}
```

### Delete Message

```
POST /sessions/:sessionId/messages/delete
```

Revokes the message for all recipients.

```json
{
  "phone": "559981769536",
  "messageId": "3EB0FFFB732591EF45FD8B"
}
```

### React to Message

```
POST /sessions/:sessionId/messages/reaction
```

Send empty `reaction` to remove an existing reaction.

```json
{
  "phone": "559981769536",
  "messageId": "3EB0FFFB732591EF45FD8B",
  "reaction": "👍"
}
```

### Mark Message as Read

```
POST /sessions/:sessionId/messages/read
```

```json
{
  "phone": "559981769536",
  "messageId": "3EB0FFFB732591EF45FD8B"
}
```

### Set Typing / Recording Presence

```
POST /sessions/:sessionId/messages/presence
```

```json
{
  "phone": "559981769536",
  "presence": "typing"
}
```

`presence` values: `typing` | `recording` | `paused`

---

## Contacts

### List Contacts

```
GET /sessions/:sessionId/contacts
```

### Check Contacts on WhatsApp

```
POST /sessions/:sessionId/contacts/check
```

```json
{ "phones": ["559981769536", "5511888888888"] }
```

**Response:**
```json
{
  "data": [
    { "exists": true, "jid": "559981769536@s.whatsapp.net", "phoneNumber": "559981769536" },
    { "exists": false, "phoneNumber": "5511888888888" }
  ]
}
```

### Get Contact Avatar

```
POST /sessions/:sessionId/contacts/avatar
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

**Response:**
```json
{ "data": { "url": "https://pps.whatsapp.net/v/...", "id": "394563270" } }
```

### Block Contact

```
POST /sessions/:sessionId/contacts/block
```

```json
{ "jid": "5511999999999@s.whatsapp.net" }
```

### Unblock Contact

```
POST /sessions/:sessionId/contacts/unblock
```

```json
{ "jid": "5511999999999@s.whatsapp.net" }
```

### Get Blocklist

```
GET /sessions/:sessionId/contacts/blocklist
```

### Get User Info

```
POST /sessions/:sessionId/contacts/info
```

```json
{ "jids": ["559981769536@s.whatsapp.net"] }
```

**Response:**
```json
{
  "data": {
    "559981769536@s.whatsapp.net": {
      "jid": "559981769536@s.whatsapp.net",
      "status": "My status text",
      "picture": "394563270",
      "devices": ["0", "51", "52"]
    }
  }
}
```

### Get Privacy Settings

```
GET /sessions/:sessionId/contacts/privacy
```

**Response:**
```json
{
  "data": {
    "GroupAdd": "all",
    "LastSeen": "all",
    "Status": "all",
    "Profile": "all",
    "ReadReceipts": "all",
    "CallAdd": "all",
    "Online": "all"
  }
}
```

### Set Profile Picture

```
POST /sessions/:sessionId/contacts/profile-picture
```

```json
{ "base64": "<base64-encoded-jpeg>" }
```

---

## Groups

### List Groups

```
GET /sessions/:sessionId/groups
```

### Create Group

```
POST /sessions/:sessionId/groups/create
```

```json
{
  "name": "My Group",
  "participants": ["5511999999999", "5511888888888"]
}
```

### Get Group Info

```
POST /sessions/:sessionId/groups/info
```

```json
{ "groupJid": "120363418922713210@g.us" }
```

### Get Group Info from Invite Code

```
POST /sessions/:sessionId/groups/invite-info
```

```json
{ "inviteCode": "EB1xRXbYyjlHqKyG7D56Q5" }
```

### Join Group with Invite Code

```
POST /sessions/:sessionId/groups/join
```

```json
{ "inviteCode": "EB1xRXbYyjlHqKyG7D56Q5" }
```

**Response:**
```json
{ "data": { "jid": "120363424675579035@g.us" } }
```

### Get Invite Link

```
POST /sessions/:sessionId/groups/invite-link
```

Optional query param `?reset=true` to generate a new link.

```json
{ "groupJid": "120363424675579035@g.us" }
```

**Response:**
```json
{ "data": { "link": "https://chat.whatsapp.com/EB1xRXbYyjlHqKyG7D56Q5" } }
```

### Leave Group

```
POST /sessions/:sessionId/groups/leave
```

```json
{ "groupJid": "120363424675579035@g.us" }
```

### Update Participants

```
POST /sessions/:sessionId/groups/participants
```

```json
{
  "groupJid": "120363424675579035@g.us",
  "participants": ["5511999999999"],
  "action": "add"
}
```

`action` values: `add` | `remove` | `promote` | `demote`

### Get Join Requests

```
POST /sessions/:sessionId/groups/requests
```

```json
{ "groupJid": "120363424675579035@g.us" }
```

### Approve / Reject Join Requests

```
POST /sessions/:sessionId/groups/requests/action
```

```json
{
  "groupJid": "120363424675579035@g.us",
  "participants": ["5511999999999"],
  "action": "approve"
}
```

`action` values: `approve` | `reject`

### Update Group Name

```
POST /sessions/:sessionId/groups/name
```

```json
{
  "groupJid": "120363424675579035@g.us",
  "text": "New Group Name"
}
```

### Update Group Description

```
POST /sessions/:sessionId/groups/description
```

```json
{
  "groupJid": "120363424675579035@g.us",
  "text": "New description"
}
```

### Update Group Photo

```
POST /sessions/:sessionId/groups/photo
```

```json
{
  "groupJid": "120363424675579035@g.us",
  "photoBase64": "<base64-encoded-jpeg>"
}
```

### Set Announce Mode

```
POST /sessions/:sessionId/groups/announce
```

Only admins can send messages when enabled.

```json
{ "groupJid": "120363424675579035@g.us", "enabled": true }
```

### Set Locked Mode

```
POST /sessions/:sessionId/groups/locked
```

Only admins can edit group info when enabled.

```json
{ "groupJid": "120363424675579035@g.us", "enabled": true }
```

### Set Join Approval

```
POST /sessions/:sessionId/groups/join-approval
```

Requires admin approval for new members.

```json
{ "groupJid": "120363424675579035@g.us", "enabled": true }
```

---

## Chat

### Archive Chat

```
POST /sessions/:sessionId/chat/archive
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

### Unarchive Chat

```
POST /sessions/:sessionId/chat/unarchive
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

### Mute Chat

```
POST /sessions/:sessionId/chat/mute
```

```json
{ "jid": "559981769536@s.whatsapp.net", "duration": 3600 }
```

`duration`: mute duration in seconds. `0` defaults to 8 hours.

### Unmute Chat

```
POST /sessions/:sessionId/chat/unmute
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

### Pin Chat

```
POST /sessions/:sessionId/chat/pin
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

### Unpin Chat

```
POST /sessions/:sessionId/chat/unpin
```

```json
{ "jid": "559981769536@s.whatsapp.net" }
```

---

## Labels

WhatsApp Business labels for organizing chats and messages.

### Add Label to Chat

```
POST /sessions/:sessionId/label/chat
```

```json
{ "jid": "559981769536@s.whatsapp.net", "labelId": "1" }
```

### Remove Label from Chat

```
POST /sessions/:sessionId/unlabel/chat
```

```json
{ "jid": "559981769536@s.whatsapp.net", "labelId": "1" }
```

### Add Label to Message

```
POST /sessions/:sessionId/label/message
```

```json
{
  "jid": "559981769536@s.whatsapp.net",
  "labelId": "1",
  "messageId": "3EB0FFFB732591EF45FD8B"
}
```

### Remove Label from Message

```
POST /sessions/:sessionId/unlabel/message
```

```json
{
  "jid": "559981769536@s.whatsapp.net",
  "labelId": "1",
  "messageId": "3EB0FFFB732591EF45FD8B"
}
```

### Edit Label

```
POST /sessions/:sessionId/label/edit
```

```json
{
  "labelId": "1",
  "name": "New Label Name",
  "color": 1,
  "deleted": false
}
```

---

## Newsletter

WhatsApp Channels (newsletters).

### Create Newsletter

```
POST /sessions/:sessionId/newsletter/create
```

```json
{
  "name": "My Channel",
  "description": "Channel description",
  "picture": "<base64-encoded-image>"
}
```

### Get Newsletter Info

```
POST /sessions/:sessionId/newsletter/info
```

```json
{ "newsletterJid": "120363407790725999@newsletter" }
```

### Get Newsletter Info from Invite Code

```
POST /sessions/:sessionId/newsletter/invite
```

```json
{ "code": "0029Vb7Muz0EVccClmotqu23" }
```

### List Subscribed Newsletters

```
GET /sessions/:sessionId/newsletter/list
```

### Get Newsletter Messages

```
POST /sessions/:sessionId/newsletter/messages
```

```json
{
  "newsletterJid": "120363407790725999@newsletter",
  "count": 20,
  "beforeId": 0
}
```

`beforeId`: message ID for pagination cursor. `0` = latest messages.

### Subscribe to Newsletter

```
POST /sessions/:sessionId/newsletter/subscribe
```

```json
{ "newsletterJid": "120363407790725999@newsletter" }
```

---

## Community

Communities are groups of groups in WhatsApp.

### Create Community

```
POST /sessions/:sessionId/community/create
```

```json
{
  "name": "My Community",
  "description": "Community description"
}
```

### Add Subgroup to Community

```
POST /sessions/:sessionId/community/participant/add
```

```json
{
  "communityJid": "120363426132481766@g.us",
  "participants": ["120362023605733675@g.us"]
}
```

### Remove Subgroup from Community

```
POST /sessions/:sessionId/community/participant/remove
```

```json
{
  "communityJid": "120363426132481766@g.us",
  "participants": ["120362023605733675@g.us"]
}
```

---

## Webhooks

Webhooks deliver real-time event notifications for a session via HTTP POST.

### Create Webhook

```
POST /sessions/:sessionId/webhooks
```

```json
{
  "url": "https://your-server.com/webhook",
  "secret": "optional-signing-secret",
  "events": ["Message", "Connected", "Disconnected"],
  "natsEnabled": false
}
```

See [Supported Event Types](#supported-event-types) for valid `events` values.

**Response** (201):
```json
{
  "data": {
    "id": "uuid",
    "sessionId": "session-uuid",
    "url": "https://your-server.com/webhook",
    "events": ["Message", "Connected", "Disconnected"],
    "enabled": true,
    "natsEnabled": false,
    "createdAt": "..."
  }
}
```

### List Webhooks

```
GET /sessions/:sessionId/webhooks
```

### Delete Webhook

```
DELETE /sessions/:sessionId/webhooks/:wid
```

---

## Webhook Payload

Events are published to NATS (`wzap.events.<sessionId>`) and delivered to registered webhook URLs.

| Field | Type | Description |
|---|---|---|
| `eventId` | string | Unique UUID for this event |
| `sessionId` | string | Session that emitted the event |
| `event` | string | Event type name |
| `timestamp` | string | ISO 8601 timestamp || _...native fields_ | any | All API event fields merged at top level |

---

## Supported Event Types

Use these values in webhook `events`. Use `"All"` to subscribe to everything.

| Category | Events |
|---|---|
| **Special** | `All` |
| **Messages** | `Message`, `UndecryptableMessage`, `MediaRetry`, `Receipt`, `DeleteForMe` |
| **Connection** | `Connected`, `Disconnected`, `ConnectFailure`, `LoggedOut`, `PairSuccess`, `PairError`, `QR`, `QRScannedWithoutMultidevice`, `StreamError`, `StreamReplaced`, `KeepAliveTimeout`, `KeepAliveRestored`, `ClientOutdated`, `TemporaryBan`, `CATRefreshError`, `ManualLoginReconnect` |
| **Contacts** | `Contact`, `Picture`, `IdentityChange`, `UserAbout`, `PushName`, `BusinessName` |
| **Groups** | `GroupInfo`, `JoinedGroup` |
| **Presence** | `Presence`, `ChatPresence` |
| **Chat State** | `Archive`, `Mute`, `Pin`, `Star`, `ClearChat`, `DeleteChat`, `MarkChatAsRead`, `UnarchiveChatsSetting` |
| **Labels** | `LabelEdit`, `LabelAssociationChat`, `LabelAssociationMessage` |
| **Calls** | `CallOffer`, `CallAccept`, `CallTerminate`, `CallOfferNotice`, `CallRelayLatency`, `CallPreAccept`, `CallReject`, `CallTransport`, `UnknownCallEvent` |
| **Newsletter** | `NewsletterJoin`, `NewsletterLeave`, `NewsletterMuteChange`, `NewsletterLiveUpdate` |
| **Sync** | `HistorySync`, `AppState`, `AppStateSyncComplete`, `AppStateSyncError`, `OfflineSyncCompleted`, `OfflineSyncPreview` |
| **Privacy** | `PrivacySettings`, `PushNameSetting`, `UserStatusMute`, `BlocklistChange`, `Blocklist` |
