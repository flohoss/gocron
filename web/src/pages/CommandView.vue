<script setup lang="ts">
import { useEventSource, useMagicKeys } from '@vueuse/core';
import CommandWindow from '../components/utils/CommandWindow.vue';
import { BackendURL } from '../main';
import { ref, watch, watchEffect } from 'vue';
import ArrowUp from '~icons/fa7-solid/arrow-up';
import ArrowDown from '~icons/fa7-solid/arrow-down';
import Terminal from '~icons/fa7-solid/terminal';
import { postCommand } from '../client/sdk.gen';
import { GetColor } from '../severity';
import { useCommands } from '../stores/useCommands';

const { responses, history, addResponse, addToHistory, clearResponses } = useCommands();

const { data, close } = useEventSource(BackendURL + '/api/events?stream=command', [], {
  autoReconnect: { delay: 100 },
});
addEventListener('beforeunload', () => {
  close();
});

watch(data, () => {
  const parsedResponse = JSON.parse(data.value);
  if (!parsedResponse) return;
  addResponse(parsedResponse);
});

const command = ref('');
const historyIndex = ref(-1);

function executeCommand() {
  postCommand({
    body: {
      command: command.value,
    },
  });
  addToHistory(command.value);
  historyIndex.value = -1;
  command.value = '';
}

function navigateHistory(direction: string) {
  if (direction === 'up') {
    if (historyIndex.value < history.value.length - 1) {
      historyIndex.value++;
      command.value = history.value[history.value.length - 1 - historyIndex.value] ?? '';
    }
  } else if (direction === 'down') {
    if (historyIndex.value > 0) {
      historyIndex.value--;
      command.value = history.value[history.value.length - 1 - historyIndex.value] ?? '';
    } else {
      historyIndex.value = -1;
      command.value = '';
    }
  }
}

const { ctrl, l } = useMagicKeys();

watchEffect(() => {
  if (ctrl?.value && l?.value) clearResponses();
});
</script>

<template>
  <CommandWindow :stickToBottom="true" title="Terminal">
    <pre v-for="(response, index) in responses" :key="index" :class="GetColor(response.severity)" class="flex">
      <code>{{ response.command }}</code>
    </pre>
    <template v-slot:top>
      <div class="hidden md:flex flex-wrap items-center gap-10 text-xs">
        <div class="flex items-center gap-2">
          Press
          <kbd class="kbd kbd-xs">ctrl</kbd>
          +
          <kbd class="kbd kbd-xs">l</kbd>
          to clear terminal
        </div>
        <div class="flex items-center gap-2">
          Press
          <kbd class="kbd kbd-xs"><ArrowUp /></kbd>
          or
          <kbd class="kbd kbd-xs"><ArrowDown /></kbd>
          to navigate history
        </div>
      </div>
    </template>
    <template v-slot:bottom>
      <div class="grid gap-5">
        <label class="input w-full">
          <Terminal />
          <input
            @keydown.up.prevent="navigateHistory('up')"
            @keydown.down.prevent="navigateHistory('down')"
            @keydown.esc="command = ''"
            @keydown.enter="executeCommand"
            v-model="command"
            autofocus
            type="text"
            placeholder="Command"
          />
        </label>
      </div>
    </template>
  </CommandWindow>
</template>
