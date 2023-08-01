<script setup lang="ts">
import { inject } from 'vue';

const props = defineProps<{
  id: string;
  title: string;
  modelValue: string | number;
  help: string;
  v$?: any;
  validate?: string;
  class?: string;
}>();
const emit = defineEmits(['update:modelValue']);
const submitted = inject('submitted', false);
</script>

<template>
  <div class="form-control w-full" :class="class">
    <label class="label">
      <span class="label-text">{{ title }}</span>
    </label>
    <div>
      <input
        :id="id"
        type="text"
        :value="props.modelValue"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement)?.value)"
        class="input input-bordered w-full"
        :class="{ 'input-error': (v$ && v$.$errors.length !== 0) || validate, 'input-success': !submitted && v$ && v$.$dirty }"
      />
    </div>
    <label class="label">
      <span class="label-text-alt select-text">
        <span v-if="(!v$ || v$.$errors.length === 0) && !validate">{{ help }}<br /></span>
        <span v-if="validate" class="text-error">{{ validate }}<br /></span>
        <span v-if="v$" v-for="error in v$.$errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
