import { component$, Slot, type PropFunction, type QwikChangeEvent, type QwikFocusEvent } from '@builder.io/qwik';

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
  classes?: string;
  example?: string;
  suffix?: string;
};

export default component$(({ label, error, classes, example, suffix, ...props }: TextInputProps) => {
  const { name, required } = props;
  return (
    <>
      <div class={`form-control w-full ${classes ? classes : ''}`}>
        {label && (
          <label class="label">
            <span class="label-text">
              {label} {required && <span>*</span>}
            </span>
          </label>
        )}
        <div class="join">
          <input
            class={`input input-bordered bg-base-200 ${error ? 'input-error' : ''} join-item w-full`}
            {...props}
            id={name}
            aria-invalid={!!error}
            aria-errormessage={`${name}-error`}
          />
          {suffix && <div class="flex justify-center items-center join-item px-4 bg-base-300 border border-gray-500/40">{suffix}</div>}
          <Slot />
        </div>
        <label class="label">
          {error ? (
            <span class="label-text-alt text-error" id={`${name}-error`}>
              {error}
            </span>
          ) : (
            example && (
              <span class="label-text-alt opacity-50 select-text" id={`${name}-example`}>
                Example: "{example}"
              </span>
            )
          )}
        </label>
      </div>
    </>
  );
});
