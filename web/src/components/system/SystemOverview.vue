<script setup lang="ts">
import { SystemService, type database_JobStats, type system_Versions } from '@/openapi';
import { computed, ref } from 'vue';
import SystemStat from './SystemStat.vue';
import BadgeList from '../ui/BadgeList.vue';
import { emptyJobStats } from '@/types';

defineProps<{ versions: system_Versions; badges: any }>();

const system = ref<database_JobStats>(emptyJobStats);

const init = () => {
  SystemService.getSystemStats()
    .then((res) => {
      system.value = res;
    })
    .catch((err) => console.log(err));
};
init();

const generalAmount = computed(() => (system.value.total_runs !== 0 ? (system.value.general_runs / system.value.total_runs) * 100 : 0));
const resticAmount = computed(() => (system.value.total_runs !== 0 ? (system.value.restic_runs / system.value.total_runs) * 100 + generalAmount.value : 0));
const checkAmount = computed(() => (system.value.total_runs !== 0 ? (system.value.check_runs / system.value.total_runs) * 100 + resticAmount.value : 0));
const pruneAmount = computed(() => (system.value.total_runs !== 0 ? (system.value.prune_runs / system.value.total_runs) * 100 + checkAmount.value : 0));
const customAmount = computed(() => (system.value.total_runs !== 0 ? (system.value.custom_runs / system.value.total_runs) * 100 + pruneAmount.value : 0));
const runStats = computed<{ percent: number; value: number; desc: string }[]>(() => [
  { percent: customAmount.value, value: system.value.custom_runs, desc: 'Custom' },
  { percent: pruneAmount.value, value: system.value.prune_runs, desc: 'Prune' },
  { percent: checkAmount.value, value: system.value.check_runs, desc: 'Check' },
  { percent: resticAmount.value, value: system.value.restic_runs, desc: 'Restic' },
  { percent: generalAmount.value, value: system.value.general_runs, desc: 'General' },
]);

const infoAmount = computed(() => (system.value.total_logs !== 0 ? (system.value.info_logs / system.value.total_logs) * 100 : 0));
const warningAmount = computed(() => (system.value.total_logs !== 0 ? (system.value.warning_logs / system.value.total_logs) * 100 + infoAmount.value : 0));
const errorAmount = computed(() => (system.value.total_logs !== 0 ? (system.value.error_logs / system.value.total_logs) * 100 + warningAmount.value : 0));
const logsStats = computed<{ percent: number; value: number; desc: string }[]>(() => [
  { percent: errorAmount.value, value: system.value.error_logs, desc: 'Error' },
  { percent: warningAmount.value, value: system.value.warning_logs, desc: 'Warning' },
  { percent: infoAmount.value, value: system.value.info_logs, desc: 'Info' },
]);
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-10">
    <SystemStat title="Runs" :value="system.total_runs" :data="runStats" />
    <SystemStat title="Logs" :value="system.total_logs" :data="logsStats" />
  </div>
  <div class="flex gap-2 mt-5 flex-wrap select-none">
    <BadgeList :badges="versions" />
    <BadgeList :badges="badges" />
  </div>
</template>
