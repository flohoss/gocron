<script setup lang="ts">
import { computed } from 'vue';
import { RouterLink } from 'vue-router';
import { useWindowSize } from '@vueuse/core';
import type { JobView } from '../client/types.gen';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faCheck, faQuestion, faTimes } from '@fortawesome/free-solid-svg-icons';

const props = defineProps<{ job: JobView }>();
const url = computed<string>(() => '/jobs/' + props.job.name);

const { width } = useWindowSize();
const isMobile = computed(() => width.value < 1024);
const runs = computed(() => {
  const runsArray = props.job?.runs ?? [];
  const amount = runsArray.length;

  if (amount === 0) return null;

  if (isMobile.value) {
    return [runsArray[amount - 1]];
  } else {
    return runsArray.slice(-Math.min(3, amount));
  }
});

enum Status {
  Running = 1,
  Stopped = 2,
  Finished = 3,
}

function getStepColor(status: Status): string {
  switch (status) {
    case Status.Running:
      return 'step-warning';
    case Status.Stopped:
      return 'step-error';
    case Status.Finished:
      return 'step-success';
    default:
      return 'step-neutral';
  }
}

function getStepIcon(status: Status) {
  switch (status) {
    case Status.Stopped:
      return faTimes;
    case Status.Finished:
      return faCheck;
    default:
      return faQuestion;
  }
}
</script>

<template>
  <RouterLink class="flex justify-between items-center group last:mb-8 lg:last:mb-0 hover:cursor-pointer" :to="url">
    <div class="pl-4 truncate">
      <div class="group-hover:text-primary hover-animation text-2xl font-medium truncate">{{ job.name }}</div>
      <div class="text-secondary text-sm truncate">{{ job.cron }}</div>
    </div>
    <div class="text-sm">
      <ul class="steps" v-if="runs">
        <li v-for="run in runs" :key="run.id" class="step" :class="getStepColor(run.status_id)">
          <span class="step-icon">
            <span v-if="run.status_id === Status.Running" class="loading loading-spinner"></span>
            <FontAwesomeIcon v-else :icon="getStepIcon(run.status_id)" class="size-6" />
          </span>
          {{ run.duration }}
        </li>
      </ul>
    </div>
  </RouterLink>
</template>

<style scoped>
.steps .step::before {
  height: 0.2rem !important;
}
</style>
