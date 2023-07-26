<script setup lang="ts">
defineProps<{ title?: string; value: string | number; desc?: string; class?: string }>();
</script>

<template>
  <div class="stat place-items-center">
    <div v-if="title" class="stat-title" :class="class">{{ title }}</div>
    <div class="stat-value" :class="class">
      <span class="countdown">
        <span :style="'--value:' + value + ';'"></span>
      </span>
    </div>
    <div v-if="desc" class="stat-desc" :class="class">{{ desc }}</div>
  </div>
</template>
