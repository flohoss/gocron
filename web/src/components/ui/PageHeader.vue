<script setup lang="ts">
import PageContent from './PageContent.vue';
defineProps<{ badges?: { name: string; icon: string }[] }>();
</script>

<template>
  <PageContent background="bg-base-300">
    <div class="flex flex-col gap-2 lg:flex-row justify-between items-center">
      <slot></slot>
    </div>
    <div class="mt-5 flex items-center flex-wrap gap-2">
      <div v-for="(badge, i) in badges" :key="i" class="badge badge-ghost badge-sm gap-1"><i :class="'fa-solid fa-' + badge.icon"></i>{{ badge.name }}</div>
    </div>
  </PageContent>
</template>
