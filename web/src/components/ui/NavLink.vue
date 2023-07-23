<script setup lang="ts">
import { RouterLink } from 'vue-router';

const props = defineProps<{
  link: { name: string; params?: { id?: number } };
  name: string;
  extra?: string;
  icon: string;
  status?: string;
  active: boolean;
  background?: 'bg-primary' | 'bg-secondary';
  text?: 'text-primary-content' | 'text-secondary-content';
  smallHidden?: boolean;
}>();
</script>

<template>
  <li :class="{ 'hidden lg:flex': smallHidden }">
    <RouterLink class="flex justify-between px-2" :class="{ active: active }" :to="link">
      <div class="flex items-center gap-2">
        <div
          class="w-10 h-10 p-2 rounded-full flex items-center justify-center shrink-0 text-lg"
          :class="background ? background + ' ' + text : 'bg-base-100 text-base-content'"
        >
          <span v-html="icon"></span>
        </div>
        <div class="flex flex-col">
          <div class="text-lg">{{ name }}</div>
          <div v-if="extra" class="opacity-50 text-xs">{{ extra }}</div>
        </div>
      </div>
      <div v-if="status" class="relative flex h-3 w-3">
        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
        <span class="relative inline-flex rounded-full h-3 w-3 bg-success"></span>
      </div>
    </RouterLink>
  </li>
</template>
