import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface EventInfo {
  idle: boolean;
  job: Job;
}

export interface Job {
  id: string;
  name: string;
  log: Log;
}

export interface Log {
  created_at: number;
  severity_id: number;
  message: string;
}

const emptyEvent: EventInfo = { idle: false, job: { id: '', name: '', log: { created_at: 0, severity_id: 0, message: '' } } };

export const useEventStore = defineStore('event', () => {
  const job = ref(null);
  const event = ref<EventInfo>(emptyEvent);

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

  return { event, job, run };
});
