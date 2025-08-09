<script setup lang="ts">
import { onUpdated, useTemplateRef } from 'vue';

defineProps<{ title?: string }>();

const scrollContainer = useTemplateRef('scrollContainer');

onUpdated(() => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight;
  }
});
</script>

<template>
  <div class="bg-base-300 w-full text-sm rounded-xl">
    <div class="flex justify-between gap-5 bg-base-200 rounded-t-xl padding">
      <div class="flex items-center gap-2">
        <div class="console-btn bg-error text-error hover:text-error-content"></div>
        <div class="console-btn bg-warning text-warning hover:text-warning-content"></div>
        <div class="console-btn bg-success text-success hover:text-success-content"></div>
      </div>
      <div v-if="title" class="text-secondary font-bold truncate max-w-full">{{ title }}</div>
    </div>
    <div ref="scrollContainer" class="flex flex-col h-[calc(100vh-12rem)] md:h-[calc(100vh-14rem)] lg:h-[calc(100vh-18rem)] overflow-scroll p-4">
      <slot></slot>
      <div class="mt-auto">
        <slot name="bottom"></slot>
      </div>
    </div>
  </div>
</template>
