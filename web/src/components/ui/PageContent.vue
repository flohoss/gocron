<script setup lang="ts">
defineProps<{ background?: string }>();
</script>

<template>
  <div class="px-4 py-2 md:p-5 lg:px-10" :class="background && background">
    <slot></slot>
  </div>
</template>
