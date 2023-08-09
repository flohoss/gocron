<script setup lang="ts">
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { useJobStore } from '@/stores/jobs';
import { computed, inject, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { CommandsService, database_Job, JobsService, type database_Run, type database_Log, database_LogSeverity } from '@/openapi';
import JobRun from '@/components/jobs/JobRun.vue';
import CustomCommand from '@/components/jobs/CustomCommand.vue';
import JobHeader from '@/components/jobs/JobHeader.vue';
import { EventType, sseKey } from '@/types';
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue';

const store = useJobStore();
const route = useRoute();
const job = computed<database_Job>(() => store.getJob(route.params.id));
const cmdModal = ref();
const loading = ref(true);
const error = ref('');

const startCommand = async (cmd: string, custom?: string) => {
  error.value = '';
  CommandsService.postCommands({
    local_directory: '',
    password_file_path: '',
    restic_remote: '',
    command: cmd,
    job_id: job.value.id,
    custom_command: custom ? custom : '',
  })
    .then(() => cmdModal.value.close())
    .catch((err) => (error.value = err));
};

const runs = ref<database_Run[]>([]);

const getRuns = async () => {
  if (job.value.id) {
    JobsService.getJobsRuns(job.value.id)
      .then((res) => {
        runs.value = res;
        loading.value = false;
      })
      .catch((err) => console.log(err.body));
  }
};
getRuns();
watch(job, () => getRuns());

const parsed = inject(sseKey, ref());
watch(parsed, (value) => {
  if (!value) return;
  console.log(value);
  switch (value.event_type) {
    case EventType.EventCreateRun: {
      runs.value.unshift(value.content as database_Run);
      break;
    }
    case EventType.EventCreateLog: {
      const log = value.content as database_Log;
      const run = runs.value.find((r) => r.id == log.run_id);
      if (run) {
        if (!run.logs) run.logs = [];
        run.logs.push(log);
      }
      break;
    }
    case EventType.EventUpdateRun: {
      const run = value.content as database_Run;
      const currentRun = runs.value.find((r) => r.id == run.id);
      if (currentRun && currentRun.logs) {
        currentRun.end_time = run.end_time;
        job.value.status = run.status;
      }
      break;
    }
  }
});

const clearRuns = () => {
  runs.value = [];
  job.value.status = database_LogSeverity.LogInfo;
};
</script>

<template>
  <div>
    <CustomCommand :error="error" @gotRef="(el) => (cmdModal = el)" @start="(c, v) => startCommand(c, v)" />
    <PageHeader>
      <JobHeader :job="job" @showModal="cmdModal.showModal()" @start="(c) => startCommand(c)" @clearRuns="clearRuns" />
    </PageHeader>
    <PageContent>
      <Transition>
        <LoadingSpinner v-if="loading" />
        <div v-else class="grid grid-cols-1 gap-5 overflow-x-auto">
          <JobRun v-for="(run, i) of runs" :key="run.id" :run="run" :checked="i === 0" />
        </div>
      </Transition>
    </PageContent>
  </div>
</template>
