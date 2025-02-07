<script setup lang="ts">
import { useEventStore } from '../stores/event';
import { PlayIcon, ChevronLeftIcon } from '@heroicons/vue/24/outline';

const event = useEventStore();
</script>

<template>
  <header class="flex justify-between mx-4 items-center md:justify-center md:gap-20 mb-5 md:mb-10">
    <a href="/" class="btn btn-soft btn-circle">
      <ChevronLeftIcon class="size-6" />
    </a>
    <img class="size-28 lg:size-36" src="/logo/logo.webp" />
    <button @click="event.run()" class="btn btn-soft btn-circle" :disabled="!event.idle">
      <PlayIcon v-if="event.idle" class="size-6" />
      <span v-else class="loading loading-spinner"></span>
    </button>
  </header>
</template>
