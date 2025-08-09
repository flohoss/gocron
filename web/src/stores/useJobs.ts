import { computed, ref } from 'vue';
import { createGlobalState } from '@vueuse/core';
import type { JobView, RunView } from '../client/types.gen';
import { getJobs, getRuns } from '../client/sdk.gen';
import { useRoute } from 'vue-router';

export type EventInfo = {
  idle: boolean;
  run?: RunView;
  jobs?: JobView[];
};

export const useJobs = createGlobalState(() => {
  const route = useRoute();

  const idle = ref(false);
  const jobs = ref(new Map<string, JobView>());
  const error = ref<string | null>(null);
  const loading = ref(false);

  const fetchSuccess = computed(() => error.value === null && loading.value === false);
  const disabled = computed(() => loading.value || !idle.value);
  const currentJob = computed(() => {
    if (route.params.id) {
      const job = jobs.value.get(route.params.id + '');
      if (job) return job;
    }
    return null;
  });

  function parseEventInfo(info: string | null): void {
    if (!info) return;

    const parsed: EventInfo = JSON.parse(info);
    idle.value = parsed.idle;

    if (parsed.jobs) {
      jobs.value.clear();
      parsed.jobs.forEach((job) => jobs.value.set(job.name, job));
    }

    if (!parsed.run) return;

    const jobName = parsed.run.job_name;
    const existingJobView = jobs.value.get(jobName);
    if (!existingJobView) return;

    const existingRuns = existingJobView.runs ?? [];

    const runIndex = existingRuns.findIndex((run) => parsed.run && run.id === parsed.run.id);

    let updatedRuns;
    if (runIndex !== -1) {
      updatedRuns = [...existingRuns];
      updatedRuns[runIndex] = parsed.run;
    } else {
      updatedRuns = [...existingRuns, parsed.run];
    }

    jobs.value.set(jobName, {
      ...existingJobView,
      runs: updatedRuns,
    });
  }

  async function fetchJobs() {
    error.value = null;
    loading.value = true;

    try {
      const result = await getJobs();
      if (!result.data) return;
      result.data.forEach((job) => jobs.value.set(job.name, job));
    } catch (err: any) {
      error.value = err.toString();
    } finally {
      loading.value = false;
    }
  }

  async function fetchJob() {
    error.value = null;
    loading.value = true;

    if (!currentJob.value) return;
    const jobName = currentJob.value.name;

    const existingJobView = jobs.value.get(jobName);
    if (!existingJobView) return;

    try {
      const result = await getRuns({ path: { job_name: jobName } });
      if (!result.data) return;

      if (existingJobView) {
        jobs.value.set(jobName, {
          ...existingJobView,
          runs: result.data,
        });
      }
    } catch (err: any) {
      error.value = err.toString();
    } finally {
      loading.value = false;
    }
  }

  return { idle, jobs, loading, disabled, fetchSuccess, currentJob, parseEventInfo, fetchJobs, fetchJob };
});
