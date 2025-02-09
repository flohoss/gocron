<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { JobsService, type services_TemplateJob } from '../openapi';
import ShortDuration from '../components/ShortDuration.vue';

const route = useRoute();

const loading = ref(false);
const job = ref<services_TemplateJob | null>(null);
const error = ref(null);

watch(() => route.params.id, fetchData, { immediate: true });

async function fetchData(id: string | string[]) {
  error.value = job.value = null;
  loading.value = true;

  try {
    job.value = await JobsService.getJobs1(id + '');
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
}

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
  <div class="bg-base-300 w-full text-sm rounded-xl">
    <div class="flex justify-start items-center gap-2 padding bg-base-200 rounded-t-xl">
      <div class="console-btn bg-error text-error hover:text-error-content"></div>
      <div class="console-btn bg-warning text-warning hover:text-warning-content"></div>
      <div class="console-btn bg-success text-success hover:text-success-content"></div>
    </div>
    <div class="overflow-x-scroll padding" v-if="job">
      <template v-for="(run, i) in job.runs" :key="i">
        <pre
          :id="`run-${i + 1}`"
          :class="getColor(Severity.Debug)"
        ><code>{{ run.fmt_start_time }}: Job <span class="text-primary font-bold">{{ job.name }}</span> started</code></pre>

        <template v-for="log in run.logs" :key="log.id">
          <span :class="[getColor(log.severity_id), 'flex']">
            <pre><code>{{ log.created_at_time }}: </code></pre>
            <pre><code>{{ log.message }}</code></pre>
          </span>
        </template>
        <pre
          :class="getColor(Severity.Debug)"
          class="mb-2 last:mb-0"
        ><code>{{ run.fmt_end_time }}: Job finished <span v-if="run.duration">(took <ShortDuration :duration="run.duration" />)</span></code></pre>
      </template>
    </div>
  </div>
</template>
