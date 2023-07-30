<script setup lang="ts">
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { useJobStore } from '@/stores/jobs';
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { CommandsService, database_Job, JobsService, type database_Run, database_LogSeverity, type database_Log } from '@/openapi';
import JobRun from '@/components/jobs/JobRun.vue';
import CustomCommand from '@/components/jobs/CustomCommand.vue';
import JobHeader from '@/components/jobs/JobHeader.vue';
import { useEventSource } from '@vueuse/core';
import { EventType, type SSEvent } from '@/types';

const store = useJobStore();
const route = useRoute();
const job = computed<database_Job>(() => store.getJob(route.params.id));
const cmdModal = ref();

const startCommand = async (cmd: string, custom?: string) => {
  CommandsService.postCommands({
    local_directory: '',
    password_file_path: '',
    restic_remote: '',
    command: cmd,
    job_id: job.value.id,
    custom_command: custom ? custom : '',
  }).catch((err) => console.log(err));
};

const runs = ref<database_Run[]>([]);

const getRuns = async () => {
  job.value.id && (runs.value = await JobsService.getJobsRuns(job.value.id));
};
getRuns();
watch(job, () => getRuns());

const { data, close } = useEventSource('/api/sse?stream=jobs');
watch(data, (value) => {
  const parsed: SSEvent = JSON.parse(value + '');
  switch (parsed.event_type) {
    case EventType.EventCreateRun: {
      runs.value.unshift(parsed.content as database_Run);
      job.value.status = database_LogSeverity.LogRunning;
      break;
    }
    case EventType.EventCreateLog: {
      const log = parsed.content as database_Log;
      const run = runs.value.find((r) => r.id == log.run_id);
      if (run) {
        if (!run.logs) run.logs = [];
        run.logs.push(log);
      }
      break;
    }
    case EventType.EventUpdateRun: {
      const run = parsed.content as database_Run;
      const currentRun = runs.value.find((r) => r.id == run.id);
      if (currentRun && currentRun.logs) {
        currentRun.end_time = run.end_time;
        let severity = database_LogSeverity.LogNone;
        for (const log of currentRun.logs) {
          if (log.log_severity > severity) {
            severity = log.log_severity;
          }
        }
        job.value.status = severity;
      }
      break;
    }
  }
});
onBeforeUnmount(() => close());
</script>

<template>
  <div>
    <CustomCommand @gotRef="(el) => (cmdModal = el)" @start="(c, v) => startCommand(c, v)" />
    <PageHeader>
      <JobHeader :job="job" @showModal="cmdModal.showModal()" @start="(c) => startCommand(c)" />
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto">
        <JobRun v-for="(run, i) of runs" :key="run.id" :run="run" :checked="i === 0" />
      </div>
    </PageContent>
  </div>
</template>
