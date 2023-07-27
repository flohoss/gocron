<script setup lang="ts">
import type { database_Job } from '@/openapi';
import NavLink from '../ui/NavLink.vue';
import { computed } from 'vue';
import { severityIcons } from '@/helper/severity';

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

const severity = computed(() => {
  if (props.job.runs?.length === 0) return '';
  if (!props.job.runs || !props.job.runs[0].end_time) {
    return severityIcons(0);
  }
  let severity = 0;
  if (props.job.runs && props.job.runs[0].logs) {
    for (let log of props.job.runs[0].logs) {
      if (severity < log.log_severity) {
        severity = log.log_severity;
      }
    }
  }
  return severityIcons(severity);
});
</script>

<template>
  <NavLink
    @click="emit('click')"
    :name="job.description"
    :extra="job.local_directory"
    :link="{ name: 'jobs', params: { id: job.id } }"
    :active="active"
    :icon="jobIcon()"
  >
    <span v-html="severity"></span>
  </NavLink>
</template>
