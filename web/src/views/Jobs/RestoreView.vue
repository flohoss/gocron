<script setup lang="ts">
import PageContent from '@/components/ui/PageContent.vue';
import PageHeader from '@/components/ui/PageHeader.vue';
import useVuelidate from '@vuelidate/core';
import { required } from '@vuelidate/validators';
import { reactive, ref } from 'vue';
import TextInput from '@/components/form/TextInput.vue';
import { useRouter } from 'vue-router';
import ErrorModal from '@/components/ui/ErrorModal.vue';
import moment from 'moment';

const router = useRouter();
const state = reactive({
  restic_remote: '',
  local_directory: '',
  password_file_path: '',
});

const rules = {
  restic_remote: { required },
  password_file_path: { required },
};

const v$ = useVuelidate(rules, state);

const error = ref<string>('');
const errorModal = ref();

const handleSubmit = async () => {
  const isFormCorrect = await v$.value.$validate();
  if (!isFormCorrect) return;

  try {
    v$.value.$reset();
    router.push({ name: 'home' });
  } catch (err: any) {
    error.value = err.body.message;
    errorModal.value.showModal();
  }
};
</script>

<template>
  <div>
    <ErrorModal :error="error" @gotRef="(el) => (errorModal = el)" />
    <PageHeader>
      <div class="text-xl font-bold">Restore</div>
    </PageHeader>
    <PageContent>
      <form class="grid gap-10" @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-x-5">
          <TextInput
            id="restic_remote"
            title="Restic Remote"
            v-model="state.restic_remote"
            help="Example: rclone:pcloud:Backups/gitea"
            :errors="v$.restic_remote.$errors"
          />
          <TextInput id="local_directory" title="Local directory" v-model="state.local_directory" help="Default: /" />
          <TextInput
            id="password_file_path"
            title="Password file"
            v-model="state.password_file_path"
            help="Example: /secrets/.restipw"
            :errors="v$.password_file_path.$errors"
          />
        </div>
        <div class="flex justify-start flex-row-reverse gap-5">
          <button class="btn btn-primary" type="submit"><i class="fa-solid fa-check"></i>Submit</button>
          <button @click.prevent="router.push({ name: 'home' })" class="btn btn-neutral" type="submit"><i class="fa-solid fa-times"></i>Cancel</button>
        </div>
      </form>
    </PageContent>
  </div>
</template>
