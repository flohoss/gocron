<script setup lang="ts">
import PageContent from './PageContent.vue';
</script>

<template>
  <PageContent background="bg-base-300" class="sticky top-0">
    <div class="flex justify-between items-center min-w-0">
      <slot></slot>
    </div>
  </PageContent>
</template>
