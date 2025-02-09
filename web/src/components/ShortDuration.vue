<script setup lang="ts">
import humanizeDuration from 'humanize-duration';
import { computed } from 'vue';

const props = defineProps<{ duration: number }>();

const shortEnglishHumanizer = humanizeDuration.humanizer({
  language: 'shortEn',
  languages: {
    shortEn: {
      y: () => 'y',
      mo: () => 'mo',
      w: () => 'w',
      d: () => 'd',
      h: () => 'h',
      m: () => 'm',
      s: () => 's',
      ms: () => 'ms',
    },
  },
});
const duration = computed(() => shortEnglishHumanizer(props.duration, { round: true, largest: 1 }));
</script>

<template>{{ duration }}</template>
