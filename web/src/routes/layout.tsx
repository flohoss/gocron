import { component$, Slot } from '@builder.io/qwik';
import { routeLoader$ } from '@builder.io/qwik-city';
import Link from '~/components/navigation/link';
import { JobsService } from '~/openapi';

export const useJobs = routeLoader$(async () => {
  const jobs = await JobsService.getJobs();
  return jobs;
});

export default component$(() => {
  const jobs = useJobs();
  return (
    <div class="drawer lg:drawer-open">
      <input id="drawer" type="checkbox" class="drawer-toggle" />
      <div class="drawer-content">
        <Slot />
      </div>
      <div class="drawer-side">
        <label for="drawer" class="drawer-overlay"></label>
        <ul class="menu p-4 w-80 h-full bg-base-200 text-base-content">
          {jobs.value.map((job) => (
            <Link name={job.description} />
          ))}
        </ul>
      </div>
    </div>
  );
});
