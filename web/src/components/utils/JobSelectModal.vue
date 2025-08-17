<script lang="ts" setup>
import { useJobs } from '../../stores/useJobs';
import { putJob } from '../../client/sdk.gen';

const { jobs, loading, checked } = useJobs();

async function withMinLoading<T>(fn: () => Promise<T>, minDuration = 500): Promise<T> {
  const start = Date.now();
  loading.value = true;

  try {
    const result = await fn();
    const elapsed = Date.now() - start;
    if (elapsed < minDuration) {
      await new Promise((resolve) => setTimeout(resolve, minDuration - elapsed));
    }
    return result;
  } finally {
    loading.value = false;
  }
}

async function toggleJob(event: Event, id: string) {
  const input = event.target as HTMLInputElement;
  const original = input.checked;

  await withMinLoading(async () => {
    try {
      await putJob({ query: { name: id } });
      if (original) {
        if (!checked.value.includes(id)) checked.value.push(id);
      } else {
        checked.value = checked.value.filter((x) => x !== id);
      }
    } catch {
      input.checked = !original;
    }
  });
}

async function enableAll() {
  await withMinLoading(async () => {
    await putJob({ query: { enable_all: true } });
    checked.value = [...jobs.value.keys()];
  });
}

async function enableScheduled() {
  await withMinLoading(async () => {
    await putJob({ query: { enable_scheduled: true } });
    const result: string[] = [];
    for (const job of jobs.value.values()) {
      if (!job.disable_cron) result.push(job.name);
    }
    checked.value = result;
  });
}

async function enableNonScheduled() {
  await withMinLoading(async () => {
    await putJob({ query: { enable_non_scheduled: true } });
    const result: string[] = [];
    for (const job of jobs.value.values()) {
      if (job.disable_cron) result.push(job.name);
    }
    checked.value = result;
  });
}
</script>

<template>
  <dialog id="selectModal" class="modal">
    <div class="modal-box grid gap-6">
      <div class="grid gap-1">
        <h3 class="text-lg font-bold">Select Jobs</h3>
        <div class="flex gap-2 text-secondary text-sm">
          <button :disabled="loading" @click="enableAll" class="link link-hover hover:text-primary">All Jobs</button>
          |
          <button :disabled="loading" @click="enableScheduled" class="link link-hover hover:text-primary">All Scheduled Jobs</button>
          |
          <button :disabled="loading" @click="enableNonScheduled" class="link link-hover hover:text-primary">All Non-Scheduled Jobs</button>
        </div>
      </div>
      <div class="grid gap-2">
        <label class="label flex gap-5" v-for="[id] in jobs.entries()" :key="id">
          <input
            @change.prevent="(e) => toggleJob(e, id)"
            :value="id"
            type="checkbox"
            class="toggle"
            :checked="checked.includes(id)"
            :class="checked.includes(id) ? 'toggle-primary' : 'toggle-neutral'"
            :disabled="loading"
          />
          {{ id }}
        </label>
      </div>
    </div>
    <form v-if="!loading" method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
