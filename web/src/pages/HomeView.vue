<script setup lang="ts">
import { computed } from 'vue';
import HomeJob from '../components/HomeJob.vue';
import HomeJobSkeleton from '../components/HomeJobSkeleton.vue';
import { useJobs } from '../stores/useJobs';

const { loadJobs, jobs } = useJobs();

const amount = computed(() => {
  if (jobs.value.size > 0) {
    return jobs.value.size;
  }
  return 2;
});
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
    <template v-if="loadJobs">
      <HomeJobSkeleton v-for="i in amount" :key="i" />
    </template>

    <template v-else>
      <HomeJob v-for="[id, job] in jobs" :key="id" :job="job" />
    </template>
  </div>
</template>
