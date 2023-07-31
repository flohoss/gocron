<script setup lang="ts">
import { ref, onMounted } from 'vue';

const showComponent = ref(false);

onMounted(() => {
  setTimeout(() => {
    showComponent.value = true;
  }, 1000);
});
</script>

<template>
  <div v-if="showComponent" class="flex justify-center mt-5">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
</template>
