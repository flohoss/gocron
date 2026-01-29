<script lang="ts" setup>
import { useJobs } from '../../stores/useJobs';
import { putJob } from '../../client/sdk.gen';
import Search from '~icons/fa7-solid/search';

const { jobs, filteredJobs, search, loading, checked } = useJobs();

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
      await putJob({ query: { name: id, action: 'toggle' } });
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

async function changeAction(action: 'disable_all' | 'enable_all' | 'enable_scheduled' | 'enable_non_scheduled') {
  await withMinLoading(async () => {
    await putJob({ query: { action: action } });
    checked.value = [...jobs.value.keys()];
    switch (action) {
      case 'disable_all':
        checked.value = [];
        break;
      case 'enable_all':
        checked.value = [...jobs.value.keys()];
        break;
      case 'enable_scheduled':
        checked.value = [];
        for (const job of jobs.value.values()) {
          if (!job.disable_cron) checked.value.push(job.name);
        }
        break;
      case 'enable_non_scheduled':
        checked.value = [];
        for (const job of jobs.value.values()) {
          if (job.disable_cron) checked.value.push(job.name);
        }
        break;
    }
  });
}
</script>

<template>
  <dialog id="selectModal" class="modal modal-bottom sm:modal-middle">
    <div class="modal-box grid gap-6">
      <div class="grid gap-1">
        <h3 class="text-lg font-bold">Select Jobs</h3>
        <div class="flex flex-wrap gap-2 text-secondary text-sm">
          <button :disabled="loading" @click="changeAction('enable_all')" class="link link-hover hover:text-primary">All Jobs</button>
          |
          <button :disabled="loading" @click="changeAction('disable_all')" class="link link-hover hover:text-primary">No Jobs</button>
          |
          <button :disabled="loading" @click="changeAction('enable_scheduled')" class="link link-hover hover:text-primary">All Scheduled Jobs</button>
          |
          <button :disabled="loading" @click="changeAction('enable_non_scheduled')" class="link link-hover hover:text-primary">All Non-Scheduled Jobs</button>
        </div>
      </div>
      <label class="input w-full">
        <Search />
        <input type="search" v-model="search" class="grow" placeholder="Search" />
      </label>
      <div class="grid sm:grid-cols-2 gap-2">
        <label class="label flex gap-5" v-for="[id] in filteredJobs.entries()" :key="id">
          <input
            @change.prevent="(e) => toggleJob(e, id)"
            :value="id"
            type="checkbox"
            class="toggle"
            :checked="checked.includes(id)"
            :class="checked.includes(id) ? 'toggle-primary' : 'toggle-neutral'"
            :disabled="loading"
          />
          <span class="truncate max-w-40">
            {{ id }}
          </span>
        </label>
      </div>
    </div>
    <form v-if="!loading" method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
