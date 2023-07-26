<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ title?: string; value: number; desc?: string; class?: string }>();

const twoDigits = computed(() => props.value % 100);
const rest = computed(() => {
  const number = Math.floor(props.value / 100);
  if (number === 0) return '';
  else return number;
});
</script>

<template>
  <div class="stat place-items-center">
    <div v-if="title" class="stat-title" :class="class">{{ title }}</div>
    <div class="stat-value" :class="class">
      <span>{{ rest }}</span>
      <span class="countdown">
        <span :style="'--value:' + twoDigits + ';'"></span>
      </span>
    </div>
    <div v-if="desc" class="stat-desc" :class="class">{{ desc }}</div>
  </div>
</template>
