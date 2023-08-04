<script setup lang="ts">
import type { Moment } from 'moment';
import moment from 'moment';
import { computed, ref } from 'vue';

const props = defineProps<{ title: string; cron: Moment }>();

const currentTime = ref(moment());

const diffDuration = computed(() => moment.duration(props.cron.diff(currentTime.value)));
const days = computed(() => diffDuration.value.days());
const hours = computed(() => diffDuration.value.hours());
const minutes = computed(() => diffDuration.value.minutes());
const seconds = computed(() => diffDuration.value.seconds());
setInterval(() => (currentTime.value = moment()), 1000);
</script>

<template>
  <div class="card compact">
    <div class="card-body">
      <div class="card-title text-lg">{{ title }}</div>
      <div class="grid grid-flow-col gap-4 text-center auto-cols-max">
        <div class="flex flex-col p-2 bg-base-200 rounded-box text-base-content">
          <span class="countdown font-mono text-4xl md:text-5xl">
            <span :style="'--value: ' + days"></span>
          </span>
          days
        </div>
        <div class="flex flex-col p-2 bg-base-200 rounded-box text-base-content">
          <span class="countdown font-mono text-4xl md:text-5xl">
            <span :style="'--value: ' + hours"></span>
          </span>
          hours
        </div>
        <div class="flex flex-col p-2 bg-base-200 rounded-box text-base-content">
          <span class="countdown font-mono text-4xl md:text-5xl">
            <span :style="'--value: ' + minutes"></span>
          </span>
          min
        </div>
        <div class="flex flex-col p-2 bg-base-200 rounded-box text-base-content">
          <span class="countdown font-mono text-4xl md:text-5xl">
            <span :style="'--value: ' + seconds"></span>
          </span>
          sec
        </div>
      </div>
    </div>
  </div>
</template>
