import { component$ } from '@builder.io/qwik';
import { database_Job } from '~/openapi';

export default component$((props: { job: database_Job }) => {
  return (
    <div>
      <div class="text-xl font-bold">{props.job.description}</div>
    </div>
  );
});
