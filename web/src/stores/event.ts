import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useEventStore = defineStore('event', () => {
  const job = ref(null);
  const idle = ref<boolean>(false);

  function run() {
    let url = '/api/jobs';
    if (job.value) {
      url = url + '/' + job.value;
    }
    fetch(url, { method: 'POST' }).then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    });
  }

  return { idle, job, run };
});
