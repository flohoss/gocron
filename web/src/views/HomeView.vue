<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import { SystemService, type database_SystemLog, type system_Data } from '@/openapi';
import { useEventSource, useWindowSize } from '@vueuse/core';
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import TerminalLog from '@/components/ui/TerminalLog.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemStat from '@/components/system/SystemStat.vue';
import { emptySystem } from '@/helper/system';

const logs = ref<database_SystemLog[]>([]);
const system = ref<system_Data>(emptySystem);
const { width } = useWindowSize();

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

const runChartData = computed(() => {
  return {
    labels: ['restic', 'check', 'prune', 'custom'],
    datasets: [
      {
        data: [system.value.job_stats.restic_runs, system.value.job_stats.check_runs, system.value.job_stats.prune_runs, system.value.job_stats.custom_runs],
        borderWidth: 0,
      },
    ],
  };
});

const logChartData = computed(() => {
  return {
    labels: ['info', 'error', 'warning'],
    datasets: [
      {
        data: [system.value.job_stats.info_logs, system.value.job_stats.error_logs, system.value.job_stats.warning_logs],
        borderWidth: 0,
      },
    ],
  };
});
</script>

<template>
  <div>
    <PageHeader>
      <div class="grid w-full grid-cols-2">
        <div class="flex">
          <SystemStat title="Runs" :value="system.job_stats.total_runs" />
          <SystemStat v-if="width >= 768" title="Restic" class="opacity-50" :value="system.job_stats.restic_runs" />
          <SystemStat v-if="width >= 768" title="Check" class="opacity-50" :value="system.job_stats.check_runs" />
          <SystemStat v-if="width >= 1280" title="Prune" class="opacity-50" :value="system.job_stats.prune_runs" />
          <SystemStat v-if="width >= 1280" title="Custom" class="opacity-50" :value="system.job_stats.custom_runs" />
        </div>
        <div class="flex">
          <SystemStat title="Logs" :value="system.job_stats.total_logs" />
          <SystemStat v-if="width >= 768" title="Info" :hidden="true" class="opacity-50" :value="system.job_stats.info_logs" />
          <SystemStat v-if="width >= 1280" title="Warning" class="opacity-50" :value="system.job_stats.warning_logs" />
          <SystemStat v-if="width >= 1280" title="Error" class="opacity-50" :value="system.job_stats.error_logs" />
        </div>
      </div>
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto">
        <TerminalLog :logs="logs" />
      </div>
    </PageContent>
  </div>
</template>
