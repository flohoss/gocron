<script setup lang="ts">
import { useEventSource } from '@vueuse/core';
import { watch } from 'vue';
import { BackendURL } from '../main';
import { useEventStore } from '../stores/event';
import { PlayIcon, ChevronLeftIcon, InformationCircleIcon } from '@heroicons/vue/24/outline';
import { useRouter } from 'vue-router';
import VersionDialog from './VersionDialog.vue';
import { postJob, postJobs } from '../client/sdk.gen';

const router = useRouter();
const store = useEventStore();

const { data, close } = useEventSource(BackendURL + '/api/events?stream=status', [], {
  autoReconnect: true,
});
addEventListener('beforeunload', () => {
  close();
});
watch(() => data.value, store.parseEventInfo);

const run = async () => {
  if (store.currentJobId === null) {
    await postJobs();
  } else {
    await postJob({ path: { name: store.currentJobId } });
  }
};
</script>

<template>
  <header class="flex justify-between items-center md:justify-center md:gap-20 mb-4 md:mb-8 mx-3">
    <button v-if="$route.name !== 'homeView'" @click="router.push('/')" class="btn btn-soft btn-circle">
      <ChevronLeftIcon class="size-6" />
    </button>
    <button v-else onclick="version_modal.showModal()" class="btn btn-soft btn-circle">
      <InformationCircleIcon class="size-6" />
    </button>

    <img class="size-28 lg:size-36" src="/static/logo.webp" />
    <button @click="run" class="btn btn-soft btn-circle" :disabled="!store.idle">
      <PlayIcon v-if="store.idle" class="size-6" />
      <span v-else class="loading loading-spinner"></span>
    </button>

    <VersionDialog id="version_modal" />
  </header>
</template>
