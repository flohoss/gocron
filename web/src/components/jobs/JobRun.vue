<script setup lang="ts">
import type { database_Run } from '@/openapi';
import { computed } from 'vue';
import moment from 'moment';
import { severityIcons, severityColor } from '@/helper/severity';
import TerminalLog from '@/components/ui/TerminalLog.vue';

const props = defineProps<{ run: database_Run; checked: boolean }>();

const startDate = computed(() => moment(props.run.start_time).format('L'));
const startTime = computed(() => moment(props.run.start_time).format('LTS'));
const startDateTime = computed(
  () =>
    startDate.value +
    ' ' +
    startTime.value +
    (props.run.end_time ? ' (' + moment.duration(props.run.end_time - (props.run.start_time || 0)).humanize() + ')' : '')
);

const severity = computed(() => {
  let severity = 0;
  if (props.run.id === 0) {
    return -1;
  }
  if (!props.run.end_time) {
    return severity;
  }
  for (let log of props.run.logs) {
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
  <div>
    <div class="flex items-center gap-2" :class="color">
      <div class="underline underline-offset-4">{{ startDateTime }}</div>
      <span v-if="severity !== 0"><i :class="' fa-solid fa-' + icon"></i></span>
      <span v-else class="loading loading-spinner"></span>
    </div>
    <TerminalLog v-if="run.logs" :logs="run.logs" />
  </div>
</template>
