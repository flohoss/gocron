<script setup lang="ts">
import { type database_Job } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { watch } from 'vue';

const route = useRoute();
const router = useRouter();

const store = useJobStore();
const storeJob = computed(() => store.getJob(route.params.id));
const job = ref<database_Job>({ ...storeJob.value });

watch(storeJob, () => (job.value = { ...storeJob.value }));

const error = ref<string>('');
const errorModal = ref();

const handleSubmit = async (job: database_Job) => {
  try {
    if (job.id === 0) {
      const created = await store.createJob(job);
      router.push({ name: 'jobs', params: { id: created.id } });
    } else {
      await store.updateJob(job);
      router.push({ name: 'jobs', params: { id: job.id } });
    }
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};

const header = computed(() => (job.value.id !== 0 ? 'Edit' : 'New') + ' Job');
</script>

<template>
  <div>
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="text-xl font-bold">{{ header }}</div>
    </PageHeader>
    <PageContent>
      <FormKit
        v-model="job"
        type="form"
        form-class="grid grid-cols-1 lg:grid-cols-2 gap-5"
        :actions="false"
        :config="{ validationVisibility: 'submit' }"
        @submit="handleSubmit"
      >
        <FormKit type="text" name="description" id="description" validation="required" label="Description" help="Example: Gitea" />
        <FormKit type="text" name="local_directory" id="local_directory" validation="required" label="Local directory" help="Example: /opt/docker/gitea" />
        <FormKit type="text" name="restic_remote" id="restic_remote" validation="required" label="Restic Remote" help="Example: rclone:pcloud:Backups/gitea" />
        <FormKit type="text" name="password_file_path" id="password_file_path" validation="required" label="Password file" help="Example: /secrets/.restipw" />
        <FormKit type="text" name="svg_icon" id="svg_icon" validation="" label="SVG-Icon" help="Example: <i class='fa-solid fa-circle-nodes'></i>" />
        <FormKit type="text" name="routine_check" id="routine_check" validation="number|between:1,100" label="Routine check" suffix=" %" help="Example: 15" />

        <div class="flex justify-start flex-row-reverse gap-5 lg:col-span-2">
          <FormKit type="submit"><i class="fa-solid fa-check"></i>Submit</FormKit>
          <FormKit type="button" @click="router.push({ name: 'home' })"><i class="fa-solid fa-times"></i>Cancel</FormKit>
        </div>
      </FormKit>
    </PageContent>
  </div>
</template>
