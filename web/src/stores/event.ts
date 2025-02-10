import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import { JobsService, type events_EventInfo, type jobs_JobsView, type services_TemplateJob } from '../openapi';

export const useEventStore = defineStore('event', () => {
  const idle = ref<boolean>(false);
  const currentJobId = ref<string | null>(null);

  function parseEventInfo(info: string | null): void {
    if (!info) return;
    const parsed: events_EventInfo = JSON.parse(info);
    idle.value = parsed.idle;
    console.log(parsed);
  }

  const homeView = reactive<{ loading: boolean; error: string | null; jobs: jobs_JobsView[] | null }>({
    loading: false,
    error: null,
    jobs: null,
  });
  const homeViewSuccess = computed(() => homeView.error === null && homeView.loading === false && homeView.jobs !== null);

  async function fetchHomeViewData() {
    homeView.error = homeView.jobs = null;
    homeView.loading = true;

    try {
      homeView.jobs = await JobsService.getJobs();
    } catch (err: any) {
      homeView.error = err.toString();
    } finally {
      homeView.loading = false;
    }
  }

  const jobView = reactive<{ loading: boolean; error: string | null; job: services_TemplateJob | null }>({
    loading: false,
    error: null,
    job: null,
  });
  const jobViewSuccess = computed(() => jobView.error === null && jobView.loading === false && jobView.job !== null);

  async function fetchJobViewData(id: string | string[]) {
    currentJobId.value = id + '';
    jobView.error = jobView.job = null;
    jobView.loading = true;

    try {
      jobView.job = await JobsService.getJobs1(id + '');
    } catch (err: any) {
      jobView.error = err.toString();
    } finally {
      jobView.loading = false;
    }
  }

  return { idle, currentJobId, parseEventInfo, homeView, fetchHomeViewData, homeViewSuccess, jobView, fetchJobViewData, jobViewSuccess };
});
