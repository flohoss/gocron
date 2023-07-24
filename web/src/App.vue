<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router';
import { useJobStore } from './stores/jobs';
import NavLink from './components/ui/NavLink.vue';
import JobLink from './components/jobs/JobLink.vue';
import { computed, ref } from 'vue';
import ErrorModal from './components/ui/ErrorModal.vue';

const store = useJobStore();
const route = useRoute();
const error = ref<string>('');
const errorModal = ref();

const init = async () => {
  try {
    await store.getJobs();
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};
init();

const drawerRef = ref();
</script>

<template>
  <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
  <div class="drawer lg:drawer-open">
    <input id="drawer" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <RouterView v-slot="{ Component }">
        <Transition mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
      <div class="my-20 lg:hidden"></div>
      <div class="lg:hidden btm-nav bg-base-200">
        <RouterLink :to="{ name: 'home' }" :class="{ active: route.name === 'home' }">
          <i class="fa-solid fa-circle-nodes"></i>
          <div class="text-xs opacity-75">Dashboard</div>
        </RouterLink>
        <label for="drawer" ref="drawerRef">
          <i class="fa-solid fa-list-ul"></i>
          <div class="text-xs opacity-75">Jobs</div>
        </label>
        <RouterLink :to="{ name: 'jobsForm' }" :class="{ active: route.name === 'jobsForm' && !route.params.id }">
          <i class="fa-solid fa-plus"></i>
          <div class="text-xs opacity-75">New</div>
        </RouterLink>
      </div>
    </div>
    <div class="drawer-side">
      <label for="drawer" class="drawer-overlay"></label>
      <ul class="menu p-2 w-80 h-full bg-base-200 text-base-content flex flex-col flex-nowrap gap-4 overflow-y-auto">
        <NavLink
          :link="{ name: 'home' }"
          name="Dashboard"
          icon="<i class='fa-solid fa-circle-nodes'></i>"
          :active="route.name === 'home'"
          :small-hidden="true"
        />
        <JobLink
          v-for="job in store.jobs"
          :key="job.id"
          :job="job"
          @click="drawerRef && drawerRef.click()"
          :active="route.name === 'jobs' && parseInt(route.params.id + '') === job.id"
        />
        <NavLink
          :link="{ name: 'jobsForm' }"
          name="New"
          icon="<i class='fa-solid fa-plus'></i>"
          :active="route.name === 'jobsForm' && !route.params.id"
          :small-hidden="true"
        />
      </ul>
    </div>
  </div>
</template>

<style>
.v-enter-active {
  transition: opacity 0.5s ease;
}

.v-enter-from {
  opacity: 0;
}
</style>
