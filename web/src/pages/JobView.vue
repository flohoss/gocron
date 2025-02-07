<script setup>
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const loading = ref(false);
const post = ref(null);
const error = ref(null);

watch(() => route.params.id, fetchData, { immediate: true });

async function fetchData(id) {
  error.value = post.value = null;
  loading.value = true;

  try {
  } catch (err) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
}
</script>

<template>Job</template>
