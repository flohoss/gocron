<script lang="ts" setup>
import { ref } from 'vue';
import { useJobs } from '../../stores/useJobs';
import { putJob } from '../../client/sdk.gen';

const { jobs, checked } = useJobs();
const loading = ref(false);

async function toggleJob(event: Event, id: string) {
  const input = event.target as HTMLInputElement;
  const original = input.checked;

  loading.value = true;

  try {
    await putJob({ path: { name: id } });
    if (original) {
      if (!checked.value.includes(id)) checked.value.push(id);
    } else {
      checked.value = checked.value.filter((x) => x !== id);
    }
  } catch {
    input.checked = !original;
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <dialog id="selectModal" class="modal">
    <div class="modal-box">
      <h3 class="text-lg font-bold">Select Jobs</h3>
      <div class="grid gap-5 mt-5">
        <label class="label" v-for="[id] in jobs.entries()" :key="id">
          <input
            @change="(e) => toggleJob(e, id)"
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
      <div class="modal-action">
        <form method="dialog">
          <button class="btn btn-error">Cancel</button>
        </form>
      </div>
    </div>
  </dialog>
</template>
