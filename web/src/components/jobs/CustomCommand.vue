<script setup lang="ts">
import { ref } from 'vue';
import TextInput from '@/components/form/TextInput.vue';

const emit = defineEmits(['gotRef', 'start']);
const cmdButton = ref();
const command = ref('');
</script>

<template>
  <dialog :ref="(el) => emit('gotRef', el)" id="my_modal" class="modal modal-bottom sm:modal-middle">
    <form method="dialog" class="modal-box">
      <h3 class="font-bold text-lg">Custom restic command</h3>
      <div class="py-4">
        <TextInput id="custom_command" title="Command" v-model="command" help="Example: restic snapshots" @keyup.enter="cmdButton.click()"> </TextInput>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button ref="cmdButton" @click="emit('start', 'custom', command)" class="btn btn-primary"><i class="fa-solid fa-check"></i>Submit</button>
          <button class="btn btn-neutral"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </div>
    </form>
  </dialog>
</template>
