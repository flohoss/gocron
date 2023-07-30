<script setup lang="ts">
defineProps<{ badges: {} }>();
</script>

<template>
  <div class="badge whitespace-nowrap" v-for="(badge, key, counter) of badges" :key="key">
    <span class="opacity-50 mr-2">{{ Object.keys(badges)[counter] }}</span> {{ badge || '-' }}
  </div>
</template>
