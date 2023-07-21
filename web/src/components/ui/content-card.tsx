import { Slot, component$ } from '@builder.io/qwik';

export default component$((props: { classes?: string }) => {
  return (
    <div class={`grid lg:gap-x-6 grid-cols-1 ${props.classes ? props.classes : ''}`}>
      <Slot />
    </div>
  );
});
