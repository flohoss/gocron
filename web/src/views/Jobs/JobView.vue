<script setup lang="ts">
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { useJobStore } from '@/stores/jobs';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { CommandsService, database_Job } from '@/openapi';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import JobRun from '@/components/jobs/JobRun.vue';

const store = useJobStore();
const route = useRoute();
const router = useRouter();
const job = computed<database_Job>(() => store.getJob(route.params.id));

const error = ref<string>('');
const errorModal = ref();

const startRun = async () => {
  try {
    await CommandsService.postCommands({ command: 'run', job_id: job.value.id });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};

const startPrune = async () => {
  try {
    await CommandsService.postCommands({ command: 'prune', job_id: job.value.id });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};

const startCheck = async () => {
  try {
    await CommandsService.postCommands({ command: 'check', job_id: job.value.id });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};

const deleteJob = async () => {
  if (job.value.id) {
    try {
      await store.deleteJob(job.value.id);
      router.push({ name: 'home' });
    } catch (err: any) {
      error.value = err.body.message;
      errorModal.value.showModal();
    }
  }
};

const badges = computed(() => {
  const badges = [
    { name: job.value.local_directory, icon: 'folder' },
    { name: job.value.restic_remote, icon: 'file-export' },
    { name: job.value.password_file_path, icon: 'key' },
  ];
  if (job.value.pre_commands) {
    for (const cmd of job.value.pre_commands) {
      badges.push({ name: cmd.command, icon: 'terminal' });
    }
  }
  if (job.value.post_commands) {
    for (const cmd of job.value.post_commands) {
      badges.push({ name: cmd.command, icon: 'terminal' });
    }
  }
  return badges;
});

const disabled = computed(() => {
  if (job.value.runs && job.value.runs?.length !== 0) {
    return !job.value.runs[0].end_time;
  }
  return true;
});
</script>

<template>
  <div>
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader :badges="badges">
      <div class="text-xl font-bold">{{ job.description }}</div>
      <div class="join">
        <button @click="startRun" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-play"></i><span class="hidden lg:block">Run</span>
        </button>
        <button @click="startPrune" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-broom"></i><span class="hidden lg:block">Prune</span>
        </button>
        <button @click="startCheck" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-check"></i><span class="hidden lg:block">Check</span>
        </button>
        <button @click="router.push({ name: 'jobsForm', params: { id: job.id } })" class="join-item btn btn-sm btn-warning" :disabled="disabled">
          <i class="fa-solid fa-pencil"></i><span class="hidden lg:block">Edit</span>
        </button>
        <button @click="deleteJob" class="join-item btn btn-sm btn-error" :disabled="disabled">
          <i class="fa-solid fa-trash"></i><span class="hidden lg:block">Delete</span>
        </button>
      </div>
    </PageHeader>
    <PageContent>
      <div class="grid grid-cols-1 gap-5 overflow-x-scroll">
        <JobRun v-for="(run, i) of job.runs" :key="run.id" :run="run" :checked="i === 0" />
      </div>
    </PageContent>
  </div>
</template>
