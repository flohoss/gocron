<script setup lang="ts">
import { ref } from 'vue';
import TextInput from '@/components/form/TextInput.vue';

defineProps<{ error: string }>();
const emit = defineEmits(['gotRef', 'start']);
const command = ref('');
const submitButton = ref();
</script>

<template>
  <dialog :ref="(el) => emit('gotRef', el)" id="my_modal" class="modal modal-bottom sm:modal-middle">
    <form method="dialog" class="modal-box bg-base-300 text-base-content">
      <h3 class="font-bold text-lg">Custom restic command</h3>
      <div class="py-4">
        <TextInput @keydown.enter.prevent="submitButton.click()" id="custom_command" title="Command" v-model="command" help="Example: restic snapshots">
        </TextInput>
        <span class="text-sm text-error">{{ error }}</span>
        <div class="flex justify-start flex-row-reverse gap-5">
          <div ref="submitButton" @click="emit('start', 'custom', command)" type="submit" class="btn btn-primary"><i class="fa-solid fa-check"></i>Submit</div>
          <button class="btn btn-neutral"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </div>
    </form>
  </dialog>
</template>
