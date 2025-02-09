<script setup lang="ts">
import { computed } from 'vue';
import { BackendURL } from '../main';
import type { jobs_JobsView } from '../openapi';
import humanizeDuration from 'humanize-duration';

const props = defineProps<{ job: jobs_JobsView }>();
const url = computed<string>(() => BackendURL + props.job.id);

const duration = (duration: number) => humanizeDuration(duration, { round: true, largest: 1 });
</script>

<template>
  <a class="flex justify-between items-center group last:mb-8 lg:last:mb-0" :href="url">
    <div class="pl-4 truncate">
      <div class="group-hover:text-primary hover-animation text-2xl font-medium truncate">{{ job.name }}</div>
      <div class="text-secondary text-sm truncate">{{ job.cron }}</div>
    </div>
    <div class="text-sm">
      <ul class="steps hidden lg:flex">
        <li v-for="run in job.runs" :key="run.id">
          <span v-if="run.duration.Valid">{{ duration(run.duration.Int64) }}</span>
        </li>
      </ul>
    </div>
  </a>
</template>
