<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faCheck, faQuestion, faTimes } from '@fortawesome/free-solid-svg-icons';
import type { RunView } from '../client/types.gen';
import { computed, onMounted, onUnmounted, ref } from 'vue';

const { run } = defineProps<{ run: RunView }>();

enum Status {
  Running = 1,
  Stopped = 2,
  Finished = 3,
}

function getStepColor(status: Status): string {
  switch (status) {
    case Status.Running:
      return 'step-warning';
    case Status.Stopped:
      return 'step-error';
    case Status.Finished:
      return 'step-success';
    default:
      return 'step-neutral';
  }
}

function getStepIcon(status: Status) {
  switch (status) {
    case Status.Stopped:
      return faTimes;
    case Status.Finished:
      return faCheck;
    default:
      return faQuestion;
  }
}

const elapsedSeconds = ref(0);
const timer = ref<ReturnType<typeof setInterval> | null>(null);

function formatDuration(seconds: number) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;

  let result = '';
  if (h > 0) result += h + 'h';
  if (m > 0) result += m + 'm';
  if (s > 0 || seconds === 0) result += s + 's';
  return result;
}

const duration = computed(() => {
  if (run.duration !== '0s') {
    return run.duration;
  } else {
    return formatDuration(elapsedSeconds.value);
  }
});

onMounted(() => {
  if (run.duration === '0s') {
    const now = Date.now();
    const startTimeMs = run.start_time_unix;
    elapsedSeconds.value = Math.floor((now - startTimeMs) / 1000);

    timer.value = setInterval(() => {
      elapsedSeconds.value++;
    }, 1000);
  }
});

onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value);
  }
});
</script>

<template>
  <li class="step" :class="getStepColor(run.status_id)">
    <span class="step-icon">
      <span v-if="run.status_id === Status.Running" class="loading loading-spinner"></span>
      <FontAwesomeIcon v-else :icon="getStepIcon(run.status_id)" class="size-6" />
    </span>
    {{ duration }}
  </li>
</template>

<style scoped>
.steps .step::before {
  height: 0.2rem !important;
}
</style>
