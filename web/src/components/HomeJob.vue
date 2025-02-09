<script setup lang="ts">
import { computed } from 'vue';
import { type jobs_JobsView } from '../openapi';
import { RouterLink } from 'vue-router';
import { useWindowSize } from '@vueuse/core';
import ShortDuration from './ShortDuration.vue';

const props = defineProps<{ job: jobs_JobsView }>();
const url = computed<string>(() => '/jobs/' + props.job.id);

const { width } = useWindowSize();
const isMobile = computed(() => width.value < 1024);
const runs = computed(() => {
  if (isMobile.value) {
    const amount = props.job.runs?.length || 0;
    if (amount > 0) {
      return [props.job.runs[amount - 1]];
    }
    return null;
  } else {
    return props.job.runs;
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

function getStepIcon(status: Status): string {
  switch (status) {
    case Status.Running:
      return '●';
    case Status.Stopped:
      return '✕';
    case Status.Finished:
      return '✓';
    default:
      return '?';
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
        <li v-for="run in runs" :key="run.id" :data-content="getStepIcon(run.status_id)" class="step" :class="getStepColor(run.status_id)">
          <ShortDuration v-if="run.duration" :duration="run.duration" />
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
