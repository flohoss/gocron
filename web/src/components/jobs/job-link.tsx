import type { PropFunction } from '@builder.io/qwik';
import { component$ } from '@builder.io/qwik';
import { useLocation } from '@builder.io/qwik-city';
import type { database_Job } from '~/openapi';
import NavLink from '../nav/nav-link';

export default component$((props: { job: database_Job; onClick$: PropFunction<() => void> }) => {
  const jobIcon = () => {
    if (props.job.svg_icon) {
      return props.job.svg_icon;
    }
    return props.job.description.charAt(0);
  };

  const loc = useLocation();
  const isActive = (job: database_Job) => {
    return loc.params.id == '' + job.id;
  };

  return (
    <>
      <NavLink
        name={props.job.description}
        extra={props.job.local_directory}
        link={'/jobs/' + props.job.id}
        onClick$={props.onClick$}
        active={isActive(props.job)}
        icon={jobIcon()}
        status="test"
        bg="bg-primary"
        text="text-primary-content"
      />
    </>
  );
});
