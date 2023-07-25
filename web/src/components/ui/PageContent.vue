<script setup lang="ts">
defineProps<{ background?: string }>();
</script>

<template>
  <div class="grid grid-cols-1 px-4 py-3 md:p-5 lg:px-10" :class="background && background">
    <slot></slot>
  </div>
</template>
