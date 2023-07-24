<script setup lang="ts">
const props = defineProps<{ modelValue: string; index: number; amount: number; errors?: any[] }>();
const emit = defineEmits(['update:modelValue', 'handleMoveUp', 'handleMoveDown', 'handleRemoveCommand']);
</script>

<template>
  <label class="label">
    <span class="label-text"><slot></slot></span>
  </label>
  <div class="join w-full">
    <input class="join-item input input-bordered w-full" :value="modelValue" @input="emit('update:modelValue', ($event.target as HTMLInputElement)?.value)" />
    <button class="join-item btn btn-neutral" @click="emit('handleMoveUp', index)" type="button" :disabled="index === 0">
      <i class="fa-solid fa-arrow-up"></i>
    </button>
    <button class="join-item btn btn-neutral" @click="emit('handleMoveDown', index)" type="button" :disabled="index === amount - 1">
      <i class="fa-solid fa-arrow-down"></i>
    </button>
    <button class="join-item btn btn-error" @click="emit('handleRemoveCommand', index)" type="button"><i class="fa-solid fa-trash"></i></button>
  </div>
  <label v-if="errors" class="label">
    <span class="label-text-alt select-text">
      <span v-for="error in errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
    </span>
  </label>
</template>
