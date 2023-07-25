<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import { SystemService, type database_SystemLog } from '@/openapi';
import { useEventSource } from '@vueuse/core';
import { onBeforeUnmount, ref, watch } from 'vue';
import figlet from 'figlet';
// @ts-ignore
import standard from 'figlet/importable-fonts/Standard.js';
import TerminalLog from '@/components/ui/TerminalLog.vue';

const logo = ref<database_SystemLog[]>([]);
const logs = ref<database_SystemLog[]>([]);

figlet.parseFont('Standard', standard);
figlet.text('GoBackup', { font: 'Standard' }, (err, data) => {
  if (err) {
    console.log('Something went wrong...');
    console.dir(err);
    return;
  }
  logo.value.push({ message: data });
});

const init = async () => {
  const response = await SystemService.getSystemLogs();
  logs.value = response;
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
  <PageContent>
    <div class="grid grid-cols-1 gap-5 overflow-x-auto">
      <TerminalLog :logs="logo" />
      <TerminalLog :logs="logs" />
    </div>
  </PageContent>
</template>
