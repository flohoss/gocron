<script setup lang="ts">
import { onBeforeMount } from 'vue';
import HomeJob from '../components/HomeJob.vue';
import { useEventStore } from '../stores/event';

const store = useEventStore();
onBeforeMount(() => store.fetchHomeViewData());
</script>

<template>
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-8" v-if="store.homeViewSuccess">
    <HomeJob v-for="job in store.homeView.jobs" :key="job.id" :job="job" />
  </div>
</template>
