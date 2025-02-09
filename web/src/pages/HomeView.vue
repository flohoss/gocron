<script setup lang="ts">
import HomeJob from '../components/HomeJob.vue';
import { BackendURL } from '../main';
import { useFetch } from '@vueuse/core';

const { data } = useFetch(BackendURL + 'api/jobs')
  .get()
  .json();
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
    <HomeJob v-for="job in data" :key="job.id" :job="job" />
  </div>
</template>
