<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ title?: string; value: number; data: { percent: number; value: number; desc: string }[] }>();

const twoDigits = (val: number) => val % 100;
const rest = (val: number) => {
  const number = Math.floor(val / 100);
  if (number === 0) return '';
  else return number;
};

const typeColor = (index: number) => {
  switch (index) {
    case 0:
      return 'progress-error badge-error';
    case 1:
      return 'progress-warning badge-warning';
    case 2:
      return 'progress-success badge-success';
    case 3:
      return 'progress-primary badge-primary';
  }
};

const barBackground = computed(() => {
  for (let val of props.data) {
    if (val.value != 0) {
      return 'bg-transparent';
    }
  }
  return 'bg-neutral';
});
</script>

<template>
  <div class="flex w-full flex-col justify-between select-none">
    <div class="flex justify-between">
      <div>
        <div v-if="title" class="text-lg opacity-75">{{ title }}</div>
        <div class="text-6xl font-bold">
          <span>{{ rest(value) }}</span>
          <span class="countdown">
            <span :style="'--value:' + twoDigits(value) + ';'"></span>
          </span>
        </div>
      </div>
      <div class="flex justify-start flex-col-reverse text-sm gap-1">
        <div class="flex justify-end items-center gap-2" v-for="(run, index) of data" :key="index">
          <span>{{ rest(run.value) }}</span>
          <span class="countdown">
            <span :style="'--value:' + twoDigits(run.value) + ';'"></span>
          </span>
          <div :class="typeColor(index)" class="badge badge-sm px-1">{{ run.desc }}</div>
        </div>
      </div>
    </div>
    <div class="relative w-full h-6 mt-2">
      <div v-for="(run, index) of data" :key="index">
        <progress class="absolute progress w-full h-6" :class="typeColor(index) + ' ' + barBackground" :value="run.percent" max="100"></progress>
      </div>
    </div>
  </div>
</template>
