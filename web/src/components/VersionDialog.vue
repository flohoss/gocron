<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import type { Information } from '../client/types.gen';
import { getVersions } from '../client/sdk.gen';

defineProps<{ id: string }>();
const versions = ref<Information[] | null>(null);

onBeforeMount(async () => {
  const res = await getVersions();
  versions.value = res.data!;
});
</script>

<template>
  <dialog :id="id" class="modal modal-bottom sm:modal-middle">
    <div class="modal-box">
      <div class="overflow-x-auto">
        <table class="table">
          <thead>
            <tr>
              <th>Software</th>
              <th>Installed</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="version in versions" :key="version.name">
              <td>{{ version.name }}</td>
              <td>
                <a target="_blank" class="link" :href="version.repo + version.version">{{ version.version }}</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
