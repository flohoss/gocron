<script setup lang="ts">
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { useJobStore } from '@/stores/jobs';
import { computed, onBeforeUnmount, watch, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { CommandsService, database_Job, type database_Run } from '@/openapi';
import JobRun from '@/components/jobs/JobRun.vue';
import { useEventSource } from '@vueuse/core';
import CustomCommand from '@/components/jobs/CustomCommand.vue';
import JobHeader from '@/components/jobs/JobHeader.vue';

const store = useJobStore();
const route = useRoute();
const job = computed<database_Job>(() => store.getJob(route.params.id));
const cmdModal = ref();

const startCommand = async (cmd: string, custom?: string) => {
  CommandsService.postCommands({
    command: cmd,
    job_id: job.value.id,
    custom_command: custom,
  }).catch((err) => console.log(err));
};

const runsEvent = useEventSource('/api/sse?stream=runs');
watch(runsEvent.data, (value) => {
  const parsed: database_Run = value && JSON.parse(value);
  store.updateOrCreateRun(parsed);
});
onBeforeUnmount(() => runsEvent.close());
const logsEvent = useEventSource('/api/sse?stream=logs');
watch(logsEvent.data, (value) => {
  const parsed: database_Run = value && JSON.parse(value);
  store.updateOrCreateRun(parsed);
});
onBeforeUnmount(() => logsEvent.close());
</script>

<template>
  <div>
    <CustomCommand @gotRef="(el) => (cmdModal = el)" @start="(c, v) => startCommand(c, v)" />
    <PageHeader>
      <JobHeader :job="job" @showModal="cmdModal.showModal()" @start="(c) => startCommand(c)" />
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto">
        <JobRun v-for="(run, i) of job.runs" :key="run.id" :run="run" :checked="i === 0" />
      </div>
    </PageContent>
  </div>
</template>
