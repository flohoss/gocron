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

const startJob = async () => {
  try {
    await CommandsService.postCommands({ command: 'start', job_id: job.value.id });
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
</script>

<template>
  <div v-if="job">
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="text-xl font-bold">{{ job.description }}</div>
      <div class="join">
        <button @click="startJob" class="join-item btn btn-sm btn-neutral"><i class="fa-solid fa-play"></i>Run</button>
        <RouterLink :to="{ name: 'jobsForm', params: { id: job.id } }" class="join-item btn btn-sm btn-neutral">
          <i class="fa-solid fa-pencil"></i>Edit
        </RouterLink>
        <button @click="deleteJob" class="join-item btn btn-sm btn-error"><i class="fa-solid fa-trash"></i>Delete</button>
      </div>
    </PageHeader>
    <PageContent>
      <div v-for="(run, i) of job.runs" :key="i"><JobRun :run="run" :checked="i === 0" /></div>
    </PageContent>
  </div>
</template>
