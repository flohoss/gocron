<script setup lang="ts">
import { watch } from 'vue';
import { useJobs } from '../stores/useJobs';
import CommandWindow from '../components/utils/CommandWindow.vue';

const { jobs, loading, fetchSuccess, currentJob, fetchJob } = useJobs();

watch(
  () => jobs.value.size,
  async (size) => {
    if (size > 0 && currentJob.value) {
      await fetchJob();
    }
  },
  { immediate: true }
);

enum Severity {
  Debug = 1,
  Info = 2,
  Warning = 3,
  Error = 4,
}

function getColor(severity: Severity): string {
  switch (severity) {
    case Severity.Info:
      return 'text-secondary';
    case Severity.Warning:
      return 'text-warning';
    case Severity.Error:
      return 'text-error';
    default:
      return 'text-base-content';
  }
}
</script>

<template>
  <CommandWindow :title="currentJob?.name">
    <div v-if="loading" class="p-4 flex justify-center items-center">
      <span class="text-secondary loading loading-dots loading-xl"></span>
    </div>
    <template v-else-if="fetchSuccess && currentJob" v-for="(run, i) in currentJob.runs" :key="i">
      <pre
        :id="`run-${i + 1}`"
        :class="getColor(Severity.Debug)"
      ><code>{{ run.start_time }}: Job <span class="text-primary font-bold">{{ currentJob.name }}</span> started</code></pre>

      <template v-for="log in run.logs" :key="log.id">
        <span :class="[getColor(log.severity_id), 'flex']">
          <pre><code>{{ log.created_at_time }}: </code></pre>
          <pre><code>{{ log.message }}</code></pre>
        </span>
      </template>
      <pre
        v-if="run.end_time !== '' && run.duration !== ''"
        :class="getColor(Severity.Debug)"
        class="mb-2 last:mb-0"
      ><code>{{ run.end_time }}: Job finished (took {{ run.duration }})</code></pre>
    </template>
  </CommandWindow>
</template>
