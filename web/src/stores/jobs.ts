import { ref } from 'vue';
import { defineStore } from 'pinia';
import { JobsService, type database_Job } from '@/openapi';
import { emptyJob } from '@/types';

export const useJobStore = defineStore('jobs', () => {
  const jobs = ref<database_Job[]>([]);

  async function fetchJobs() {
    jobs.value = await JobsService.getJobs();
  }

  async function createJob(job: database_Job) {
    const response = await JobsService.postJobs(job);
    jobs.value.push(response);
    sortJobs();
    return response;
  }

  async function updateJob(job: database_Job) {
    const response = await JobsService.putJobs(job);
    jobs.value = jobs.value.map((j) => {
      if (j.id === job.id) return response;
      else return j;
    });
    sortJobs();
  }

  async function deleteJob(id: number) {
    await JobsService.deleteJobs(id);
    jobs.value = jobs.value.filter((job) => job.id !== id);
  }

  function getJob(strId: string | string[]) {
    const id = parseInt(strId + '');
    const job = jobs.value.find((job) => job.id === id);
    if (job) return job;
    else return { ...emptyJob };
  }

  function sortJobs() {
    jobs.value.sort((a: database_Job, b: database_Job) => {
      const nameA = a.description.toUpperCase();
      const nameB = b.description.toUpperCase();
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    });
  }

  return { jobs, fetchJobs, createJob, updateJob, deleteJob, getJob };
});
