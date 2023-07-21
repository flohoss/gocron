import { component$ } from '@builder.io/qwik';
import { routeLoader$, type DocumentHead } from '@builder.io/qwik-city';
import { OpenAPI, SystemService, system_Data } from '~/openapi';

export const useSystemLoader = routeLoader$(async () => {
  OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
  const system = await SystemService.getSystem();
  return system;
});

export default component$(() => {
  const system = useSystemLoader();
  return (
    <div>
      <div class="text-2xl">{system.value.configuration?.hostname}</div>
      <div>{system.value.versions?.go}</div>
    </div>
  );
});

export const head: DocumentHead = {
  title: 'Welcome to Qwik',
  meta: [
    {
      name: 'description',
      content: 'Qwik site description',
    },
  ],
};
