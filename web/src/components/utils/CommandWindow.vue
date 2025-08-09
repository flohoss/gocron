<script setup lang="ts">
import { onUpdated, useTemplateRef } from 'vue';

const scrollContainer = useTemplateRef('scrollContainer');

onUpdated(() => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight;
  }
});
</script>

<template>
  <div class="bg-base-300 w-full text-sm rounded-xl">
    <div class="flex justify-start items-center gap-2 padding bg-base-200 rounded-t-xl">
      <div class="console-btn bg-error text-error hover:text-error-content"></div>
      <div class="console-btn bg-warning text-warning hover:text-warning-content"></div>
      <div class="console-btn bg-success text-success hover:text-success-content"></div>
    </div>
    <div ref="scrollContainer" class="flex flex-col h-[calc(100vh-12rem)] md:h-[calc(100vh-14rem)] lg:h-[calc(100vh-18rem)] overflow-scroll p-4">
      <slot></slot>
      <div class="mt-auto">
        <slot name="bottom"></slot>
      </div>
    </div>
  </div>
</template>
