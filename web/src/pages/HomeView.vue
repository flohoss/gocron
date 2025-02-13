<script setup lang="ts">
import { onBeforeMount } from 'vue';
import HomeJob from '../components/HomeJob.vue';
import { useEventStore } from '../stores/event';

const store = useEventStore();
onBeforeMount(() => store.fetchJobs());
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8" v-if="store.fetchSuccess">
    <HomeJob v-for="[id, job] in store.state.jobs" :key="id" :job="job" />
  </div>
</template>
