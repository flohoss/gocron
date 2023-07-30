<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import { SystemService, type database_SystemLog } from '@/openapi';
import { useEventSource } from '@vueuse/core';
import { onBeforeUnmount, ref, watch } from 'vue';
import TerminalLog from '@/components/ui/TerminalLog.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import SystemOverview from '@/components/system/SystemOverview.vue';

const logs = ref<database_SystemLog[]>([]);

const init = () => {
  SystemService.getSystemLogs()
    .then((res) => {
      logs.value = res;
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
      <SystemOverview />
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto">
        <TerminalLog :logs="logs" />
      </div>
    </PageContent>
  </div>
</template>
