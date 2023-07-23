<script setup lang="ts">
defineProps<{ background?: string }>();
</script>

<template>
  <div class="p-2 lg:p-5" :class="background && background">
    <slot></slot>
  </div>
</template>
