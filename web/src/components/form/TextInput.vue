<script setup lang="ts">
const props = defineProps<{ title: string; modelValue: string; help: string; errors?: any[] }>();
const emit = defineEmits(['update:modelValue']);
</script>

<template>
  <div class="form-control w-full">
    <label class="label">
      <span class="label-text">{{ title }}</span>
    </label>
    <input
      type="text"
      :value="props.modelValue"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement)?.value)"
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
