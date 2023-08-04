<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import useVuelidate from '@vuelidate/core';
import { required } from '@vuelidate/validators';
import { onBeforeUnmount, reactive, ref, watch } from 'vue';
import TextInput from '@/components/form/TextInput.vue';
import { CommandsService, SystemService, type database_SystemLog } from '@/openapi';
import { useEventSource } from '@vueuse/core';
import TerminalLog from '@/components/ui/TerminalLog.vue';

const state = reactive({
  restic_remote: '',
  local_directory: '',
  password_file_path: '',
});

const rules = {
  restic_remote: { required },
  password_file_path: { required },
  local_directory: {},
};

const v$ = useVuelidate(rules, state);
const validate = ref({ ResticRemote: '', PasswordFilePath: '', LocalDirectory: '' });

const handleSubmit = async () => {
  const isFormCorrect = await v$.value.$validate();
  if (!isFormCorrect) return;

  validate.value = { ResticRemote: '', PasswordFilePath: '', LocalDirectory: '' };

  CommandsService.postCommands({
    command: 'restore',
    restic_remote: state.restic_remote,
    local_directory: state.local_directory,
    password_file_path: state.password_file_path,
  })
    .then(() => {
      v$.value.$reset();
    })
    .catch((err) => (validate.value = err.body));
};

const logs = ref<database_SystemLog[]>([]);

const init = () => {
  SystemService.getSystemLogs()
    .then((res) => (logs.value = res))
    .catch((err) => console.log(err));
};
init();

const { data, close } = useEventSource('/api/sse?stream=restore_logs');
watch(data, (value) => {
  const parsed: database_SystemLog = value && JSON.parse(value);
  if (logs.value.length > 5) {
    logs.value = logs.value.slice(1);
  }
  logs.value.push(parsed);
});
onBeforeUnmount(() => close());
</script>

<template>
  <div>
    <PageHeader>
      <div class="text-xl font-bold">Restore</div>
    </PageHeader>
    <PageContent>
      <form class="grid gap-10" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 xl:grid-cols-3 gap-x-5">
          <TextInput
            id="restic_remote"
            title="Restic Remote"
            v-model="v$.restic_remote.$model"
            help="Example: rclone:pcloud:Backups/gitea"
            :v$="v$.restic_remote"
            :validate="validate.ResticRemote"
          />
          <TextInput id="local_directory" title="Local directory" v-model="v$.local_directory.$model" help="Default: /" :validate="validate.LocalDirectory" />
          <TextInput
            id="password_file_path"
            title="Password file"
            v-model="v$.password_file_path.$model"
            help="Example: /secrets/.restipw"
            :v$="v$.password_file_path"
            :validate="validate.PasswordFilePath"
          />
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Restore</button>
        </div>
      </form>
      <div class="grid grid-cols-1 gap-5 overflow-x-auto mt-10">
        <TerminalLog :logs="logs" />
      </div>
    </PageContent>
  </div>
</template>
