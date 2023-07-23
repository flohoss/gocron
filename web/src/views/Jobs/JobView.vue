<script setup lang="ts">
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { useJobStore } from '@/stores/jobs';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { CommandsService, database_Job } from '@/openapi';
import ErrorModal from '@/components/ui/ErrorModal.vue';

const store = useJobStore();
const route = useRoute();
const router = useRouter();
const job = computed<database_Job>(() => store.getJob(route.params.id));

const error = ref<string>('');
const errorModal = ref();

const startJob = async () => {
  try {
    await CommandsService.postCommands({ command: 'start', job_id: job.value.id });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};

const deleteJob = async () => {
  try {
    await store.deleteJob(job.value.id);
    router.push({ name: 'home' });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};
</script>

<template>
  <div v-if="job">
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="flex flex-col gap-1 items-center lg:items-start">
        <div class="text-xl font-bold">{{ job.description }}</div>
        <div class="flex items-center flex-wrap">
          <div class="badge badge-xs gap-1"><i class="fa-solid fa-file-export"></i>{{ job.local_directory }}</div>
          <div class="badge badge-xs gap-1"><i class="fa-solid fa-file-import"></i> {{ job.restic_remote }}</div>
        </div>
      </div>
      <div class="join">
        <button @click="startJob" class="join-item btn btn-sm btn-neutral"><i class="fa-solid fa-play"></i>Run</button>
        <RouterLink :to="{ name: 'jobsForm', params: { id: job.id } }" class="join-item btn btn-sm btn-neutral">
          <i class="fa-solid fa-pencil"></i>Edit
        </RouterLink>
        <button @click="deleteJob" class="join-item btn btn-sm btn-error"><i class="fa-solid fa-trash"></i>Delete</button>
      </div>
    </PageHeader>
    <PageContent>
      <div v-for="run of job.runs" :key="run.id">{{ run }}</div>
    </PageContent>
  </div>
</template>
