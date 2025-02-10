<script setup lang="ts">
import { useEventSource } from '@vueuse/core';
import { onBeforeUnmount, watch } from 'vue';
import { BackendURL } from '../main';
import { useEventStore } from '../stores/event';
import { PlayIcon, ChevronLeftIcon } from '@heroicons/vue/24/outline';
import { useRouter } from 'vue-router';
import { JobsService } from '../openapi';

const router = useRouter();
const store = useEventStore();

const { data, close } = useEventSource(BackendURL + 'api/events?stream=status', [], {
  autoReconnect: true,
});
onBeforeUnmount(() => close());
watch(() => data.value, store.parseEventInfo);

const run = async () => {
  if (store.currentJobId === null) {
    await JobsService.postJobs();
  } else {
    await JobsService.postJobs1(store.currentJobId);
  }
};
</script>

<template>
  <header class="flex justify-between mx-4 items-center md:justify-center md:gap-20 mb-5 md:mb-10">
    <button @click="router.push('/')" class="btn btn-soft btn-circle">
      <ChevronLeftIcon class="size-6" />
    </button>
    <img class="size-28 lg:size-36" src="/logo/logo.webp" />
    <button @click="run" class="btn btn-soft btn-circle" :disabled="!store.idle">
      <PlayIcon v-if="store.idle" class="size-6" />
      <span v-else class="loading loading-spinner"></span>
    </button>
  </header>
</template>
