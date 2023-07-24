<script setup lang="ts">
import type { database_PostCommand, database_PreCommand } from '@/openapi';

const props = defineProps<{ title: string; commands: database_PreCommand[] | database_PostCommand[] }>();

const handleAddCommand = () => {
  props.commands.push({ command: '' });
};

const handleRemoveCommand = (index: number) => {
  props.commands.splice(index, 1);
};

const handleMoveUp = (index: number) => {
  if (index > 0) {
    [props.commands[index - 1], props.commands[index]] = [props.commands[index], props.commands[index - 1]];
  }
};

const handleMoveDown = (index: number) => {
  if (index < props.commands.length - 1) {
    [props.commands[index], props.commands[index + 1]] = [props.commands[index + 1], props.commands[index]];
  }
};

const handleChange = (index: number, newValue: string) => {
  props.commands[index].command = newValue;
};
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="label-text">{{ title }}</div>
    <div v-for="(command, index) in commands" :key="index" class="join">
      <input class="join-item input input-bordered w-full" v-model="command.command" @input="handleChange(index, command.command)" />
      <button class="join-item btn btn-neutral" @click="handleMoveUp(index)" type="button" :disabled="index === 0"><i class="fa-solid fa-arrow-up"></i></button>
      <button class="join-item btn btn-neutral" @click="handleMoveDown(index)" type="button" :disabled="index === commands.length - 1">
        <i class="fa-solid fa-arrow-down"></i>
      </button>
      <button class="join-item btn btn-error" @click="handleRemoveCommand(index)" type="button"><i class="fa-solid fa-trash"></i></button>
    </div>
    <button class="btn btn-sm" @click="handleAddCommand" type="button">Add Command</button>
  </div>
</template>
