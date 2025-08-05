import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import type { JobView } from '../client/types.gen';
import { getJobs, getRuns } from '../client/sdk.gen';

export type EventInfo = {
  idle: boolean;
};

export const useEventStore = defineStore('event', () => {
  const idle = ref<boolean>(false);
  const currentJobId = ref<string | null>(null);
  const currentJob = computed(() => state.jobs.get(currentJobId.value + ''));
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
  }

  async function fetchJobs() {
    currentJobId.value = null;
    state.error = null;
    state.loading = true;

    try {
      const result = await getJobs();
      result.data?.map((job: JobView) => state.jobs.set(job.name, job));
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  async function fetchJob(jobName: string | string[]) {
    currentJobId.value = jobName + '';
    state.error = null;
    state.loading = true;

    try {
      const result = await getRuns({ path: { job_name: currentJobId.value } });
      const existingJobView = state.jobs.get(currentJobId.value);
      if (existingJobView) {
        state.jobs.set(currentJobId.value, {
          ...existingJobView,
          runs: result.data ? result.data : [],
        });
      }
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  return { idle, currentJobId, parseEventInfo, state, fetchJobs, fetchSuccess, fetchJob, currentJob };
});
