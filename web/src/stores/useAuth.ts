import { createGlobalState } from '@vueuse/core';
import { computed, ref } from 'vue';
import { getAuthMe, postAuthLogout } from '../client/sdk.gen';
import type { CurrentUserBody } from '../client/types.gen';

export const useAuth = createGlobalState(() => {
  const user = ref<CurrentUserBody | null>(null);
  const ready = ref(false);

  let loading: Promise<void> | null = null;

  const authEnabled = computed(() => user.value?.auth_enabled ?? false);
  const authenticated = computed(() => user.value?.authenticated ?? false);
  const username = computed(() => user.value?.username ?? '');
  const email = computed(() => user.value?.email ?? '');

  const canUseApp = computed(() => ready.value && (!authEnabled.value || authenticated.value));

  function fetchCurrentUser(): Promise<void> {
    loading ??= (async () => {
      try {
        const { data } = await getAuthMe();
        user.value = data ?? null;
      } catch {
        user.value = null;
      } finally {
        ready.value = true;
      }
    })();

    return loading;
  }

  async function logout(): Promise<void> {
    const { data } = await postAuthLogout();
    user.value = null;
    window.location.href = data?.logout_url || '/';
  }

  return { authEnabled, authenticated, username, email, ready, canUseApp, fetchCurrentUser, logout };
});
