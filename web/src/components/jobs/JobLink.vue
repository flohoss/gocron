<script setup lang="ts">
import type { database_Job } from '@/openapi';
import NavLink from '../ui/NavLink.vue';
import { computed } from 'vue';
import { severityColor, severityIcons } from '@/helper/severity';

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
  let severity = 0;
  if (props.job.runs.length === 0) {
    return -1;
  }
  if (!props.job.runs[0].end_time) {
    return severity;
  }
  for (let log of props.job.runs[0].logs) {
    if (severity < log.log_severity) {
      severity = log.log_severity;
    }
  }
  return severity;
});

const color = computed(() => severityColor(severity.value));
const icon = computed(() => severityIcons(severity.value));
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
    <span v-if="severity !== 0"><i :class="color + ' fa-solid fa-' + icon"></i></span>
    <span v-else class="loading loading-spinner"></span>
  </NavLink>
</template>
