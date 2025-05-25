<script setup lang="ts">
import { useEventSource } from '@vueuse/core';
import { computed, watch } from 'vue';
import { BackendURL } from '../main';
import { useEventStore } from '../stores/event';
import { PlayIcon, ChevronLeftIcon, InformationCircleIcon } from '@heroicons/vue/24/outline';
import { useRouter } from 'vue-router';
import VersionDialog from './VersionDialog.vue';
import { postJob, postJobs } from '../client/sdk.gen';

const router = useRouter();
const store = useEventStore();

const { data, close } = useEventSource(BackendURL + '/api/events?stream=status', [], {
  autoReconnect: {
    retries: 10,
    delay: 100,
    onFailed() {
      alert('Failed to connect EventSource after 10 retries');
    },
  },
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

const playLabel = computed(() => 'run ' + (store.currentJobId !== null ? store.currentJobId : 'all jobs'));

const disabled = computed(() => store.state.loading || !store.idle);
</script>

<template>
  <header class="flex justify-between items-center md:justify-center md:gap-20 mb-4 md:mb-10 mx-3">
    <div v-if="$route.name !== 'homeView'" class="tooltip" data-tip="back">
      <button @click="router.push('/')" class="btn btn-soft btn-circle">
        <ChevronLeftIcon class="size-6" />
      </button>
    </div>
    <div v-else class="tooltip" data-tip="software details">
      <button onclick="version_modal.showModal()" class="btn btn-soft btn-circle">
        <InformationCircleIcon class="size-6" />
      </button>
    </div>

    <img class="h-28 lg:h-36" src="/static/logo.webp" />

    <div class="tooltip" :data-tip="playLabel">
      <button @click="run" class="btn btn-soft btn-circle" :disabled="disabled">
        <PlayIcon v-if="!disabled || store.state.loading" class="size-6" />
        <span v-else class="loading loading-spinner"></span>
      </button>
    </div>

    <VersionDialog id="version_modal" />
  </header>
</template>
