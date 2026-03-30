<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed } from 'vue'

const route = useRoute()

const isNarrowSettingsPage = computed(() => route.path === '/settings')
const isInboxWizardPage = computed(() => route.path === '/settings/inbox-new')

const links = [[{
  label: 'General',
  icon: 'i-lucide-settings',
  to: '/settings',
  exact: true
}, {
  label: 'Inboxes',
  icon: 'i-lucide-smartphone',
  to: '/settings/inboxes'
}, {
  label: 'Agents',
  icon: 'i-lucide-user-check',
  to: '/settings/agents'
}, {
  label: 'Teams',
  icon: 'i-lucide-users-2',
  to: '/settings/teams'
}, {
  label: 'Labels',
  icon: 'i-lucide-tag',
  to: '/settings/labels'
}, {
  label: 'Canned Responses',
  icon: 'i-lucide-message-square-text',
  to: '/settings/canned-responses'
}]] satisfies NavigationMenuItem[][]
</script>

<template>
  <UDashboardPanel v-if="isInboxWizardPage" id="settings" :ui="{ body: 'p-0 sm:p-0' }">
    <template #body>
      <NuxtPage />
    </template>
  </UDashboardPanel>

  <UDashboardPanel v-else id="settings" :ui="{ body: 'lg:py-12' }">
    <template #header>
      <UDashboardNavbar title="Settings">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
      </UDashboardNavbar>

      <UDashboardToolbar>
        <!-- NOTE: The `-mx-1` class is used to align with the `DashboardSidebarCollapse` button here. -->
        <UNavigationMenu :items="links" highlight class="-mx-1 flex-1" />
      </UDashboardToolbar>
    </template>

    <template #body>
      <div
        :class="[
          'flex w-full flex-col gap-4 sm:gap-6 lg:gap-12',
          isNarrowSettingsPage ? 'mx-auto lg:max-w-2xl' : ''
        ]"
      >
        <NuxtPage />
      </div>
    </template>
  </UDashboardPanel>
</template>
