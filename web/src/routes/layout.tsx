import { component$, Slot } from '@builder.io/qwik';
import { Link, routeLoader$ } from '@builder.io/qwik-city';
import { JobsService, OpenAPI } from '~/openapi';

export const useJobs = routeLoader$(async () => {
  OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
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
            <li>
              <Link href={'/jobs/' + job.id}>{job.description}</Link>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
});
