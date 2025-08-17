<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { postJob, postJobs } from '../client/sdk.gen';
import { useJobs } from '../stores/useJobs';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTerminal, faChevronLeft, faPlay, faListCheck } from '@fortawesome/free-solid-svg-icons';
import JobSelectModal from './utils/JobSelectModal.vue';

const { disabled, loading, currentJob, checked, jobsUnchecked } = useJobs();
const router = useRouter();

const run = async () => {
  if (currentJob.value === null) {
    await postJobs();
  } else {
    await postJob({ path: { name: currentJob.value.name } });
  }
};

const playLabel = computed(() => {
  if (checked.value.length === 0) {
    return 'no job selected';
  } else if (jobsUnchecked.value) {
    return 'run selected jobs';
  }
  return 'run ' + (currentJob.value !== null ? currentJob.value.name : 'all jobs');
});
</script>

<template>
  <header class="flex justify-between items-center md:justify-center md:gap-20 mb-4 md:mb-10 mx-3">
    <div v-if="$route.name !== 'homeView'" class="tooltip" data-tip="back">
      <button @click="router.push('/')" class="btn btn-soft btn-circle">
        <FontAwesomeIcon :icon="faChevronLeft" />
      </button>
    </div>
    <div v-else class="tooltip" data-tip="execute command">
      <button @click="router.push('/commands')" class="btn btn-soft btn-circle">
        <FontAwesomeIcon :icon="faTerminal" />
      </button>
    </div>

    <img class="h-28 lg:h-36" src="/static/logo.webp" />

    <div class="join">
      <div class="tooltip" :data-tip="playLabel">
        <button @click="run" class="btn join-item btn-soft rounded-l-full" :disabled="disabled || checked.length === 0">
          <FontAwesomeIcon v-if="!disabled || loading" :icon="faPlay" />
          <span v-else class="loading loading-spinner"></span>
        </button>
      </div>
      <div class="tooltip" data-tip="select jobs">
        <button
          onclick="selectModal.showModal()"
          class="btn join-item px-2 rounded-r-full"
          :class="jobsUnchecked ? 'btn-primary' : 'btn-soft btn-secondary'"
          :disabled="disabled"
        >
          <FontAwesomeIcon :icon="faListCheck" />
        </button>
      </div>
    </div>
    <JobSelectModal />
  </header>
</template>
