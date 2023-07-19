import { component$ } from '@builder.io/qwik';

export default component$((props: { name: string }) => {
  return (
    <li>
      <a>{props.name}</a>
    </li>
  );
});
