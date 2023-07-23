<script setup lang="ts">
import PageContent from "./PageContent.vue";
</script>

<template>
  <PageContent background="bg-base-300">
    <div class="flex flex-col gap-2 lg:flex-row justify-between items-center lg:items-end">
      <slot></slot>
    </div>
  </PageContent>
</template>
