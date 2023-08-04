<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router';
import { useJobStore } from './stores/jobs';
import NavLink from './components/ui/NavLink.vue';
import JobLink from './components/jobs/JobLink.vue';
import { ref, provide, onBeforeUnmount, watch } from 'vue';
import { useEventSource } from '@vueuse/core';
import { EventType, sseKey, type SSEvent } from './types';
import { type database_Run, database_LogSeverity } from './openapi';
import LoadingSpinner from './components/ui/LoadingSpinner.vue';

const store = useJobStore();
const route = useRoute();

const drawerRef = ref();
const loading = ref(true);
store.fetchJobs().finally(() => (loading.value = false));

const { data, close } = useEventSource('/api/sse?stream=jobs');
const parsed = ref<SSEvent>();
watch(data, (value) => {
  parsed.value = JSON.parse(value + '');
  if (!parsed.value) return;
  switch (parsed.value.event_type) {
    case EventType.EventCreateRun: {
      const run = parsed.value.content as database_Run;
      const job = store.jobs.find((j) => j.id === run.job_id);
      if (job) job.status = database_LogSeverity.LogNone;
      break;
    }
    case EventType.EventUpdateRun: {
      const run = parsed.value.content as database_Run;
      const job = store.jobs.find((j) => j.id === run.job_id);
      if (job) job.status = run.status;
      break;
    }
  }
});
provide(sseKey, parsed);
onBeforeUnmount(() => close());
</script>

<template>
  <div class="drawer lg:drawer-open">
    <input id="drawer" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <RouterView v-slot="{ Component }">
        <Transition mode="out-in">
          <component :key="route.fullPath" :is="Component" />
        </Transition>
      </RouterView>
      <div class="my-20 lg:hidden"></div>
      <div class="lg:hidden btm-nav bg-base-200">
        <RouterLink :to="{ name: 'home' }" :class="{ active: route.name === 'home' }">
          <i class="fa-solid fa-circle-nodes"></i>
          <div class="text-xs opacity-75">System</div>
        </RouterLink>
        <label for="drawer" ref="drawerRef" :class="{ active: route.name === 'jobs' }">
          <i class="fa-solid fa-list-ul"></i>
          <div class="text-xs opacity-75">Jobs</div>
        </label>
        <RouterLink :to="{ name: 'jobsForm' }" :class="{ active: route.name === 'jobsForm' && !route.params.id }">
          <i class="fa-solid fa-plus"></i>
          <div class="text-xs opacity-75">New</div>
        </RouterLink>
        <RouterLink :to="{ name: 'jobsRestore' }" :class="{ active: route.name === 'jobsRestore' }">
          <i class="fa-solid fa-file-arrow-down"></i>
          <div class="text-xs opacity-75">Restore</div>
        </RouterLink>
      </div>
    </div>
    <div class="drawer-side z-50">
      <label for="drawer" class="drawer-overlay"></label>
      <ul class="menu p-2 w-80 h-full bg-base-200 text-base-content flex flex-col flex-nowrap overflow-y-auto">
        <NavLink :link="{ name: 'home' }" name="System" icon="<i class='fa-solid fa-circle-nodes'></i>" :active="route.name === 'home'" :small-hidden="true" />
        <div class="my-2"></div>
        <Transition>
          <div class="grid gap-1" v-if="!loading">
            <JobLink
              v-for="job in store.jobs"
              :key="job.id"
              :job="job"
              @click="drawerRef && drawerRef.click()"
              :active="route.name === 'jobs' && parseInt(route.params.id + '') === job.id"
            />
          </div>
          <LoadingSpinner v-else />
        </Transition>
        <div class="flex-grow my-2"></div>
        <div class="grid gap-1">
          <NavLink
            :link="{ name: 'jobsForm' }"
            name="New"
            icon="<i class='fa-solid fa-plus'></i>"
            :active="route.name === 'jobsForm' && !route.params.id"
            :small-hidden="true"
          />
          <NavLink
            :link="{ name: 'jobsRestore' }"
            name="Restore"
            icon="<i class='fa-solid fa-file-arrow-down'></i>"
            :active="route.name === 'jobsRestore'"
            :small-hidden="true"
          />
        </div>
      </ul>
    </div>
  </div>
</template>

<style>
.v-enter-active {
  transition: all 0.3s ease;
}

.v-enter-from {
  opacity: 0;
  transform: translateY(20px);
}
</style>
