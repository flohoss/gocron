<script setup lang="ts">
defineProps<{ badges: {} }>();
</script>

<template>
  <div class="badge whitespace-nowrap min-w-0" v-for="(badge, key, counter) of badges" :key="key">
    <span class="opacity-50 mr-2">{{ Object.keys(badges)[counter] }}</span> <span class="truncate">{{ badge || '-' }}</span>
  </div>
</template>
