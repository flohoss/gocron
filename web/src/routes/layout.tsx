import { component$, createContextId, Slot, useContextProvider, useSignal, useStore, useTask$ } from '@builder.io/qwik';
import { Link, useLocation } from '@builder.io/qwik-city';
import { isServer } from '@builder.io/qwik/build';
import JobLink from '~/components/jobs/job-link';
import NavLink from '~/components/nav/nav-link';
import type { database_Job } from '~/openapi';
import { JobsService, OpenAPI } from '~/openapi';

export interface JobContextType {
  jobs: database_Job[];
}

export const JobContext = createContextId<JobContextType>('jobs');

export default component$(() => {
  const store = useStore<JobContextType>({ jobs: [] });
  const drawerRef = useSignal<HTMLElement>();

  useTask$(async () => {
    if (isServer) {
      OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
      store.jobs = await JobsService.getJobs();
    }
  });

  const loc = useLocation();
  const isActive = (when: string) => {
    if (loc.url.pathname === when) {
      return true;
    } else if (when !== '/' && loc.url.pathname.startsWith(when)) {
      return true;
    }
    return false;
  };

  useContextProvider(JobContext, store);

  return (
    <div class="drawer lg:drawer-open">
      <input id="drawer" type="checkbox" class="drawer-toggle" />
      <div class="drawer-content p-2 md:p-5 lg:p-10">
        <Slot />
        <div class="my-20 lg:hidden"></div>
        <div class="lg:hidden btm-nav bg-base-200">
          <Link class={`${isActive('/') ? 'active' : ''}`} href="/">
            <i class="fa-solid fa-circle-nodes"></i>
            <div class="text-xs opacity-75">Dashboard</div>
          </Link>
          <label for="drawer" ref={drawerRef}>
            <i class="fa-solid fa-list-ul"></i>
            <div class="text-xs opacity-75">Jobs</div>
          </label>
          <Link class={`${isActive('/jobs/form') ? 'active' : ''}`} href="/jobs/form">
            <i class="fa-solid fa-plus"></i>
            <div class="text-xs opacity-75">New</div>
          </Link>
        </div>
      </div>
      <div class="drawer-side">
        <label for="drawer" class="drawer-overlay"></label>
        <ul class="menu p-2 w-80 h-full bg-base-200 text-base-content flex flex-col flex-nowrap gap-4 overflow-y-auto">
          <NavLink link="/" name="Dashboard" icon={`<i class="fa-solid fa-circle-nodes"></i>`} active={isActive('/')} hidden={true} />
          <div class="flex flex-col flex-nowrap gap-2">
            {store.jobs.map((job) => (
              <JobLink key={job.id} job={job} onClick$={() => drawerRef.value && drawerRef.value.click()} />
            ))}
          </div>
          <NavLink link="/jobs/form" name="New" icon={`<i class="fa-solid fa-plus"></i>`} active={isActive('/jobs/form')} hidden={true} />
        </ul>
      </div>
    </div>
  );
});
