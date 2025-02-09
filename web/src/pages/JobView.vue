<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { JobsService, type services_TemplateJob } from '../openapi';

const route = useRoute();

const loading = ref(false);
const job = ref<services_TemplateJob | null>(null);
const error = ref(null);

watch(() => route.params.id, fetchData, { immediate: true });

async function fetchData(id: string | string[]) {
  error.value = job.value = null;
  loading.value = true;

  try {
    job.value = await JobsService.getJobs1(id + '');
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
}
</script>

<template>{{ job }}</template>
