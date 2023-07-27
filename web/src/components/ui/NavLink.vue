<script setup lang="ts">
import { RouterLink } from 'vue-router';

defineProps<{
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
        <div class="grid grid-cols-1">
          <div class="text-lg truncate">{{ name }}</div>
          <div v-if="extra" class="opacity-50 text-xs truncate">{{ extra }}</div>
        </div>
      </div>
      <div class="pr-2">
        <slot></slot>
      </div>
    </RouterLink>
  </li>
</template>
