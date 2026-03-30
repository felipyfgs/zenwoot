<script setup lang="ts">
defineProps<{
  collapsed?: boolean
}>()
</script>

<template>
  <div class="flex items-center gap-2 px-2 py-1.5">
    <div class="size-8 rounded-lg bg-primary flex items-center justify-center shrink-0">
      <UIcon name="i-lucide-zap" class="size-4 text-white" />
    </div>
    <span v-if="!collapsed" class="font-semibold text-highlighted truncate">Zenwoot</span>
  </div>
</template>
