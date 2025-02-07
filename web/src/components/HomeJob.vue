<script setup lang="ts">
import { computed } from 'vue';
import { BackendURL } from '../main';
import type { JobsService } from '../openapi';

const props = defineProps<{ job:  }>();
const url = computed<JobsService>(() => BackendURL + props.id);
</script>

<template>
  <a class="flex justify-between items-center group last:mb-8 lg:last:mb-0" :href="url">
    <div class="pl-4 truncate">
      <div class="group-hover:text-primary hover-animation text-2xl font-medium truncate">{ result.Name }</div>
      <div class="text-secondary text-sm truncate">{ parseCron(result.Cron, parser) }</div>
    </div>
    <div class="text-sm">
      <ul class="steps hidden lg:flex">
        <li>{ parseDuration(run.Duration.String) }</li>
      </ul>
    </div>
  </a>
</template>
