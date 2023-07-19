import { component$ } from '@builder.io/qwik';
import { routeLoader$ } from '@builder.io/qwik-city';
import { JobsService } from '~/openapi';

export const useJob = routeLoader$(async (requestEvent) => {
  const jobs = await JobsService.getJobs1(parseInt(requestEvent.params.id));
  return jobs;
});

export default component$(() => {
  const job = useJob();
  return <>{job.value.description}</>;
});
