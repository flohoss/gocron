<script setup lang="ts">
import type { Moment } from 'moment';
import moment from 'moment';
import { computed, ref } from 'vue';

const props = defineProps<{ title: string; cron: Moment }>();

const currentTime = ref(moment());

const diffDuration = computed(() => moment.duration(props.cron.diff(currentTime.value)));
const hours = computed(() => diffDuration.value.hours());
const minutes = computed(() => diffDuration.value.minutes());
const seconds = computed(() => diffDuration.value.seconds());
setInterval(() => (currentTime.value = moment()), 1000);
</script>

<template>
  <div class="card compact">
    <div class="card-body">
      <h2 class="card-title">{{ title }}</h2>
      <div class="grid grid-flow-col gap-5 text-center auto-cols-max">
        <div class="flex flex-col p-2 bg-neutral rounded-box text-neutral-content">
          <span class="countdown font-mono text-5xl">
            <span :style="'--value: ' + hours"></span>
          </span>
          hours
        </div>
        <div class="flex flex-col p-2 bg-neutral rounded-box text-neutral-content">
          <span class="countdown font-mono text-5xl">
            <span :style="'--value: ' + minutes"></span>
          </span>
          min
        </div>
        <div class="flex flex-col p-2 bg-neutral rounded-box text-neutral-content">
          <span class="countdown font-mono text-5xl">
            <span :style="'--value: ' + seconds"></span>
          </span>
          sec
        </div>
      </div>
    </div>
  </div>
</template>
