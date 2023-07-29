<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router';
import { useJobStore } from './stores/jobs';
import NavLink from './components/ui/NavLink.vue';
import JobLink from './components/jobs/JobLink.vue';
import { ref } from 'vue';

const store = useJobStore();
const route = useRoute();

const drawerRef = ref();
store.fetchJobs();
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
        <label for="drawer" ref="drawerRef">
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
    <div class="drawer-side">
      <label for="drawer" class="drawer-overlay"></label>
      <ul class="menu p-2 w-80 h-full bg-base-200 text-base-content flex flex-col flex-nowrap overflow-y-auto">
        <NavLink :link="{ name: 'home' }" name="System" icon="<i class='fa-solid fa-circle-nodes'></i>" :active="route.name === 'home'" :small-hidden="true" />
        <div class="my-2"></div>
        <div class="grid gap-1">
          <JobLink
            v-for="job in store.jobs"
            :key="job.id"
            :job="job"
            @click="drawerRef && drawerRef.click()"
            :active="route.name === 'jobs' && parseInt(route.params.id + '') === job.id"
          />
        </div>
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
