<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { postJob, postJobs } from '../client/sdk.gen';
import { useJobs } from '../stores/useJobs';
import JobSelectModal from './utils/JobSelectModal.vue';
import IconTerminal from '~icons/fa7-solid/terminal';
import ChevronLeft from '~icons/fa7-solid/chevron-left';
import ListCheck from '~icons/fa7-solid/list-check';
import Play from '~icons/fa7-solid/play';
import OpenApi from '~icons/simple-icons/openapiinitiative';

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
    return 'No Jobs Selected';
  } else if (currentJob.value === null && jobsUnchecked.value) {
    return 'Run Selected Jobs';
  }
  return 'Run ' + (currentJob.value !== null ? currentJob.value.name : 'All Jobs');
});

const showExtraButtons = computed(() => currentJob.value === null);
</script>

<template>
  <header class="mx-auto mb-4 md:mb-10 relative max-w-3xl flex justify-center">
    <div class="absolute top-1/2 -translate-y-1/2 left-3">
      <div v-if="$route.name !== 'homeView'" class="tooltip" data-tip="Back" data-test-id="back-button">
        <button @click="router.push('/')" class="btn btn-soft btn-circle">
          <ChevronLeft />
        </button>
      </div>
      <div v-else class="join">
        <div class="tooltip" data-tip="Terminal" data-test-id="terminal-button">
          <button @click="router.push('/commands')" class="btn px-[0.7rem] btn-soft join-item rounded-l-full">
            <IconTerminal />
          </button>
        </div>
        <div data-test-id="openapi-button" class="tooltip" data-tip="OpenAPI Documentation">
          <a href="/api/docs" class="btn px-3 btn-soft btn-secondary join-item rounded-r-full">
            <OpenApi />
          </a>
        </div>
      </div>
    </div>

    <img class="h-24 lg:h-36" src="/static/logo.webp" />

    <Transition mode="out-in">
      <div class="join absolute top-1/2 -translate-y-1/2 right-3" v-if="$route.name !== 'commandView'">
        <div class="tooltip" data-tip="Select Jobs" v-if="showExtraButtons" data-test-id="select-button">
          <button
            onclick="selectModal.showModal()"
            class="btn px-3 btn-soft rounded-l-full"
            :class="[jobsUnchecked ? 'btn-primary' : 'btn-secondary', showExtraButtons ? 'join-item' : '']"
            :disabled="disabled"
          >
            <ListCheck />
          </button>
        </div>
        <div class="tooltip" :data-tip="playLabel" data-test-id="run-button">
          <button
            @click="run"
            class="btn px-[0.6rem] btn-soft"
            :disabled="disabled || checked.length === 0"
            :class="showExtraButtons ? 'join-item rounded-r-full' : 'btn-circle'"
          >
            <Play v-if="!disabled || loading" />
            <span v-else class="loading loading-spinner w-[1.2rem]"></span>
          </button>
        </div>
      </div>
    </Transition>
    <JobSelectModal />
  </header>
</template>
