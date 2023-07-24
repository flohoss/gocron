<script setup lang="ts">
import { useSystemStore } from '@/stores/system';
import { ref } from 'vue';

const store = useSystemStore();
const error = ref<string>('');
const errorModal = ref();

const init = async () => {
  try {
    await store.getSystem();
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};
init();
</script>

<template>
  <div class="grid gap-5"></div>
</template>
