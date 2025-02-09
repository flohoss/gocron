<script setup lang="ts">
import HomeJob from '../components/HomeJob.vue';
import { JobsService, type jobs_JobsView } from '../openapi';
import { ref } from 'vue';

const loading = ref(false);
const jobs = ref<jobs_JobsView[] | null>(null);
const error = ref(null);

async function fetchData() {
  error.value = jobs.value = null;
  loading.value = true;

  try {
    jobs.value = await JobsService.getJobs();
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
}

fetchData();
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
    <HomeJob v-for="job in jobs" :key="job.id" :job="job" />
  </div>
</template>
