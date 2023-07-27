<script setup lang="ts">
import PageContent from './PageContent.vue';
</script>

<template>
  <PageContent background="bg-base-300" class="sticky top-0">
    <slot></slot>
  </PageContent>
</template>
