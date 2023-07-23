<script setup lang="ts">
import type { database_Job } from "@/openapi";
import NavLink from "../ui/NavLink.vue";

const props = defineProps<{
  job: database_Job;
  active: boolean;
}>();
const emit = defineEmits(["click"]);

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
    :active="active"
    :icon="jobIcon()"
    status="test"
    background="bg-primary"
    text="text-primary-content"
  />
</template>
