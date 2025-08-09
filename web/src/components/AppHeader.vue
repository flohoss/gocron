<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { postJob, postJobs } from '../client/sdk.gen';
import { useJobs } from '../stores/useJobs';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTerminal, faChevronLeft, faPlay } from '@fortawesome/free-solid-svg-icons';

const { disabled, loading, currentJob } = useJobs();
const router = useRouter();

const run = async () => {
  if (currentJob.value === null) {
    await postJobs();
  } else {
    await postJob({ path: { name: currentJob.value.name } });
  }
};

const playLabel = computed(() => 'run ' + (currentJob.value !== null ? currentJob.value.name : 'all jobs'));
</script>

<template>
  <header class="flex justify-between items-center md:justify-center md:gap-20 mb-4 md:mb-10 mx-3">
    <div v-if="$route.name !== 'homeView'" class="tooltip" data-tip="back">
      <button @click="router.push('/')" class="btn btn-soft btn-circle">
        <FontAwesomeIcon :icon="faChevronLeft" />
      </button>
    </div>
    <div v-else class="tooltip" data-tip="execute command">
      <button onclick="commandModal.showModal()" class="btn btn-soft btn-circle">
        <FontAwesomeIcon :icon="faTerminal" />
      </button>
    </div>

    <img class="h-28 lg:h-36" src="/static/logo.webp" />

    <div class="tooltip" :data-tip="playLabel">
      <button @click="run" class="btn btn-soft btn-circle" :disabled="disabled">
        <FontAwesomeIcon v-if="!disabled || loading" :icon="faPlay" />
        <span v-else class="loading loading-spinner"></span>
      </button>
    </div>
  </header>
</template>
