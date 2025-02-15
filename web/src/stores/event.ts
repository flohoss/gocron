import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import { JobsService, type events_EventInfo, type jobs_JobsView } from '../openapi';

export const useEventStore = defineStore('event', () => {
  const idle = ref<boolean>(false);
  const currentJobId = ref<string | null>(null);
  const currentJob = computed(() => state.jobs.get(currentJobId.value + ''));
  const state = reactive<{ loading: boolean; error: string | null; jobs: Map<string, jobs_JobsView> }>({
    loading: false,
    error: null,
    jobs: new Map<string, jobs_JobsView>(),
  });
  const fetchSuccess = computed(() => state.error === null && state.loading === false && state.jobs !== null);

  function parseEventInfo(info: string | null): void {
    if (!info) return;
    const parsed: events_EventInfo = JSON.parse(info);
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
      const result = await JobsService.getJobs();
      result.map((job) => state.jobs.set(job.id, job));
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
      const result = await JobsService.getJobs1(id + '');
      state.jobs.set(result.id, result);
    } catch (err: any) {
      state.error = err.toString();
    } finally {
      state.loading = false;
    }
  }

  return { idle, currentJobId, parseEventInfo, state, fetchJobs, fetchSuccess, fetchJob, currentJob };
});
