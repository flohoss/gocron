<script setup lang="ts">
defineProps<{ error: string }>();
const emit = defineEmits(['gotRef']);
</script>

<template>
  <dialog :ref="(el) => emit('gotRef', el)" id="my_modal_5" class="modal modal-bottom sm:modal-middle">
    <form method="dialog" class="modal-box">
      <h3 class="font-bold text-lg">Error!</h3>
      <p class="py-4 whitespace-pre-line">{{ error }}</p>
      <div class="modal-action">
        <button class="btn">Close</button>
      </div>
    </form>
  </dialog>
</template>
