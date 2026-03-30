# Template Map

This map helps you quickly choose which local dashboard-template file to inspect before changing Zenwoot frontend code.

## Main shell

- `.local/dashboard-template/app/app.vue`
  - global app wrapper
  - `UApp`, loading indicator, layout/page mounting

- `.local/dashboard-template/app/layouts/default.vue`
  - sidebar shell
  - top-level navigation groups
  - search button and command palette
  - footer/header slot patterns

## Core menus

- `.local/dashboard-template/app/components/TeamsMenu.vue`
  - workspace/team switcher behavior
  - compact vs expanded sidebar header behavior

- `.local/dashboard-template/app/components/UserMenu.vue`
  - user dropdown structure
  - theme controls
  - color and appearance switching

- `.local/dashboard-template/app/components/NotificationsSlideover.vue`
  - slideover structure
  - notifications-side interaction pattern

## Page references

- `.local/dashboard-template/app/pages/index.vue`
  - dashboard-style landing page
  - navbar + toolbar + panel body composition

- `.local/dashboard-template/app/pages/customers.vue`
  - table/list page reference
  - filters, actions, pagination, row menus

- `.local/dashboard-template/app/pages/inbox.vue`
  - message/inbox layout reference
  - useful for conversation-oriented split views

- `.local/dashboard-template/app/pages/settings.vue`
  - settings shell
  - settings navigation inside toolbar

## Home/dashboard components

- `.local/dashboard-template/app/components/home/HomeStats.vue`
  - stat cards
  - metric layout and icon treatment

- `.local/dashboard-template/app/components/home/HomeChart.client.vue`
  - chart rendering pattern
  - client-side visualization setup

- `.local/dashboard-template/app/components/home/HomeDateRangePicker.vue`
  - toolbar filters and date selection pattern

- `.local/dashboard-template/app/components/home/HomePeriodSelect.vue`
  - compact segmented period selector

- `.local/dashboard-template/app/components/home/HomeSales.vue`
  - data table/card hybrid presentation

## Official Nuxt UI validation

Before creating a new component or changing how a Nuxt UI component is used, consult the official Nuxt UI docs through the MCP server.

Use the official MCP lookups to:
 
- confirm the correct component to use
- confirm official usage and examples
- confirm props, slots, events, and API behavior
- avoid guessing component capabilities from memory alone

## Domain mappings for Zenwoot

Use these translations when adapting the template to product terminology:

- `Home` -> `Conversations` or Zenwoot landing route
- `Inbox` -> conversation workspace / inbox-driven messaging UI
- `Customers` -> `Contacts`
- `Members` -> `Agents` or `Team members`
- `Settings` remains `Settings`

## How to use this map

1. Identify the kind of screen or component you need to change.
2. Read the closest file from `.local/dashboard-template`.
3. If components are involved, read the official Nuxt UI docs through MCP.
4. Read the matching Zenwoot backend route and handler.
5. Only then edit `frontend/`.
