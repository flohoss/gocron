<script setup lang="ts">
import { CompressionTypesService, RetentionPoliciesService, type database_CompressionType, type database_Job, type database_SelectOption } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { watch } from 'vue';
import CommandsInput from '@/components/form/CommandsInput.vue';
import TextInput from '@/components/form/TextInput.vue';
import { useVuelidate } from '@vuelidate/core';
import { required, integer, helpers } from '@vuelidate/validators';
import SelectInput from '@/components/form/SelectInput.vue';

const route = useRoute();
const router = useRouter();

const store = useJobStore();
const storeJob = computed(() => store.getJob(route.params.id));
const job = ref<database_Job>({ ...storeJob.value });
watch(storeJob, () => (job.value = { ...storeJob.value }));

const svg = helpers.regex(/^<(svg|i).*<\/(svg|i)>$/);
const rules = {
  compression_type_id: { required },
  retention_policy_id: { required },
  description: { required },
  local_directory: { required },
  password_file_path: { required },
  restic_remote: { required },
  svg_icon: { svg: helpers.withMessage('This field should be a valid SVG', svg) },
  routine_check: { integer },
};

// @ts-ignore
const v$ = useVuelidate(rules, job);

const error = ref<string>('');
const errorModal = ref();

const handleSubmit = async () => {
  const isFormCorrect = await v$.value.$validate();
  if (!isFormCorrect) return;
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

const compressionTypes = ref<database_SelectOption[]>([]);
const retentionPolicies = ref<database_SelectOption[]>([]);
const init = async () => {
  compressionTypes.value = await CompressionTypesService.getCompressionTypes();
  retentionPolicies.value = await RetentionPoliciesService.getRetentionPolicies();
};
init();
</script>

<template>
  <div>
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="text-xl font-bold">{{ header }}</div>
    </PageHeader>
    <PageContent>
      <form class="grid gap-10" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <TextInput title="Description" v-model="job.description" help="Example: Gitea" :errors="v$.description.$errors" />
          <TextInput title="Local directory" v-model="job.local_directory" help="Example: /opt/docker/gitea" :errors="v$.local_directory.$errors" />
          <TextInput title="Restic Remote" v-model="job.restic_remote" help="Example: rclone:pcloud:Backups/gitea" :errors="v$.restic_remote.$errors" />
          <TextInput title="Password file" v-model="job.password_file_path" help="Example: /secrets/.restipw" :errors="v$.password_file_path.$errors" />
          <SelectInput title="Compression" v-model="job.compression_type_id" :errors="v$.compression_type_id.$errors" :options="compressionTypes" />
          <SelectInput title="Retention policy" v-model="job.retention_policy_id" :errors="v$.retention_policy_id.$errors" :options="retentionPolicies" />
          <TextInput
            v-if="job.svg_icon !== undefined"
            title="SVG-Icon"
            v-model="job.svg_icon"
            help="Example: <i class='fa-solid fa-circle-nodes'></i>"
            :errors="v$.svg_icon.$errors"
          />
          <TextInput
            v-if="job.routine_check !== undefined"
            title="Routine check"
            v-model="job.routine_check"
            help="Example: 15"
            :errors="v$.routine_check.$errors"
          />
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <CommandsInput v-if="job.post_commands !== undefined" title="Commands before backup" :commands="job.post_commands" />
          <CommandsInput v-if="job.pre_commands !== undefined" title="Commands after backup" :commands="job.pre_commands" />
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Submit</button>
          <button @click.prevent="router.push({ name: 'home' })" class="btn btn-neutral" type="submit"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </form>
    </PageContent>
  </div>
</template>
