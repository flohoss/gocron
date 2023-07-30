<script setup lang="ts">
import type { database_Job } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { useConfirmDialog } from '@vueuse/core';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

const store = useJobStore();
const router = useRouter();
const props = defineProps<{ job: database_Job }>();
const emit = defineEmits(['start', 'showModal']);

const disabled = computed(() => {
  if (props.job.runs && props.job.runs?.length !== 0) {
    return !props.job.runs[0].end_time;
  }
  return false;
});

const { reveal, confirm, cancel, onReveal, onConfirm, onCancel } = useConfirmDialog();
const confirmModal = ref();
onReveal(() => confirmModal.value.showModal());
onCancel(() => confirmModal.value.close());
onConfirm(() => {
  if (props.job.id) {
    store
      .deleteJob(props.job.id)
      .then(() => router.push({ name: 'home' }))
      .catch((err) => console.log(err));
  }
});
</script>

<template>
  <div class="flex justify-between items-center min-w-0">
    <div class="text-xl font-bold truncate">{{ job.description }}</div>
    <div class="flex-shrink-0 flex flex-col items-end gap-1">
      <div class="join flex-shrink-0 flex-wrap">
        <button @click="emit('start', 'prune')" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-broom"></i><span class="hidden lg:block">Prune</span>
        </button>
        <button @click="emit('start', 'check')" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-check"></i><span class="hidden lg:block">Check</span>
        </button>
        <button @click="emit('showModal')" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-terminal"></i><span class="hidden lg:block">Custom</span>
        </button>
        <button @click="emit('start', 'run')" class="join-item btn btn-sm btn-success" :disabled="disabled">
          <i class="fa-solid fa-play"></i><span class="hidden xl:block">Run</span>
        </button>
        <button @click="router.push({ name: 'jobsForm', params: { id: job.id } })" class="join-item btn btn-sm btn-warning" :disabled="disabled">
          <i class="fa-solid fa-pencil"></i><span class="hidden xl:block">Edit</span>
        </button>
        <button @click="reveal" class="join-item btn btn-sm btn-error" :disabled="disabled">
          <i class="fa-solid fa-trash"></i><span class="hidden xl:block">Delete</span>
        </button>
      </div>
    </div>

    <teleport to="body">
      <dialog ref="confirmModal" id="delete_modal" class="modal modal-bottom sm:modal-middle">
        <form method="dialog" class="modal-box">
          <p class="py-4">Do you want to delete {{ job.description }}?</p>
          <div class="modal-action">
            <button type="button" @click="cancel" class="btn btn-error">Cancel</button>
            <button type="button" @click="confirm" class="btn btn-success">Yes</button>
          </div>
        </form>
      </dialog>
    </teleport>
  </div>
</template>
