<script setup lang="ts">
import type { database_SelectOption } from '@/openapi';

const props = defineProps<{ title: string; modelValue: number; options: database_SelectOption[]; errors?: any[] }>();
const emit = defineEmits(['update:modelValue']);
</script>

<template>
  <div class="form-control w-full">
    <label class="label">
      <span class="label-text">{{ title }}</span>
    </label>
    <select
      class="select select-bordered w-full"
      :value.number="props.modelValue"
      @input="emit('update:modelValue', parseInt(($event.target as HTMLInputElement)?.value))"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">{{ option.name }}</option>
    </select>
    <label class="label">
      <span class="label-text-alt select-text">
        <span v-for="error in errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
