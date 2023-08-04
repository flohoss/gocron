<script setup lang="ts">
import type { database_Log, database_SystemLog } from '@/openapi';
import { severityColor } from '@/helper/severity';
import moment from 'moment';
import { getIcon } from '@/helper/logType';

defineProps<{ logs: database_Log[] | database_SystemLog[] }>();
const formatDate = (ts: number) => moment(ts).format('LTS');
</script>

<template>
  <div class="font-mono text-xs">
    <div v-for="log of logs" :key="log.id" class="flex items-start gap-2" :class="severityColor(log.log_severity)">
      <div v-if="log.created_at" class="whitespace-nowrap">{{ formatDate(log.created_at) }}</div>
      <div v-if="log.hasOwnProperty('log_type')" v-html="getIcon((log as database_Log).log_type)"></div>
      <div class="whitespace-pre">{{ log.message }}</div>
    </div>
  </div>
</template>
