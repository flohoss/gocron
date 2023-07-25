<script setup lang="ts">
import { CompressionTypesService, RetentionPoliciesService, type database_Job, type database_Command, type database_SelectOption } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { watch } from 'vue';
import CommandInput from '@/components/form/CommandInput.vue';
import TextInput from '@/components/form/TextInput.vue';
import { useVuelidate } from '@vuelidate/core';
import { required, integer, helpers, between } from '@vuelidate/validators';
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
  routine_check: { required, integer, between: between(0, 100) },
  pre_commands: {
    $each: helpers.forEach({
      command: { required },
    }),
  },
  post_commands: {
    $each: helpers.forEach({
      command: { required },
    }),
  },
};

// @ts-ignore
const v$ = useVuelidate(rules, job);

const error = ref<string>('');
const errorModal = ref();

const handleSubmit = async () => {
  const isFormCorrect = await v$.value.$validate();
  if (!isFormCorrect) return;

  // exclude runs from update
  job.value.runs = [];
  job.value.routine_check = parseInt(job.value.routine_check + '');
  try {
    if (job.value.id === 0) {
      const created = await store.createJob(job.value);
      v$.value.$reset();
      router.push({ name: 'jobs', params: { id: created.id } });
    } else {
      await store.updateJob(job.value);
      v$.value.$reset();
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

const handleAddCommand = (type: number, commands: database_Command[] | undefined) => {
  commands && commands.push({ command: '', sort_id: commands.length + 1, type: type });
};

const handleRemoveCommand = async (index: number, commands: database_Command[] | undefined) => {
  commands && commands.splice(index, 1);
  setSortIds(commands);
};

const handleMoveUp = (index: number, commands: database_Command[] | undefined) => {
  commands && commands.splice(index - 1, 0, commands.splice(index, 1)[0]);
  setSortIds(commands);
};

const handleMoveDown = (index: number, commands: database_Command[] | undefined) => {
  commands && commands.splice(index + 1, 0, commands.splice(index, 1)[0]);
  setSortIds(commands);
};

const setSortIds = (commands: database_Command[] | undefined) => {
  if (commands) {
    for (let i = 0; i < commands.length; i++) {
      commands[i].sort_id = i + 1;
    }
  }
};
</script>

<template>
  <div>
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="text-xl font-bold">{{ header }}</div>
    </PageHeader>
    <PageContent>
      <form class="grid gap-10" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-x-5">
          <TextInput id="description" title="Description" v-model="job.description" help="Example: Gitea" :errors="v$.description.$errors" />
          <TextInput
            id="local_directory"
            title="Local directory"
            v-model="job.local_directory"
            help="Example: /opt/docker/gitea"
            :errors="v$.local_directory.$errors"
          />
          <TextInput
            id="restic_remote"
            title="Restic Remote"
            v-model="job.restic_remote"
            help="Example: rclone:pcloud:Backups/gitea"
            :errors="v$.restic_remote.$errors"
          />
          <TextInput
            id="password_file_path"
            title="Password file"
            v-model="job.password_file_path"
            help="Example: /secrets/.restipw"
            :errors="v$.password_file_path.$errors"
          />
          <SelectInput
            id="compression_type_id"
            title="Compression"
            v-model="job.compression_type_id"
            help="How data is compressed"
            :errors="v$.compression_type_id.$errors"
            :options="compressionTypes"
          />
          <SelectInput
            id="retention_policy_id"
            title="Retention policy"
            v-model="job.retention_policy_id"
            help="Policy for which snapshots to keep"
            :errors="v$.retention_policy_id.$errors"
            :options="retentionPolicies"
          />
          <TextInput
            v-if="job.routine_check !== undefined"
            id="routine_check"
            title="Routine check"
            v-model="job.routine_check"
            help="Range: 0-100 (0: disabled)"
            :errors="v$.routine_check.$errors"
          />
          <TextInput
            v-if="job.svg_icon !== undefined"
            id="svg_icon"
            title="SVG-Icon"
            v-model="job.svg_icon"
            help="Example: <i class='fa-solid fa-circle-nodes'></i>"
            :errors="v$.svg_icon.$errors"
          />
          <div class="grid gap-x-5 grid-cols-1 lg:grid-cols-2 mt-5 gap-y-5 col-span-1 lg:col-span-2 xl:col-span-4">
            <div v-if="job.pre_commands !== undefined">
              <div v-for="(command, index) in job.pre_commands" :key="command.id">
                <CommandInput
                  id="pre_commands"
                  v-model="command.command"
                  :index="index"
                  :amount="job.pre_commands.length"
                  @handleRemoveCommand="(index) => handleRemoveCommand(index, job.pre_commands)"
                  @handleMoveUp="(index) => handleMoveUp(index, job.pre_commands)"
                  @handleMoveDown="(index) => handleMoveDown(index, job.pre_commands)"
                  :errors="v$.pre_commands.$each.$response.$errors[index].command"
                  >{{ index + 1 }}. Command</CommandInput
                >
              </div>
              <button type="button" class="btn btn-sm btn-neutral" @click="handleAddCommand(1, job.pre_commands)">
                <i class="fa-solid fa-plus"></i>Add Command before backup
              </button>
            </div>
            <div v-if="job.post_commands !== undefined">
              <div v-for="(command, index) in job.post_commands" :key="command.id">
                <CommandInput
                  id="post_commands"
                  v-model="command.command"
                  :index="index"
                  :amount="job.post_commands.length"
                  @handleRemoveCommand="(index) => handleRemoveCommand(index, job.post_commands)"
                  @handleMoveUp="(index) => handleMoveUp(index, job.post_commands)"
                  @handleMoveDown="(index) => handleMoveDown(index, job.post_commands)"
                  :errors="v$.post_commands.$each.$response.$errors[index].command"
                  >{{ index + 1 }}. Command</CommandInput
                >
              </div>
              <button type="button" class="btn btn-sm btn-neutral" @click="handleAddCommand(2, job.post_commands)">
                <i class="fa-solid fa-plus"></i>Add Command after backup
              </button>
            </div>
          </div>
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Submit</button>
          <button @click.prevent="router.push({ name: 'home' })" class="btn btn-neutral" type="submit"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </form>
    </PageContent>
  </div>
</template>
