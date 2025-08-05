import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import type { JobView, RunView } from '../client/types.gen';
import { getJobs, getRuns } from '../client/sdk.gen';

export type EventInfo = {
  idle: boolean;
  run: RunView;
};

export const useEventStore = defineStore('event', () => {
  const idle = ref<boolean>(false);
  const currentJobName = ref<string | null>(null);
  const currentJob = computed(() => state.jobs.get(currentJobName.value + ''));
  const state = reactive<{ loading: boolean; error: string | null; jobs: Map<string, JobView> }>({
    loading: false,
    error: null,
    jobs: new Map<string, JobView>(),
  });
  const fetchSuccess = computed(() => state.error === null && state.loading === false && state.jobs !== null);

  function parseEventInfo(info: string | null): void {
    if (!info) return;
    const parsed: EventInfo = JSON.parse(info);
    idle.value = parsed.idle;

    if (!parsed.run) return;
    const existingJobView = state.jobs.get(parsed.run.job_name);
    console.log(existingJobView);
    if (existingJobView) {
      // Append the run to existing job's runs
      const updatedRuns = existingJobView.runs ? [...existingJobView.runs, parsed.run] : [parsed.run];
      state.jobs.set(parsed.run.job_name, {
        ...existingJobView,
        runs: updatedRuns,
      });
    }
  }

  async function fetchJobs() {
    currentJobName.value = null;
    state.error = null;
    state.loading = true;

    try {
      const result = await getJobs();
      if (!result.data) return;
      result.data.forEach((job) => state.jobs.set(job.name, job));
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  async function fetchJob(jobName: string | string[]) {
    currentJobName.value = jobName + '';
    state.error = null;
    state.loading = true;

    try {
      const result = await getRuns({ path: { job_name: currentJobName.value } });
      if (!result.data) return;
      const existingJobView = state.jobs.get(currentJobName.value);
      if (existingJobView) {
        state.jobs.set(currentJobName.value, {
          ...existingJobView,
          runs: result.data,
        });
      }
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  return { idle, currentJobName, parseEventInfo, state, fetchJobs, fetchSuccess, fetchJob, currentJob };
});
