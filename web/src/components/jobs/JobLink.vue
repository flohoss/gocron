<script setup lang="ts">
import type { database_Job } from '@/openapi';
import NavLink from '../ui/NavLink.vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const props = defineProps<{
  job: database_Job;
}>();
const emit = defineEmits(['click']);

const isActive = computed(() => route.params.id === '' + props.job);

const jobIcon = () => {
  if (props.job.svg_icon) {
    return props.job.svg_icon;
  }
  return props.job.description.charAt(0);
};
</script>

<template>
  <NavLink
    @click="emit('click')"
    :name="job.description"
    :extra="job.local_directory"
    :link="{ name: 'jobs', params: { id: job.id } }"
    :active="isActive"
    :icon="jobIcon()"
    status="test"
    background="bg-primary"
    text="text-primary-content"
  />
</template>
