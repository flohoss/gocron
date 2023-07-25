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

const status = (run: database_Run | undefined) => {
  if (!run || !run.end_time) {
    return severityIcons(0);
  }
  let severity = 0;
  if (run.logs) {
    for (let log of run.logs) {
      if (log.log_severity_id && severity < log.log_severity_id) {
        severity = log.log_severity_id;
      }
    }
  }
  return severityIcons(severity);
};

const color = (run: database_Run | undefined) => {
  if (!run || !run.end_time) {
    return '';
  }
  let severity = 0;
  if (run.logs) {
    for (let log of run.logs) {
      if (log.log_severity_id && severity < log.log_severity_id) {
        severity = log.log_severity_id;
      }
    }
  }
  return severityColor(severity);
};
</script>

<template>
  <div>
    <div class="flex items-center gap-2" :class="color(run)">
      <div class="underline underline-offset-4">{{ startDateTime }}</div>
      <div class="flex items-center" v-html="status(run)"></div>
    </div>
    <TerminalLog v-if="run.logs" :logs="run.logs" />
  </div>
</template>
