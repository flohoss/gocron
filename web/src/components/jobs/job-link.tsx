import { PropFunction, component$ } from '@builder.io/qwik';
import { Link, useLocation } from '@builder.io/qwik-city';
import { database_Job } from '~/openapi';

export default component$((props: { job: database_Job; onClick$: PropFunction<() => void> }) => {
  const jobIcon = () => {
    if (props.job.svg_icon) {
      return props.job.svg_icon;
    }
    return props.job.description.charAt(0);
  };

  const loc = useLocation();
  const isActive = (job: database_Job) => {
    if (loc.params.id == '' + job.id) return 'active';
  };

  return (
    <li>
      <Link onClick$={props.onClick$} href={'/jobs/' + props.job.id} class={`flex justify-between px-2 ${isActive(props.job)}`}>
        <div class="flex items-center gap-2">
          <div
            class="w-10 h-10 p-2 rounded-full bg-primary flex items-center justify-center shrink-0 text-lg text-primary-content"
            dangerouslySetInnerHTML={jobIcon()}
          ></div>
          <div class="flex flex-col">
            <div class="text-lg">{props.job.description}</div>
            <div class="opacity-50 text-xs">{props.job.local_directory}</div>
          </div>
        </div>
        <div class="relative flex h-3 w-3">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
          <span class="relative inline-flex rounded-full h-3 w-3 bg-success"></span>
        </div>
      </Link>
    </li>
  );
});
