import { component$, useComputed$, useContext } from '@builder.io/qwik';
import { useLocation } from '@builder.io/qwik-city';
import JobCard from '~/components/jobs/job-card';
import { database_Job } from '~/openapi';
import { JobContext } from '~/routes/layout';
import { emptyJob } from '~/types';

export default component$(() => {
  const ctx = useContext(JobContext);
  const loc = useLocation();
  const job = useComputed$<database_Job>(() => ctx.jobs.find((j) => '' + j.id == loc.params.id) || emptyJob);
  return (
    <>
      <JobCard job={job.value} />
    </>
  );
});
