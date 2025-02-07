<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const loading = ref(false);
const post = ref(null);
const error = ref(null);

watch(() => route.params.id, fetchData, { immediate: true });

async function fetchData(id: string | string[]) {
  error.value = post.value = null;
  loading.value = true;

  try {
    id = id as string;
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
}
</script>

<template>Job</template>
