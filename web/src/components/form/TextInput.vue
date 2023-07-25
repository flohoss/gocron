<script setup lang="ts">
const props = defineProps<{ id: string; title: string; modelValue: string | number; help: string; errors?: any[]; class?: string }>();
const emit = defineEmits(['update:modelValue']);
</script>

<template>
  <div class="form-control w-full" :class="class">
    <label class="label">
      <span class="label-text">{{ title }}</span>
    </label>
    <input
      :id="id"
      type="text"
      :value="props.modelValue"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement)?.value)"
      @keyup.esc="emit('update:modelValue', '')"
      class="input input-bordered w-full"
    />
    <label class="label">
      <span class="label-text-alt select-text">
        <span>{{ help }}</span>
        <br />
        <span v-for="error in errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
