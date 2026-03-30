---
name: zenwoot-frontend-study
description: Study and implement Zenwoot frontend work by always consulting the local Nuxt UI Dashboard reference at .local/dashboard-template and the real backend contracts in backend/internal/server before editing pages, components, layouts, composables, or frontend UX.
---

# Zenwoot Frontend Study

Use this skill whenever the task involves studying, creating, editing, refactoring, or reviewing frontend code in `frontend/`.

The canonical visual reference is the local copy of the official Nuxt UI Dashboard template at:

```text
/home/obsidian/dev/zenwoot/.local/dashboard-template
```

The canonical product and data-contract reference is the Zenwoot backend at:

```text
/home/obsidian/dev/zenwoot/backend
```

## Core rule

Before making any frontend change, always do both:

1. Consult the closest matching file in `.local/dashboard-template`
2. Consult the relevant backend route, handler, and response shape in `backend/internal/server`

Do not skip either side.

When the task involves validating a Nuxt UI component, selecting between Nuxt UI components, or creating a new component that uses Nuxt UI primitives, also consult the official Nuxt UI documentation through the Nuxt UI MCP server before editing.
Use the MCP tools to inspect the relevant official component docs and metadata, especially component usage, API, props, slots, and examples.

## What this skill is for

- Study the architecture and design system of the dashboard template
- Recreate Zenwoot UI on top of the dashboard template patterns
- Add or edit pages, components, layouts, tables, forms, panels, toolbars, and navigation
- Align frontend behavior with backend routes, payloads, pagination, and websocket flows
- Keep frontend work reusable and consistent with Nuxt UI dashboard conventions

## Required workflow

### 1. Identify the target surface

Classify the request first:

- Layout or shell
- Navigation or settings structure
- Dashboard panel or split-view screen
- Table or list page
- Detail panel or slideover
- Form or CRUD screen
- API/composable/type work

### 2. Consult the local dashboard reference

Read the nearest equivalent under `.local/dashboard-template` before editing.

Typical reference files:

- `app/layouts/default.vue`
- `app/pages/index.vue`
- `app/pages/customers.vue`
- `app/pages/settings.vue`
- `app/components/UserMenu.vue`
- `app/components/TeamsMenu.vue`
- `app/components/home/HomeStats.vue`
- `app/components/home/HomeChart.client.vue`
- `app/components/home/HomeSales.vue`

Also read `TEMPLATE_MAP.md` in this skill folder for a quick mapping.

### 2.5. Consult the official Nuxt UI docs through MCP when components are involved

If you are creating, replacing, validating, or composing Nuxt UI components, consult the official Nuxt UI MCP tools before implementation.

Typical checks:

- use component documentation lookup for official usage and examples
- use component metadata lookup for props, slots, and events
- use documentation page lookup when deciding between related components or patterns

This is mandatory for new components and for any edit that changes component API usage.

### 3. Consult the backend contract

Before wiring data, inspect the real API in Zenwoot:

- `backend/internal/server/router.go`
- relevant handlers in `backend/internal/handler/*.go`
- DTOs or request/response structures when present in `backend/internal/dto`

Assume backend responses use the standard envelope:

```json
{
  "success": true,
  "data": {},
  "error": "",
  "message": "success"
}
```

The frontend must unwrap `data` and handle unsuccessful responses consistently.

### 4. Rebuild using template patterns, not template mock data

Follow the template's:

- layout composition
- Nuxt UI component usage
- spacing and visual rhythm
- panel and navbar structure
- table composition
- navigation grouping
- slideover and responsive behavior
- theme and color handling

But do **not** rely on the template's fake `server/api` data as source of truth.
Replace template mocks with real Zenwoot backend integrations.

### 5. Prefer reusable Zenwoot abstractions

When implementing:

- extract reusable components instead of page-local duplication
- create composables for API access and stateful frontend logic
- define or update shared types in `frontend/app/types`
- keep routes and naming aligned with Zenwoot domain terms such as conversations, contacts, inboxes, labels, teams, and canned responses

## Zenwoot-specific guidance

### Navigation and shell

Zenwoot should keep the dashboard-template shell but adapt the information architecture to product routes like:

- `/conversations`
- `/contacts`
- `/settings`
- `/settings/inboxes`
- `/settings/agents`
- `/settings/teams`
- `/settings/labels`
- `/settings/canned-responses`

### Conversations

For conversation-heavy screens:

- prefer `UDashboardPanel` split layouts
- use responsive side panels and `USlideover` on small screens
- align list/detail interactions with inbox and status filters
- inspect backend conversation endpoints before adding actions like read, assign, labels, snooze, or priority

### Contacts

For contacts screens:

- use the customers table patterns from the template as structure inspiration
- replace template customer fields with Zenwoot contact fields
- validate pagination, search, and row actions against backend contacts routes

### Settings

For settings screens:

- keep the template's settings navigation and toolbar patterns
- adapt pages to Zenwoot entities and backend capabilities
- build forms and management views around the real routes for inboxes, users, teams, labels, and canned responses

## Non-negotiable constraints

- Always read `.local/dashboard-template` before changing frontend UI
- Always read backend routes and handlers before binding data or actions
- Never invent backend endpoints or payloads
- Never use template mock APIs as final implementation
- Keep changes visually consistent with the dashboard template
- Favor reusable components and composables over one-off page logic
- Always consult official Nuxt UI documentation through MCP when component behavior or API usage is changed

## Suggested read order

1. `TEMPLATE_MAP.md`
2. matching file under `.local/dashboard-template`
3. official Nuxt UI MCP docs if components are involved
4. `backend/internal/server/router.go`
5. matching backend handler(s)
6. target file(s) in `frontend/`

## Done criteria

A frontend change is only complete when:

- the design clearly matches dashboard template patterns
- official Nuxt UI docs were consulted through MCP when component behavior or API usage was changed
- the data flow matches the real backend contract
- the code is reusable and aligned with Zenwoot naming
- mobile and desktop behavior are both considered
- loading, empty, and error states are not ignored
