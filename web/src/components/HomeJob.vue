<script setup lang="ts">
import { computed } from 'vue';
import { RouterLink } from 'vue-router';
import { useWindowSize } from '@vueuse/core';
import type { JobView, RunView } from '../client/types.gen';
import JobStep from './JobStep.vue';

const { job } = defineProps<{ job: JobView }>();
const url = computed<string>(() => '/jobs/' + job.name);

const { width } = useWindowSize();
const isMobile = computed(() => width.value < 1024);
const runs = computed<RunView[]>(() => {
  const runsArray = job.runs ?? [];
  const amount = runsArray.length;

  if (amount === 0) return [];

  if (isMobile.value) {
    const lastRun = runsArray[amount - 1];
    return lastRun ? [lastRun] : [];
  } else {
    return runsArray.slice(-Math.min(3, amount));
  }
});

const cron = computed(() => (job.disable_cron ? '-' : job.cron));
</script>

<template>
  <RouterLink
    data-test-id="job-link"
    :data-test-name="job.name"
    class="flex justify-between items-center group last:mb-8 lg:last:mb-0 hover:cursor-pointer"
    :to="url"
  >
    <div class="pl-4 truncate">
      <div class="group-hover:text-primary hover-animation text-2xl font-medium truncate">{{ job.name }}</div>
      <div class="text-secondary text-sm truncate">{{ cron }}</div>
    </div>
    <div class="text-sm">
      <ul class="steps" v-if="runs">
        <JobStep v-for="run in runs" :key="run.id" :run="run" />
      </ul>
    </div>
  </RouterLink>
</template>
