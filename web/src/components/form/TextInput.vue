<script setup lang="ts">
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
</script>

<template>
  <div class="form-control w-full" :class="class" v-if="v$">
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
        :class="{ 'input-error': v$.$errors.length !== 0 || validate, 'input-warning': v$.$dirty }"
      />
    </div>
    <label class="label">
      <span class="label-text-alt select-text">
        <span v-if="v$.$errors.length === 0 && !validate">{{ help }}<br /></span>
        <span v-if="validate" class="text-error">{{ validate }}<br /></span>
        <span v-for="error in v$.$errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
