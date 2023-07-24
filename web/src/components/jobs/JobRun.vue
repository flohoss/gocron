<script setup lang="ts">
import type { database_Run } from '@/openapi';
import { computed } from 'vue';
import moment from 'moment';
import { severityHTML, severityColor } from '@/helper/severity';
import { getIcon } from '@/helper/logType';

const props = defineProps<{ run: database_Run; checked: boolean }>();

const startDate = computed(() => moment(props.run.start_time).format('L'));
const duration = computed(() => moment.duration((props.run.end_time || 0) - (props.run.start_time || 0)).humanize());

const status = (run: database_Run | undefined) => {
  if (!run || !run.end_time) {
    return `<div class="text-info"><i class="fa-solid fa-play"></div>`;
  }
  let severity = 0;
  if (run.logs) {
    for (let log of run.logs) {
      if (log.log_severity_id && severity < log.log_severity_id) {
        severity = log.log_severity_id;
      }
    }
  }
  return severityHTML(severity);
};

const formatDate = (ts: number | undefined) => moment(ts).format('LTS');
</script>

<template>
  <div class="join-item collapse">
    <input type="radio" name="accordion" :checked="checked" />
    <div class="collapse-title flex items-center justify-between px-5">
      <div class="flex flex-col">
        <div>{{ startDate }}</div>
        <div class="text-sm opacity-50">{{ duration }}</div>
      </div>
      <div v-html="status(run)"></div>
    </div>
    <div class="collapse-content font-mono text-sm px-5">
      <div v-for="log of run.logs" :key="log.id" class="flex items-start gap-2" :class="severityColor(log.log_severity_id)">
        <div v-html="getIcon(log.log_type_id)"></div>
        <div class="whitespace-nowrap">{{ formatDate(log.created_at) }}</div>
        <div class="whitespace-pre-line">{{ log.message }}</div>
      </div>
    </div>
  </div>
</template>
