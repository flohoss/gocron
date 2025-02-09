<script setup lang="ts">
import { computed } from 'vue';
import { type jobs_JobsView } from '../openapi';
import humanizeDuration from 'humanize-duration';
import { RouterLink } from 'vue-router';

const props = defineProps<{ job: jobs_JobsView }>();
const url = computed<string>(() => '/jobs/' + props.job.id);

const shortEnglishHumanizer = humanizeDuration.humanizer({
  language: 'shortEn',
  languages: {
    shortEn: {
      y: () => 'y',
      mo: () => 'mo',
      w: () => 'w',
      d: () => 'd',
      h: () => 'h',
      m: () => 'm',
      s: () => 's',
      ms: () => 'ms',
    },
  },
});
const duration = (duration: number) => shortEnglishHumanizer(duration, { round: true, largest: 1 });

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
      <ul class="steps hidden lg:flex">
        <li v-for="run in job.runs" :key="run.id" :data-content="getStepIcon(run.status_id)" class="step" :class="getStepColor(run.status_id)">
          <span v-if="run.duration.Valid">{{ duration(run.duration.Int64) }}</span>
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
