<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const starting = ref(false);

const errorMessages: Record<string, string> = {
  rejected: 'The identity provider rejected the login.',
  expired: 'The login session expired, please try again.',
  unavailable: 'The identity provider is not reachable.',
  failed: 'Authentication failed, please try again.',
};

const errorMessage = computed(() => {
  const code = route.query.error;
  if (typeof code !== 'string' || !(code in errorMessages)) return '';
  return errorMessages[code];
});

function startLogin() {
  starting.value = true;
  const redirect = route.query.redirect;
  const target = typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '';
  window.location.href = '/api/auth/login' + (target ? `?redirect=${encodeURIComponent(target)}` : '');
}
</script>

<template>
  <div class="hero min-h-dvh bg-base-200">
    <div class="hero-content w-full flex-col">
      <div class="card w-full max-w-sm bg-base-100 shadow-xl">
        <div class="card-body gap-4">
          <img class="mx-auto h-20" src="/static/logo.webp" alt="GoCron" />

          <div v-if="errorMessage" role="alert" class="alert alert-error" data-test-id="login-error">
            <span class="icon-[bi--exclamation-triangle] size-5"></span>
            <span>{{ errorMessage }}</span>
          </div>

          <button type="button" class="btn btn-primary w-full gap-2" :disabled="starting" data-test-id="login-button" @click="startLogin">
            <span v-if="starting" class="loading loading-spinner loading-sm"></span>
            <span class="size-5 icon-[cib--open-id]"></span>
            <span>Sign in with OIDC</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
