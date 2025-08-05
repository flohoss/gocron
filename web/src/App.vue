<script setup lang="ts">
import { RouterView } from 'vue-router';
import AppHeader from './components/AppHeader.vue';
import { useEventStore } from './stores/event';
import { onBeforeMount } from 'vue';

const store = useEventStore();
onBeforeMount(() => store.fetchJobs());
</script>

<template>
  <AppHeader />
  <main>
    <RouterView v-slot="{ Component }">
      <Transition mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </main>
</template>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease-out;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
