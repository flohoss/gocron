<script setup lang="ts">
import { database_LogSeverity, type database_Job, JobsService } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { useConfirmDialog } from '@vueuse/core';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import BadgeList from '../ui/BadgeList.vue';
import { CompressionOptions, RetentionPolicyOptions } from '@/types';

const store = useJobStore();
const router = useRouter();
const props = defineProps<{ job: database_Job }>();
const emit = defineEmits(['start', 'showModal', 'clearRuns']);

const disabled = computed(() => {
  for (const job of store.jobs) {
    if (job.status === database_LogSeverity.LogNone) return true;
  }
  return false;
});

const {
  reveal: revealDelete,
  confirm: confirmDelete,
  cancel: cancelDelete,
  onReveal: onDeleteReveal,
  onConfirm: onDeleteConfirm,
  onCancel: onDeleteCancel,
} = useConfirmDialog();
const deleteModal = ref();
onDeleteReveal(() => deleteModal.value.showModal());
onDeleteCancel(() => deleteModal.value.close());
onDeleteConfirm(() => {
  if (props.job.id) {
    store
      .deleteJob(props.job.id)
      .then(() => router.push({ name: 'home' }))
      .catch((err) => console.log(err));
  }
});

const {
  reveal: revealClearRuns,
  confirm: confirmClearRuns,
  cancel: cancelClearRuns,
  onReveal: onClearRunsReveal,
  onConfirm: onClearRunsConfirm,
  onCancel: onClearRunsCancel,
} = useConfirmDialog();
const clearRunsModal = ref();
onClearRunsReveal(() => clearRunsModal.value.showModal());
onClearRunsCancel(() => clearRunsModal.value.close());
onClearRunsConfirm(() => {
  if (props.job.id) {
    JobsService.deleteJobsRuns(props.job.id)
      .then(() => emit('clearRuns'))
      .catch((err) => console.log(err))
      .finally(() => clearRunsModal.value.close());
  }
});

const re = new RegExp(/((?:--password|PASSWORD|-p)[="'\s]+)(.+?)(["'\s])/g);
function anonymisePasswords(text: string): string {
  return (text + ' ').replace(re, (match, p1, p2, p3) => `${p1}****${p3}`).trim();
}

const badges = computed(() => {
  const list: { [key: string]: any } = {
    'local directory': props.job.local_directory,
    'restic remote': props.job.restic_remote,
    compression: CompressionOptions[props.job.compression_type - 1].description,
    'retention policy': RetentionPolicyOptions[props.job.retention_policy - 1].description,
    'routine check': props.job.routine_check + '%',
  };
  props.job.pre_commands?.forEach((value, index) => {
    index++;
    list[index + '. pre command'] = anonymisePasswords(value.command);
  });
  props.job.post_commands?.forEach((value, index) => {
    index++;
    list[index + '. post command'] = anonymisePasswords(value.command);
  });
  return list;
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
        <button @click="revealClearRuns()" class="join-item btn btn-sm btn-neutral" :disabled="disabled">
          <i class="fa-solid fa-eraser"></i><span class="hidden lg:block">Clear runs</span>
        </button>
        <button @click="emit('start', 'run')" class="join-item btn btn-sm btn-success" :disabled="disabled">
          <i class="fa-solid fa-play"></i><span class="hidden xl:block">Run</span>
        </button>
        <button @click="router.push({ name: 'jobsForm', params: { id: job.id } })" class="join-item btn btn-sm btn-warning" :disabled="disabled">
          <i class="fa-solid fa-pencil"></i><span class="hidden xl:block">Edit</span>
        </button>
        <button @click="revealDelete()" class="join-item btn btn-sm btn-error" :disabled="disabled">
          <i class="fa-solid fa-trash"></i><span class="hidden xl:block">Delete</span>
        </button>
      </div>
    </div>
  </div>

  <div class="flex gap-2 mt-5 flex-wrap select-none">
    <badge-list :badges="badges" />
  </div>

  <teleport to="body">
    <dialog ref="deleteModal" id="delete_modal" class="modal modal-bottom sm:modal-middle">
      <form method="dialog" class="modal-box bg-base-300 text-base-content">
        <p class="py-4">Are you sure you want to delete {{ job.description }} ?</p>
        <div class="modal-action">
          <button type="button" @click="cancelDelete" class="btn btn-error">No, wait!</button>
          <button type="button" @click="confirmDelete" class="btn btn-success">Sure!</button>
        </div>
      </form>
    </dialog>
    <dialog ref="clearRunsModal" id="clear_runs_modal" class="modal modal-bottom sm:modal-middle">
      <form method="dialog" class="modal-box bg-base-300 text-base-content">
        <p class="py-4">Are you sure you want to clear all runs of {{ job.description }} ?</p>
        <div class="modal-action">
          <button type="button" @click="cancelClearRuns" class="btn btn-error">No, wait!</button>
          <button type="button" @click="confirmClearRuns" class="btn btn-success">Sure!</button>
        </div>
      </form>
    </dialog>
  </teleport>
</template>
