import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { events_EventInfo } from '../openapi';

export const useEventStore = defineStore('event', () => {
  const job = ref(null);
  const event = ref<events_EventInfo | null | undefined>(null);

  return { event, job };
});
