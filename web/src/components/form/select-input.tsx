import { component$, QwikFocusEvent, type PropFunction, type QwikChangeEvent } from '@builder.io/qwik';
import type { database_SelectOption } from '~/openapi';

type SelectInputProps = {
  name: string;
  label?: string;
  value: string | undefined;
  error: string;
  required?: boolean;
  ref: PropFunction<(element: Element) => void>;
  onInput$: PropFunction<(event: Event, element: HTMLSelectElement) => void>;
  onChange$: PropFunction<(event: QwikChangeEvent<HTMLSelectElement>, element: HTMLSelectElement) => void>;
  onBlur$: PropFunction<(event: QwikFocusEvent<HTMLSelectElement>, element: HTMLSelectElement) => void>;
  classes?: string;
  options: database_SelectOption[];
};

export default component$(({ label, error, classes, options, ...props }: SelectInputProps) => {
  const { name, required } = props;
  return (
    <div class={`form-control w-full ${classes}`}>
      <label class="label">
        <span class="label-text">
          {label} {required && <span>*</span>}
        </span>
      </label>
      <select class="select select-bordered" {...props} id={name} aria-invalid={!!error} aria-errormessage={`${name}-error`}>
        {options.map(({ value, name }) => (
          <option key={value} value={value} selected={value === props.value}>
            {name}
          </option>
        ))}
      </select>
      <label class="label">
        {error && (
          <span id={`${name}-error`} class="label-text-alt text-error">
            {error}
          </span>
        )}
      </label>
    </div>
  );
});
