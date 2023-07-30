<script setup lang="ts">
import { database_LogSeverity, type database_Job } from '@/openapi';
import NavLink from '../ui/NavLink.vue';
import { computed } from 'vue';
import { severityColor, severityIcons } from '@/helper/severity';

const props = defineProps<{
  job: database_Job;
  active: boolean;
}>();
const emit = defineEmits(['click']);

const jobIcon = () => props.job.description.charAt(0);

const color = computed(() => severityColor(props.job.status));
const icon = computed(() => severityIcons(props.job.status));
</script>

<template>
  <NavLink
    @click="emit('click')"
    :name="job.description"
    :extra="job.local_directory"
    :link="{ name: 'jobs', params: { id: job.id } }"
    :active="active"
    :icon="jobIcon()"
  >
    <span v-if="job.status === database_LogSeverity.LogNone" class="loading loading-spinner"></span>
    <span v-else><i :class="color + ' fa-solid fa-' + icon"></i></span>
  </NavLink>
</template>
