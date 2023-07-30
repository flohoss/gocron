<script setup lang="ts">
defineProps<{ id: string; command: string; fileOutput: string; index: number; amount: number; errors?: any[] }>();
const emit = defineEmits(['update:command', 'update:fileOutput', 'handleMoveUp', 'handleMoveDown', 'handleRemoveCommand']);
</script>

<template>
  <div class="flex gap-x-5 flex-col lg:flex-row items-end lg:items-center">
    <div class="form-control w-full lg:w-4/6">
      <label class="label">
        <span class="label-text"><slot></slot></span>
      </label>
      <input
        :id="id + '_cmd_' + index"
        type="text"
        class="input input-bordered w-full"
        :value="command"
        @input="emit('update:command', ($event.target as HTMLInputElement)?.value)"
      />
      <label class="label">
        <span class="label-text-alt select-text">
          <span v-if="errors?.length == 0">Example: docker compose up</span>
          <br v-if="errors?.length == 0" />
          <span v-for="error in errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
        </span>
      </label>
    </div>
    <div class="form-control w-full lg:w-2/6">
      <label class="label">
        <span class="label-text">File output</span>
      </label>
      <input
        :id="id + '_file_' + index"
        type="text"
        class="input input-bordered w-full"
        :value="fileOutput"
        @input="emit('update:fileOutput', ($event.target as HTMLInputElement)?.value)"
      />
      <label class="label">
        <span class="label-text-alt">Output will be redirected to file</span>
      </label>
    </div>
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
