<script setup lang="ts">
import { onBeforeMount, computed } from 'vue';
import HomeJob from '../components/HomeJob.vue';
import HomeJobSkeleton from '../components/HomeJobSkeleton.vue';
import { useEventStore } from '../stores/event';

const store = useEventStore();
onBeforeMount(() => store.fetchJobs());

const amount = computed(() => {
  if (store.state.jobs.size > 0) {
    return store.state.jobs.size;
  }
  return 2;
});
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
    <template v-if="store.state.loading">
      <HomeJobSkeleton v-for="i in amount" :key="i" />
    </template>

    <template v-else-if="store.fetchSuccess">
      <HomeJob v-for="[id, job] in store.state.jobs" :key="id" :job="job" />
    </template>
  </div>
</template>
