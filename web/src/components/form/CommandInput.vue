<script setup lang="ts">
import useVuelidate from '@vuelidate/core';
import { required } from '@vuelidate/validators';
import TextInput from './TextInput.vue';
import { reactive, watch } from 'vue';

const props = defineProps<{ id: string; command: string; fileOutput: string; index: number; amount: number }>();
const emit = defineEmits(['update:command', 'update:fileOutput', 'handleMoveUp', 'handleMoveDown', 'handleRemoveCommand']);

const state = reactive({
  command: props.command,
  fileOutput: props.fileOutput,
});

const rules = {
  command: { required },
  fileOutput: {},
};
const v$ = useVuelidate(rules, state);

watch(state, () => {
  emit('update:command', state.command);
  emit('update:fileOutput', state.fileOutput);
});
</script>

<template>
  <div class="flex gap-x-5 flex-col lg:flex-row items-end lg:items-center">
    <TextInput
      class="lg:w-4/6"
      :id="id + '[' + index + '].command'"
      :title="index + 1 + '. Command'"
      v-model="v$.command.$model"
      help="Example: docker compose up"
      :v$="v$.command"
    />
    <TextInput
      class="lg:w-2/6"
      :id="id + '[' + index + '].file'"
      :title="index + 1 + '. File output'"
      v-model="v$.fileOutput.$model"
      help="Example: ./dbBackup.sql"
      :v$="v$.fileOutput"
    />
    <div class="join lg:join-vertical pt-1">
      <button class="join-item btn btn-neutral" @click="emit('handleMoveUp', index)" type="button" v-if="index !== 0">
        <i class="fa-solid fa-arrow-up"></i>
      </button>
      <button class="join-item btn btn-error" @click="emit('handleRemoveCommand', index)" type="button"><i class="fa-solid fa-trash"></i></button>
      <button class="join-item btn btn-neutral" @click="emit('handleMoveDown', index)" type="button" v-if="index !== amount - 1">
        <i class="fa-solid fa-arrow-down"></i>
      </button>
    </div>
  </div>
</template>
