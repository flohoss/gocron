<script setup lang="ts">
import { useRoute } from 'vue-router';
import ShortDuration from '../components/ShortDuration.vue';
import { useEventStore } from '../stores/event';
import { onUpdated, useTemplateRef, watch } from 'vue';

const route = useRoute();
const store = useEventStore();
const scrollContainer = useTemplateRef('scrollContainer');

onUpdated(() => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight;
  }
});

watch(() => route.params.id, store.fetchJob, { immediate: true });

enum Severity {
  Debug = 1,
  Info = 2,
  Warning = 3,
  Error = 4,
}

function getColor(severity: Severity): string {
  switch (severity) {
    case Severity.Info:
      return 'text-secondary';
    case Severity.Warning:
      return 'text-warning';
    case Severity.Error:
      return 'text-error';
    default:
      return 'text-base-content';
  }
}
</script>

<template>
  <div class="bg-base-300 w-full text-sm rounded-xl">
    <div class="flex justify-start items-center gap-2 padding bg-base-200 rounded-t-xl">
      <div class="console-btn bg-error text-error hover:text-error-content"></div>
      <div class="console-btn bg-warning text-warning hover:text-warning-content"></div>
      <div class="console-btn bg-success text-success hover:text-success-content"></div>
    </div>
    <div
      ref="scrollContainer"
      class="h-[calc(100vh-12rem)] md:h-[calc(100vh-14rem)] lg:h-[calc(100vh-18rem)] overflow-scroll padding"
      v-if="store.fetchSuccess"
    >
      <template v-for="(run, i) in store.currentJob!.runs" :key="i">
        <pre
          :id="`run-${i + 1}`"
          :class="getColor(Severity.Debug)"
        ><code>{{ run.fmt_start_time }}: Job <span class="text-primary font-bold">{{ store.currentJob!.name }}</span> started</code></pre>

        <template v-for="log in run.logs" :key="log.id">
          <span :class="[getColor(log.severity_id), 'flex']">
            <pre><code>{{ log.created_at_time }}: </code></pre>
            <pre><code>{{ log.message }}</code></pre>
          </span>
        </template>
        <pre
          v-if="run.duration.Valid && run.fmt_end_time.Valid"
          :class="getColor(Severity.Debug)"
          class="mb-2 last:mb-0"
        ><code>{{ run.fmt_end_time.String }}: Job finished (took <ShortDuration :duration="run.duration.Int64" />)</code></pre>
      </template>
    </div>
  </div>
</template>
