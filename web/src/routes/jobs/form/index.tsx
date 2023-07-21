import { component$, $ } from '@builder.io/qwik';
import { routeLoader$, z } from '@builder.io/qwik-city';
import type { InitialValues, SubmitHandler } from '@modular-forms/qwik';
import { formAction$, insert, move, remove, useForm, zodForm$ } from '@modular-forms/qwik';
import SelectInput from '~/components/form/select-input';
import TextInput from '~/components/form/text-input';
import { CompressionTypesService, JobsService, OpenAPI, RetentionPoliciesService } from '~/openapi';

const required = 'This field is required';

const jobSchema = z.object({
  description: z.string().trim().min(1, required),
  local_directory: z.string().trim().min(1, required),
  restic_remote: z.string().trim().min(1, required),
  password_file_path: z.string().trim().min(1, required),
  svg_icon: z.string().trim(),
  compression_type_id: z.string(),
  retention_policy_id: z.string(),
  pre_commands: z.array(z.object({ command: z.string() })),
  post_commands: z.array(z.object({ command: z.string() })),
});

type JobForm = z.infer<typeof jobSchema>;

export const useFormLoader = routeLoader$<InitialValues<JobForm>>(async () => {
  return {
    id: undefined,
    description: '',
    local_directory: '',
    restic_remote: '',
    password_file_path: '',
    svg_icon: '',
    compression_type_id: '1',
    retention_policy_id: '1',
    pre_commands: [],
    post_commands: [],
  };
});

export const useFormAction = formAction$<JobForm>(() => {}, zodForm$(jobSchema));

export const useCompressionTypesLoader = routeLoader$(async () => {
  OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
  const compression_types = await CompressionTypesService.getCompressionTypes();
  return compression_types;
});

export const useRetentionPoliciesLoader = routeLoader$(async () => {
  OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
  const retention_policies = await RetentionPoliciesService.getRetentionPolicies();
  return retention_policies;
});

export default component$(() => {
  const [jobForm, { Form, Field, FieldArray }] = useForm<JobForm>({
    loader: useFormLoader(),
    fieldArrays: ['pre_commands', 'post_commands'],
    action: useFormAction(),
    validate: zodForm$(jobSchema),
  });

  const handleSubmit$: SubmitHandler<JobForm> = $(async (values: any) => {
    OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
    const response = await JobsService.postJobs(values);
    console.log(response);
  });

  const compression_types = useCompressionTypesLoader();
  const retention_policies = useRetentionPoliciesLoader();

  return (
    <Form onSubmit$={handleSubmit$}>
      <div class="grid gap-x-2 grid-cols-1 lg:grid-cols-2">
        <Field name="description">
          {(field, props) => <TextInput {...props} type="text" label="Description" value={field.value} error={field.error} required />}
        </Field>
        <Field name="local_directory">
          {(field, props) => <TextInput {...props} type="text" label="Local Directory" value={field.value} error={field.error} required />}
        </Field>
        <Field name="restic_remote">
          {(field, props) => <TextInput {...props} type="text" label="Restic remote" value={field.value} error={field.error} required />}
        </Field>
        <Field name="password_file_path">
          {(field, props) => <TextInput {...props} type="text" label="Password file" value={field.value} error={field.error} required />}
        </Field>
        <Field name="svg_icon">
          {(field, props) => <TextInput {...props} type="text" label="SVG-Icon" value={field.value} error={field.error} classes="col-span-1 lg:col-span-2" />}
        </Field>
        <Field name="compression_type_id">
          {(field, props) => <SelectInput {...props} label="Compression" value={field.value} options={compression_types.value} error={field.error} required />}
        </Field>
        <Field name="retention_policy_id">
          {(field, props) => (
            <SelectInput {...props} label="Retention policy" value={field.value} options={retention_policies.value} error={field.error} required />
          )}
        </Field>
        <div class="flex flex-col">
          <label class="label">
            <span class="label-text">Commands before backup</span>
          </label>
          <FieldArray name="pre_commands">
            {(fieldArray) => (
              <>
                {fieldArray.items.map((item, index) => (
                  <div key={item} class="flex space-x-5">
                    <Field name={`${fieldArray.name}.${index}.command`}>
                      {(field, props) => (
                        <TextInput {...props} type="text" value={field.value} error={field.error}>
                          <div
                            class="btn btn-error"
                            onClick$={() =>
                              remove(jobForm, fieldArray.name, {
                                at: index,
                              })
                            }
                          >
                            <i class="fa-solid fa-xmark"></i>
                          </div>
                        </TextInput>
                      )}
                    </Field>
                  </div>
                ))}
                <div class="btn btn-ghost gap-2" onClick$={() => insert(jobForm, fieldArray.name, { value: { command: '' } })}>
                  <i class="fa-solid fa-plus"></i>
                  New command
                </div>
              </>
            )}
          </FieldArray>
        </div>
        <div class="flex flex-col">
          <label class="label">
            <span class="label-text">Commands after backup</span>
          </label>
          <FieldArray name="post_commands">
            {(fieldArray) => (
              <>
                {fieldArray.items.map((item, index) => (
                  <div key={item}>
                    <Field name={`${fieldArray.name}.${index}.command`}>
                      {(field, props) => (
                        <TextInput {...props} type="text" value={field.value} error={field.error}>
                          <div
                            class="btn btn-error"
                            onClick$={() =>
                              remove(jobForm, fieldArray.name, {
                                at: index,
                              })
                            }
                          >
                            <i class="fa-solid fa-xmark"></i>
                          </div>
                        </TextInput>
                      )}
                    </Field>
                  </div>
                ))}
                <div class="btn btn-ghost gap-2" onClick$={() => insert(jobForm, fieldArray.name, { value: { command: '' } })}>
                  <i class="fa-solid fa-plus"></i>
                  New command
                </div>
              </>
            )}
          </FieldArray>
        </div>
      </div>
      <div class="mt-5 flex gap-4 justify-start flex-row-reverse">
        <button class="btn btn-primary" type="submit">
          Save
        </button>
        <button class="btn" type="submit">
          Cancel
        </button>
      </div>
    </Form>
  );
});
