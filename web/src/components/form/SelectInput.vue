<script setup lang="ts">
import type { SelectOption } from '@/types';
import { inject } from 'vue';

const props = defineProps<{ id: string; title: string; modelValue: number; options: SelectOption[]; help: string; v$?: any; class?: string }>();
const emit = defineEmits(['update:modelValue']);
const submitted = inject('submitted', false);
</script>

<template>
  <div class="form-control w-full" :class="class">
    <label class="label">
      <span class="label-text">{{ title }}</span>
    </label>
    <select
      :id="id"
      class="select select-bordered w-full"
      :value="props.modelValue"
      @input="emit('update:modelValue', parseInt(($event.target as HTMLInputElement)?.value))"
      :class="{ 'select-error': v$ && v$.$errors.length !== 0, 'select-success': !submitted && v$ && v$.$dirty }"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">{{ option.description }}</option>
    </select>
    <label class="label">
      <span class="label-text-alt select-text">
        <span>{{ help }}<br /></span>
        <span v-if="v$" v-for="error in v$.$errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
