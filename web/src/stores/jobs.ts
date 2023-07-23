import { ref } from 'vue';
import { defineStore } from 'pinia';
import { JobsService, type database_Job } from '@/openapi';

export const useJobStore = defineStore('jobs', () => {
  const jobs = ref<database_Job[]>([]);

  async function getJobs() {
    const response = await JobsService.getJobs();
    jobs.value = response;
  }

  async function createJob(job: database_Job) {
    const response = await JobsService.postJobs(job);
    jobs.value.push(response);
    jobs.value.sort();
  }

  return { jobs, getJobs, createJob };
});
