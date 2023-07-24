<script setup lang="ts">
const props = defineProps<{ title: string; modelValue: string; help: string }>();
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
      <span class="label-text-alt">{{ help }}</span>
    </label>
  </div>
</template>
