<script setup lang="ts">
import { useSystemStore } from "@/stores/system";
import { computed } from "vue";
import { ref } from "vue";

const store = useSystemStore();
const error = ref<string>("");
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

const storagePercentage = computed(() => ((store.system.disk.used / store.system.disk.total) * 100).toFixed(0));
</script>

<template>
  <div class="grid gap-5">{{ storagePercentage }}</div>
</template>
