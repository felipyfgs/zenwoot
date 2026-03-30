<script setup lang="ts">
import type { User } from '~/types'

defineProps<{
  members: User[]
}>()
</script>

<template>
  <ul role="list" class="divide-y divide-default">
    <li
      v-for="member in members"
      :key="member.id"
      class="flex items-center gap-3 py-3 px-4 sm:px-6"
    >
      <UAvatar :alt="member.displayName" size="md" />
      <div class="text-sm min-w-0 flex-1">
        <p class="text-highlighted font-medium truncate">
          {{ member.displayName }}
        </p>
        <p class="text-muted truncate">
          {{ member.email }}
        </p>
      </div>
      <UBadge variant="subtle" size="xs">
        {{ member.role }}
      </UBadge>
    </li>
  </ul>
</template>
