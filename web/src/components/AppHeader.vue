<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { postJob, postJobs } from '../client/sdk.gen';
import { useJobs } from '../stores/useJobs';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTerminal, faChevronLeft, faPlay, faListCheck } from '@fortawesome/free-solid-svg-icons';
import JobSelectModal from './utils/JobSelectModal.vue';
import { faGithub } from '@fortawesome/free-brands-svg-icons';

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

const showExtraButtons = computed(() => currentJob.value === null);
</script>

<template>
  <header class="mx-auto mb-4 md:mb-10 relative max-w-3xl flex justify-center">
    <div class="absolute top-1/2 -translate-y-1/2 left-3">
      <div v-if="$route.name !== 'homeView'" class="tooltip" data-tip="back">
        <button @click="router.push('/')" class="btn btn-soft btn-circle">
          <FontAwesomeIcon :icon="faChevronLeft" />
        </button>
      </div>
      <div v-else class="join">
        <div class="tooltip" data-tip="execute command">
          <button @click="router.push('/commands')" class="btn px-[0.6rem] btn-soft join-item rounded-l-full">
            <FontAwesomeIcon :icon="faTerminal" />
          </button>
        </div>
        <div class="tooltip" data-tip="github repository">
          <a href="https://github.com/flohoss/gocron" target="_blank" class="btn px-3 btn-soft btn-secondary join-item rounded-r-full">
            <FontAwesomeIcon :icon="faGithub" />
          </a>
        </div>
      </div>
    </div>

    <img class="h-24 lg:h-36" src="/static/logo.webp" />

    <Transition mode="out-in">
      <div class="join absolute top-1/2 -translate-y-1/2 right-3" v-if="$route.name !== 'commandView'">
        <div class="tooltip" data-tip="select jobs" v-if="showExtraButtons">
          <button
            onclick="selectModal.showModal()"
            class="btn px-3 btn-soft rounded-l-full"
            :class="[jobsUnchecked ? 'btn-primary' : 'btn-secondary', showExtraButtons ? 'join-item' : '']"
            :disabled="disabled"
          >
            <FontAwesomeIcon :icon="faListCheck" />
          </button>
        </div>
        <div class="tooltip" :data-tip="playLabel">
          <button
            @click="run"
            class="btn px-[0.6rem] join-item btn-soft rounded-r-full"
            :disabled="disabled || checked.length === 0"
            :class="showExtraButtons ? 'join-item rounded-r-full' : 'btn-circle'"
          >
            <FontAwesomeIcon v-if="!disabled || loading" :icon="faPlay" />
            <span v-else class="loading loading-spinner"></span>
          </button>
        </div>
      </div>
    </Transition>
    <JobSelectModal />
  </header>
</template>
