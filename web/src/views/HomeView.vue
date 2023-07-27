<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import { SystemService, type database_SystemLog, type system_Data } from '@/openapi';
import { useEventSource } from '@vueuse/core';
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import TerminalLog from '@/components/ui/TerminalLog.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemStat from '@/components/system/SystemStat.vue';
import { emptySystem } from '@/helper/system';

const logs = ref<database_SystemLog[]>([]);
const system = ref<system_Data>(emptySystem);

const init = () => {
  SystemService.getSystemLogs()
    .then((res) => {
      logs.value = res;
    })
    .catch((err) => console.log(err));
  SystemService.getSystem()
    .then((res) => {
      system.value = res;
    })
    .catch((err) => console.log(err));
};
init();

const { data, close } = useEventSource('/api/sse?stream=system_logs');
watch(data, (value) => {
  const parsed: database_SystemLog = value && JSON.parse(value);
  logs.value.unshift(parsed);
});
onBeforeUnmount(() => close());

const resticAmount = computed(() =>
  system.value.job_stats.total_runs !== 0 ? (system.value.job_stats.restic_runs / system.value.job_stats.total_runs) * 100 : 0
);
const checkAmount = computed(() =>
  system.value.job_stats.total_runs !== 0 ? (system.value.job_stats.check_runs / system.value.job_stats.total_runs) * 100 + resticAmount.value : 0
);
const pruneAmount = computed(() =>
  system.value.job_stats.total_runs !== 0 ? (system.value.job_stats.prune_runs / system.value.job_stats.total_runs) * 100 + checkAmount.value : 0
);
const customAmount = computed(() =>
  system.value.job_stats.total_runs !== 0 ? (system.value.job_stats.custom_runs / system.value.job_stats.total_runs) * 100 + pruneAmount.value : 0
);
const runStats = computed<{ percent: number; value: number; desc: string }[]>(() => [
  { percent: customAmount.value, value: system.value.job_stats.custom_runs, desc: 'Custom' },
  { percent: pruneAmount.value, value: system.value.job_stats.prune_runs, desc: 'Prune' },
  { percent: checkAmount.value, value: system.value.job_stats.check_runs, desc: 'Check' },
  { percent: resticAmount.value, value: system.value.job_stats.restic_runs, desc: 'Restic' },
]);

const infoAmount = computed(() => (system.value.job_stats.total_logs !== 0 ? (system.value.job_stats.info_logs / system.value.job_stats.total_logs) * 100 : 0));
const warningAmount = computed(() =>
  system.value.job_stats.total_logs !== 0 ? (system.value.job_stats.warning_logs / system.value.job_stats.total_logs) * 100 + infoAmount.value : 0
);
const errorAmount = computed(() =>
  system.value.job_stats.total_logs !== 0 ? (system.value.job_stats.error_logs / system.value.job_stats.total_logs) * 100 + warningAmount.value : 0
);
const logsStats = computed<{ percent: number; value: number; desc: string }[]>(() => [
  { percent: errorAmount.value, value: system.value.job_stats.error_logs, desc: 'Error' },
  { percent: warningAmount.value, value: system.value.job_stats.warning_logs, desc: 'Warning' },
  { percent: infoAmount.value, value: system.value.job_stats.info_logs, desc: 'Info' },
]);
</script>

<template>
  <div>
    <PageHeader>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-10">
        <SystemStat title="Runs" :value="system.job_stats.total_runs" :data="runStats" />
        <SystemStat title="Logs" :value="system.job_stats.total_logs" :data="logsStats" />
      </div>
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto">
        <TerminalLog :logs="logs" />
      </div>
    </PageContent>
  </div>
</template>
