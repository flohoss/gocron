import { createGlobalState } from '@vueuse/core';
import { ref } from 'vue';
import type { Severity } from '../severity';

type CommandInfo = {
  severity: Severity;
  command: string;
};

export const useCommands = createGlobalState(() => {
  const responses = ref<CommandInfo[]>([]);
  const history = ref<string[]>([]);

  function addResponse(response: CommandInfo) {
    responses.value.push(response);
  }

  function addToHistory(command: string) {
    history.value.push(command);
  }

  function clearResponses() {
    responses.value = [];
  }

  return {
    responses,
    history,
    addResponse,
    addToHistory,
    clearResponses,
  };
});
