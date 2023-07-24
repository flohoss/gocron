<script setup lang="ts">
import type { database_Job, database_Run } from '@/openapi';
import NavLink from '../ui/NavLink.vue';
import { severityHTML } from '@/helper/severity';

const props = defineProps<{
  job: database_Job;
  active: boolean;
}>();
const emit = defineEmits(['click']);

const jobIcon = () => {
  if (props.job.svg_icon) {
    return props.job.svg_icon;
  }
  return props.job.description.charAt(0);
};

const status = (runs: database_Run[] | undefined) => {
  if (!runs || runs.length === 0) {
    return `<div class="text-info"><i class="fa-solid fa-play"></div>`;
  }
  let severity = 0;
  if (runs[0].logs) {
    for (let log of runs[0].logs) {
      if (log.log_severity_id && severity < log.log_severity_id) {
        severity = log.log_severity_id;
      }
    }
  }
  return severityHTML(severity);
};
</script>

<template>
  <NavLink
    @click="emit('click')"
    :name="job.description"
    :extra="job.local_directory"
    :link="{ name: 'jobs', params: { id: job.id } }"
    :active="active"
    :icon="jobIcon()"
    :status="status(job.runs)"
    background="bg-primary"
    text="text-primary-content"
  />
</template>
