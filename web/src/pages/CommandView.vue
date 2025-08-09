<script setup lang="ts">
import { useEventSource } from '@vueuse/core';
import CommandWindow from '../components/utils/CommandWindow.vue';
import { BackendURL } from '../main';
import { ref, watch } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTerminal } from '@fortawesome/free-solid-svg-icons';
import { postCommand } from '../client/sdk.gen';
import { GetColor, Severity } from '../severity';

type CommandInfo = {
  severity: Severity;
  command: string;
};

const { data, close } = useEventSource(BackendURL + '/api/events?stream=command', [], {
  autoReconnect: { delay: 100 },
});
addEventListener('beforeunload', () => {
  close();
});

watch(data, () => {
  const parsedResponse = JSON.parse(data.value);
  if (!parsedResponse) return;
  responses.value.push(parsedResponse);
});

const responses = ref<CommandInfo[]>([]);
const command = ref('');

const executeCommand = () => {
  postCommand({
    body: {
      command: command.value,
    },
  });
  command.value = '';
};
</script>

<template>
  <CommandWindow :stickToBottom="true" title="Terminal">
    <pre v-for="(response, index) in responses" :key="index" :class="GetColor(response.severity)" class="flex">
      <code>{{ response.command }}</code>
    </pre>
    <template v-slot:bottom>
      <label class="input w-full">
        <FontAwesomeIcon :icon="faTerminal" />
        <input @keydown.esc="command = ''" @keydown.enter="executeCommand" v-model="command" autofocus type="text" placeholder="Command" />
      </label>
    </template>
  </CommandWindow>
</template>
