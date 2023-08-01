<script setup lang="ts">
import { type database_Job, type database_Command } from '@/openapi';
import { useJobStore } from '@/stores/jobs';
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import PageHeader from '@/components/ui/PageHeader.vue';
import PageContent from '@/components/ui/PageContent.vue';
import { watch } from 'vue';
import CommandInput from '@/components/form/CommandInput.vue';
import TextInput from '@/components/form/TextInput.vue';
import { useVuelidate } from '@vuelidate/core';
import { required, integer, helpers, between } from '@vuelidate/validators';
import SelectInput from '@/components/form/SelectInput.vue';
import { CompressionOptions, RetentionPolicyOptions } from '@/types';

const route = useRoute();
const router = useRouter();

const store = useJobStore();
const storeJob = computed(() => store.getJob(route.params.id));
const job = ref<database_Job>({ ...storeJob.value });
watch(storeJob, () => (job.value = { ...storeJob.value }));

const svg = helpers.regex(/^<(svg|i).*<\/(svg|i)>$/);
const rules = {
  compression_type: { required },
  retention_policy: { required },
  description: { required },
  local_directory: { required },
  password_file_path: { required },
  restic_remote: { required },
  svg_icon: { svg: helpers.withMessage('This field should be a valid SVG', svg) },
  routine_check: { required, integer, between: between(0, 100) },
};

// @ts-ignore
const v$ = useVuelidate(rules, job);
const validate = ref({
  Description: '',
  LocalDirectory: '',
  ResticRemote: '',
  PasswordFilePath: '',
  SvgIcon: '',
  RoutineCheck: '',
});

const handleSubmit = async () => {
  if (v$.value.$errors.length !== 0) return;

  job.value.routine_check = parseInt(job.value.routine_check + '');
  if (job.value.id === 0) {
    store
      .createJob(job.value)
      .then((res) => {
        v$.value.$reset();
        router.push({ name: 'jobs', params: { id: res.id } });
      })
      .catch((err) => (validate.value = err.body));
  } else {
    store
      .updateJob(job.value)
      .then(() => {
        v$.value.$reset();
        router.push({ name: 'jobs', params: { id: job.value.id } });
      })
      .catch((err) => (validate.value = err.body));
  }
};

const header = computed(() => (job.value.id !== 0 ? 'Edit' : 'New') + ' Job');

const handleAddCommand = (type: number, commands: database_Command[] | undefined) => {
  commands && commands.push({ id: 0, job_id: 0, command: '', sort_id: commands.length + 1, type: type, file_output: '' });
};

const handleRemoveCommand = (index: number, commands: database_Command[] | undefined) => {
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
    <PageHeader>
      <div class="text-xl font-bold">{{ header }}</div>
    </PageHeader>
    <PageContent>
      <form class="grid gap-10" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-x-5">
          <TextInput id="description" title="Description" v-model="v$.description.$model" help="Example: Gitea" :v$="v$.description" />
          <TextInput
            id="local_directory"
            title="Local directory"
            v-model="v$.local_directory.$model"
            help="Example: /opt/docker/gitea"
            :v$="v$.local_directory"
            :validate="validate.LocalDirectory"
          />
          <TextInput
            id="restic_remote"
            title="Restic Remote"
            v-model="v$.restic_remote.$model"
            help="Example: rclone:pcloud:Backups/gitea"
            :v$="v$.restic_remote"
            :validate="validate.ResticRemote"
          />
          <TextInput
            id="password_file_path"
            title="Password file"
            v-model="v$.password_file_path.$model"
            help="Example: /secrets/.restipw"
            :v$="v$.password_file_path"
            :validate="validate.PasswordFilePath"
          />
          <SelectInput
            id="compression_type"
            title="Compression"
            v-model="v$.compression_type.$model"
            help="How data is compressed"
            :v$="v$.compression_type"
            :options="CompressionOptions"
          />
          <TextInput
            v-if="job.routine_check !== undefined"
            id="routine_check"
            title="Routine check"
            v-model="v$.routine_check.$model"
            help="Range: 0-100 (0: disabled)"
            :v$="v$.routine_check"
            :validate="validate.RoutineCheck"
          />
          <SelectInput
            id="retention_policy"
            class="col-span-1 lg:col-span-2"
            title="Retention policy"
            v-model="v$.retention_policy.$model"
            help="Policy for which snapshots to keep"
            :v$="v$.retention_policy"
            :options="RetentionPolicyOptions"
          />
          <div class="grid gap-x-5 grid-cols-1 mt-5 gap-y-5 col-span-1 lg:col-span-2">
            <div v-if="job.pre_commands !== undefined">
              <div v-for="(command, index) in job.pre_commands" :key="command.id">
                <CommandInput
                  v-if="command.file_output !== undefined"
                  id="pre_commands"
                  v-model:command="command.command"
                  v-model:fileOutput="command.file_output"
                  :index="index"
                  :amount="job.pre_commands.length"
                  @handleRemoveCommand="(index) => handleRemoveCommand(index, job.pre_commands)"
                  @handleMoveUp="(index) => handleMoveUp(index, job.pre_commands)"
                  @handleMoveDown="(index) => handleMoveDown(index, job.pre_commands)"
                />
              </div>
              <button type="button" class="btn btn-sm btn-neutral mt-2" @click="handleAddCommand(1, job.pre_commands)">
                <i class="fa-solid fa-plus"></i>Add Command before backup
              </button>
            </div>
            <div v-if="job.post_commands !== undefined">
              <div v-for="(command, index) in job.post_commands" :key="command.id">
                <CommandInput
                  v-if="command.file_output !== undefined"
                  id="post_commands"
                  v-model:command="command.command"
                  v-model:fileOutput="command.file_output"
                  :index="index"
                  :amount="job.post_commands.length"
                  @handleRemoveCommand="(index) => handleRemoveCommand(index, job.post_commands)"
                  @handleMoveUp="(index) => handleMoveUp(index, job.post_commands)"
                  @handleMoveDown="(index) => handleMoveDown(index, job.post_commands)"
                />
              </div>
              <button type="button" class="btn btn-sm btn-neutral mt-2" @click="handleAddCommand(2, job.post_commands)">
                <i class="fa-solid fa-plus"></i>Add Command after backup
              </button>
            </div>
          </div>
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Submit</button>
          <button @click.prevent="router.push({ name: 'home' })" class="btn btn-neutral" type="button"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </form>
    </PageContent>
  </div>
</template>
