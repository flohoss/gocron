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
  const checked = ref<string[]>([]);

  const fetchSuccess = computed(() => error.value === null && loading.value === false);
  const disabled = computed(() => loading.value || !idle.value);
  const currentJob = computed(() => {
    if (route.params.id) {
      const job = jobs.value.get(route.params.id + '');
      if (job) return job;
    }
    return null;
  });
  const jobsUnchecked = computed(() => {
    return jobs.value.size - checked.value.length;
  });

  function setJobs(newJobs: JobView[]) {
    jobs.value.clear();
    newJobs.forEach((job) => jobs.value.set(job.name, job));
    checked.value = newJobs.reduce<string[]>((acc, job) => {
      if (!job.disabled) acc.push(job.name);
      return acc;
    }, []);
  }

  function parseEventInfo(info: string | null): void {
    if (!info) return;

    const parsed: EventInfo = JSON.parse(info);
    idle.value = parsed.idle;

    if (parsed.jobs) {
      setJobs(parsed.jobs);
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
      setJobs(result.data);
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

  return {
    idle,
    jobs,
    loading,
    disabled,
    fetchSuccess,
    currentJob,
    checked,
    jobsUnchecked,
    parseEventInfo,
    fetchJobs,
    fetchJob,
  };
});
