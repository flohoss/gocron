<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import { SystemService, type database_SystemLog, type system_Data } from '@/openapi';
import { useEventSource } from '@vueuse/core';
import { onBeforeUnmount, ref, watch } from 'vue';
import TerminalLog from '@/components/ui/TerminalLog.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemStat from '@/components/system/SystemStat.vue';
import { emptySys } from '@/helper/system';

const logs = ref<database_SystemLog[]>([]);
const system = ref<system_Data>(emptySys);

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
</script>

<template>
  <div>
    <PageHeader>
      <div class="w-full grid grid-cols-1 gap-5">
        <div class="stats bg-transparent">
          <SystemStat title="Total Runs" :value="system.job_stats.total_runs" />
          <SystemStat title="Restic Runs" :value="system.job_stats.restic_runs" />
          <SystemStat title="Check Runs" :value="system.job_stats.check_runs" />
          <SystemStat title="Prune Runs" :value="system.job_stats.prune_runs" />
          <SystemStat title="Custom Runs" :value="system.job_stats.custom_runs" />
        </div>
        <div class="stats bg-transparent">
          <SystemStat title="Logs" :value="system.job_stats.total_logs" class="text-info" />
          <SystemStat title="Errors" :value="system.job_stats.error_logs" class="text-error" />
          <SystemStat title="Warnings" :value="system.job_stats.warning_logs" class="text-warning" />
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
