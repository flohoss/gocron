<script setup lang="ts">
import { useEventSource } from '@vueuse/core';
import CommandWindow from '../components/utils/CommandWindow.vue';
import { BackendURL } from '../main';
import { ref, watch } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTerminal } from '@fortawesome/free-solid-svg-icons';

const { data, close } = useEventSource(BackendURL + '/api/events?stream=status', [], {
  autoReconnect: { delay: 100 },
});
addEventListener('beforeunload', () => {
  close();
});
watch(data.value, () => responses.value.push(data.value));

const responses = ref<string[]>([]);
</script>

<template>
  <CommandWindow :stickToBottom="true">
    <pre v-for="(response, index) in responses" :key="index"><code>{{ response }}</code></pre>
    <template v-slot:bottom>
      <label class="input w-full">
        <FontAwesomeIcon :icon="faTerminal" />
        <input type="text" placeholder="Command" />
      </label>
    </template>
  </CommandWindow>
</template>
