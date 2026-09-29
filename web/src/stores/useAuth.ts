import { createGlobalState } from '@vueuse/core';
import { computed, ref } from 'vue';
import { getAuthMe, postAuthLogout } from '../client/sdk.gen';
import type { CurrentUserBody } from '../client/types.gen';

export const useAuth = createGlobalState(() => {
  const user = ref<CurrentUserBody | null>(null);
  const ready = ref(false);
  const unavailable = ref(false);

  let loading: Promise<void> | null = null;

  const authEnabled = computed(() => unavailable.value || (user.value?.auth_enabled ?? false));
  const authenticated = computed(() => user.value?.authenticated ?? false);

  const canUseApp = computed(() => ready.value && (!authEnabled.value || authenticated.value));

  function fetchCurrentUser(): Promise<void> {
    loading ??= (async () => {
      try {
        const response = await getAuthMe();
        if (response.error || !response.data) throw new Error('me request failed');
        user.value = response.data;
      } catch {
        user.value = null;
        unavailable.value = true;
      } finally {
        ready.value = true;
      }
    })();

    return loading;
  }

  async function logout(): Promise<void> {
    unavailable.value = false;
    const { data } = await postAuthLogout();
    user.value = null;
    window.location.href = data?.logout_url || '/';
  }

  return { authEnabled, authenticated, ready, canUseApp, fetchCurrentUser, logout };
});
