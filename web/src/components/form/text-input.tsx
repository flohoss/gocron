import { component$, type PropFunction, type QwikChangeEvent, type QwikFocusEvent } from '@builder.io/qwik';

type TextInputProps = {
  name: string;
  type: 'text' | 'email' | 'tel' | 'password' | 'url' | 'date';
  label?: string;
  placeholder?: string;
  value: string | undefined;
  error: string;
  required?: boolean;
  ref: PropFunction<(element: Element) => void>;
  onInput$: PropFunction<(event: Event, element: HTMLInputElement) => void>;
  onChange$: PropFunction<(event: QwikChangeEvent<HTMLInputElement>, element: HTMLInputElement) => void>;
  onBlur$: PropFunction<(event: QwikFocusEvent<HTMLInputElement>, element: HTMLInputElement) => void>;
};

export default component$(({ label, error, ...props }: TextInputProps) => {
  const { name, required } = props;
  return (
    <div class="form-control w-full">
      <label class="label">
        <span class="label-text">
          {label} {required && <span>*</span>}
        </span>
      </label>
      <input class={`input input-bordered ${error && 'input-error'} w-full`} {...props} id={name} aria-invalid={!!error} aria-errormessage={`${name}-error`} />
      <label class="label">
        {error && (
          <span class="label-text-alt text-error" id={`${name}-error`}>
            {error}
          </span>
        )}
      </label>
    </div>
  );
});
