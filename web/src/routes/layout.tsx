import { component$, createContextId, Slot, useContextProvider, useSignal, useStore, useTask$ } from '@builder.io/qwik';
import { useLocation } from '@builder.io/qwik-city';
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
    return loc.url.pathname == when;
  };

  useContextProvider(JobContext, store);

  return (
    <div class="drawer lg:drawer-open">
      <input id="drawer" type="checkbox" class="drawer-toggle" />
      <div class="drawer-content p-2 lg:p-4">
        <Slot />
        <div class="lg:hidden btm-nav bg-base-200">
          <button>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              />
            </svg>
          </button>
          <label class="active" for="drawer" ref={drawerRef}>
            <i class="fa-solid fa-list-ul"></i>
            <div class="text-xs opacity-75">Jobs</div>
          </label>
          <button>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              />
            </svg>
          </button>
        </div>
      </div>
      <div class="drawer-side">
        <label for="drawer" class="drawer-overlay"></label>
        <ul class="menu p-2 w-80 h-full bg-base-200 text-base-content flex flex-col flex-nowrap gap-2 overflow-y-auto">
          <NavLink
            link="/"
            name="Dashboard"
            icon={`<i class="fa-solid fa-circle-nodes"></i>`}
            active={isActive('/')}
            onClick$={() => drawerRef.value && drawerRef.value.click()}
          />
          <div class="my-2"></div>
          {store.jobs.map((job) => (
            <JobLink key={job.id} job={job} onClick$={() => drawerRef.value && drawerRef.value.click()} />
          ))}
          <div class="my-2"></div>
          <NavLink
            link="/jobs/form"
            name="Add"
            icon={`<i class="fa-solid fa-plus"></i>`}
            active={isActive('/jobs/form/')}
            onClick$={() => drawerRef.value && drawerRef.value.click()}
          />
        </ul>
      </div>
    </div>
  );
});
