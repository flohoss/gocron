<script setup lang="ts">
import { type database_Job } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { watch } from 'vue';
import CommandsInput from '@/components/form/CommandsInput.vue';
import TextInput from '@/components/form/TextInput.vue';

const route = useRoute();
const router = useRouter();

const store = useJobStore();
const storeJob = computed(() => store.getJob(route.params.id));
const job = ref<database_Job>({ ...storeJob.value });

watch(storeJob, () => (job.value = { ...storeJob.value }));

const error = ref<string>('');
const errorModal = ref();

const handleSubmit = async () => {
  try {
    if (job.value.id === 0) {
      const created = await store.createJob(job.value);
      router.push({ name: 'jobs', params: { id: created.id } });
    } else {
      await store.updateJob(job.value);
      router.push({ name: 'jobs', params: { id: job.value.id } });
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
      <form class="grid gap-5" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <TextInput title="Description" v-model="job.description" help="Example: Gitea" />
          <TextInput title="Local directory" v-model="job.local_directory" help="Example: /opt/docker/gitea" />
          <TextInput title="Restic Remote" v-model="job.restic_remote" help="Example: rclone:pcloud:Backups/gitea" />
          <TextInput title="Password file" v-model="job.password_file_path" help="Example: /secrets/.restipw" />
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <TextInput v-if="job.svg_icon !== undefined" title="SVG-Icon" v-model="job.svg_icon" help="Example: <i class='fa-solid fa-circle-nodes'></i>" />
          <TextInput v-if="job.routine_check !== undefined" title="Routine check" v-model="job.routine_check" help="Example: 15" />
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <CommandsInput title="Commands before backup" :commands="job.post_commands ? job.post_commands : []" />
          <CommandsInput title="Commands after backup" :commands="job.pre_commands ? job.pre_commands : []" />
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Submit</button>
          <button @click.prevent="router.push({ name: 'home' })" class="btn btn-neutral" type="submit"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </form>
    </PageContent>
  </div>
</template>
