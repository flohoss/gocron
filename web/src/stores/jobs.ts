import { ref } from 'vue';
import { defineStore } from 'pinia';
import { JobsService, type database_Job, type database_Log, type database_Run } from '@/openapi';

export const emptyJob: database_Job = {
  compression_type_id: 1,
  description: '',
  id: 0,
  local_directory: '',
  password_file_path: '',
  post_commands: [],
  pre_commands: [],
  restic_remote: '',
  retention_policy_id: 1,
  routine_check: 0,
  svg_icon: '',
};

export const useJobStore = defineStore('jobs', () => {
  const jobs = ref<database_Job[]>([]);

  function updateOrCreateRun(run: database_Run) {
    const jobindex = jobs.value.findIndex((job) => job.id === run.job_id);
    if (jobindex === -1) return;
    const runIndex = jobs.value[jobindex].runs!.findIndex((r) => r.id === run.id);

    if (runIndex === -1) {
      jobs.value[jobindex].runs?.unshift(run);
    } else {
      jobs.value[jobindex].runs![runIndex] = run;
    }
  }

  function updateOrCreateLog(jobId: number | undefined, log: database_Log) {
    const jobindex = jobs.value.findIndex((job) => job.id === jobId);
    if (jobindex === -1) return;
    const runIndex = jobs.value[jobindex].runs!.findIndex((r) => r.id === log.run_id);

    if (!jobs.value[jobindex].runs![runIndex].logs) {
      jobs.value[jobindex].runs![runIndex].logs = [];
    }
    jobs.value[jobindex].runs![runIndex].logs?.push(log);
  }

  async function getJobs() {
    const response = await JobsService.getJobs();
    jobs.value = response;
  }

  async function createJob(job: database_Job) {
    const response = await JobsService.postJobs(job);
    jobs.value.push(response);
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
    return response;
  }

  async function updateJob(job: database_Job) {
    const response = await JobsService.putJobs(job);
    jobs.value = jobs.value.map((j) => {
      if (j.id === job.id) return response;
      else return j;
    });
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

  return { jobs, updateOrCreateRun, updateOrCreateLog, getJobs, createJob, updateJob, deleteJob, getJob };
});
