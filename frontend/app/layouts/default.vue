<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

useAuth()
const open = ref(false)

function closeNav() {
  open.value = false
}

const links = [[{
  label: 'Conversations',
  icon: 'i-lucide-message-circle',
  to: '/conversations',
  onSelect: closeNav
}, {
  label: 'Contacts',
  icon: 'i-lucide-users',
  to: '/contacts',
  onSelect: closeNav
}, {
  label: 'Settings',
  to: '/settings',
  icon: 'i-lucide-settings',
  defaultOpen: true,
  type: 'trigger' as const,
  children: [{
    label: 'General',
    to: '/settings',
    exact: true,
    onSelect: closeNav
  }, {
    label: 'Inboxes',
    to: '/settings/inboxes',
    onSelect: closeNav
  }, {
    label: 'Agents',
    to: '/settings/agents',
    onSelect: closeNav
  }, {
    label: 'Teams',
    to: '/settings/teams',
    onSelect: closeNav
  }, {
    label: 'Labels',
    to: '/settings/labels',
    onSelect: closeNav
  }, {
    label: 'Canned Responses',
    to: '/settings/canned-responses',
    onSelect: closeNav
  }]
}]] satisfies NavigationMenuItem[][]

const groups = computed(() => [{
  id: 'links',
  label: 'Go to',
  items: links.flat()
}])
</script>

<template>
  <UDashboardGroup unit="rem">
    <UDashboardSidebar
      id="default"
      v-model:open="open"
      collapsible
      resizable
      class="bg-elevated/25"
      :ui="{ footer: 'lg:border-t lg:border-default' }"
    >
      <template #header="{ collapsed }">
        <TeamsMenu :collapsed="collapsed" />
      </template>

      <template #default="{ collapsed }">
        <UDashboardSearchButton :collapsed="collapsed" class="bg-transparent ring-default" />

        <UNavigationMenu
          :collapsed="collapsed"
          :items="links[0]"
          orientation="vertical"
          tooltip
          popover
        />
      </template>

      <template #footer="{ collapsed }">
        <UserMenu :collapsed="collapsed" />
      </template>
    </UDashboardSidebar>

    <UDashboardSearch :groups="groups" />

    <slot />
  </UDashboardGroup>
</template>
