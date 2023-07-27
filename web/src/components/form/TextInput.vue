<script setup lang="ts">
const props = defineProps<{ id: string; title: string; modelValue: string | number; help: string; errors?: any[]; validate?: string; class?: string }>();
const emit = defineEmits(['update:modelValue']);
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
        @keyup.esc="emit('update:modelValue', '')"
        class="input input-bordered w-full"
      />
    </div>
    <label class="label">
      <span class="label-text-alt select-text">
        <span>{{ help }}</span>
        <br />
        <span v-if="validate" class="text-error">{{ validate }}</span>
        <br v-if="validate" />
        <span v-for="error in errors" :key="error.$uid" class="text-error">{{ error.$message }}<br /></span>
      </span>
    </label>
  </div>
</template>
