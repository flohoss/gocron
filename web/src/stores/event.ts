import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import type { JobsView } from '../client/types.gen';
import { getJob, getJobs } from '../client/sdk.gen';

export type EventInfo = {
  idle: boolean;
  data: JobsView;
};

export const useEventStore = defineStore('event', () => {
  const idle = ref<boolean>(false);
  const currentJobId = ref<string | null>(null);
  const currentJob = computed(() => state.jobs.get(currentJobId.value + ''));
  const state = reactive<{ loading: boolean; error: string | null; jobs: Map<string, JobsView> }>({
    loading: false,
    error: null,
    jobs: new Map<string, JobsView>(),
  });
  const fetchSuccess = computed(() => state.error === null && state.loading === false && state.jobs !== null);

  function parseEventInfo(info: string | null): void {
    if (!info) return;
    const parsed: EventInfo = JSON.parse(info);
    idle.value = parsed.idle;
    if (parsed.data) {
      state.jobs.set(parsed.data.id, parsed.data);
    }
  }

  async function fetchJobs() {
    currentJobId.value = null;
    state.error = null;
    state.loading = true;

    try {
      const result = await getJobs();
      result.data?.map((job: JobsView) => state.jobs.set(job.id, job));
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  async function fetchJob(id: string | string[]) {
    currentJobId.value = id + '';
    state.error = null;
    state.loading = true;

    try {
      const result = await getJob({ path: { name: id + '' } });
      state.jobs.set(result.data!.id, result.data!);
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  return { idle, currentJobId, parseEventInfo, state, fetchJobs, fetchSuccess, fetchJob, currentJob };
});
