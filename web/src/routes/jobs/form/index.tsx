import { component$, $, useContext } from '@builder.io/qwik';
import { routeLoader$, useLocation, z, useNavigate, Link } from '@builder.io/qwik-city';
import type { InitialValues, SubmitHandler } from '@modular-forms/qwik';
import { formAction$, insert, move, remove, useForm, zodForm$ } from '@modular-forms/qwik';
import SelectInput from '~/components/form/select-input';
import TextInput from '~/components/form/text-input';
import ContentCard from '~/components/ui/content-card';
import { CompressionTypesService, JobsService, OpenAPI, RetentionPoliciesService } from '~/openapi';
import { JobContext } from '~/routes/layout';

const jobSchema = z.object({
  description: z.string().trim().nonempty(),
  local_directory: z.string().trim().nonempty(),
  restic_remote: z.string().trim().nonempty(),
  password_file_path: z.string().trim().nonempty(),
  svg_icon: z.string().trim().optional(),
  routine_check: z.string().trim().optional(),
  compression_type_id: z.string(),
  retention_policy_id: z.string(),
  pre_commands: z.array(z.object({ command: z.string().nonempty() })).optional(),
  post_commands: z.array(z.object({ command: z.string().nonempty() })).optional(),
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
    routine_check: '',
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
  const ctx = useContext(JobContext);
  const nav = useNavigate();
  const [jobForm, { Form, Field, FieldArray }] = useForm<JobForm>({
    loader: useFormLoader(),
    fieldArrays: ['pre_commands', 'post_commands'],
    action: useFormAction(),
    validate: zodForm$(jobSchema),
  });

  const handleSubmit$: SubmitHandler<JobForm> = $(async (values: any) => {
    OpenAPI.BASE = import.meta.env.PUBLIC_API_URL;
    values = {
      ...values,
      compression_type_id: parseInt(values.compression_type_id),
      retention_policy_id: parseInt(values.retention_policy_id),
    };
    JobsService.postJobs(values)
      .then((res) => {
        ctx.jobs.push(res);
        ctx.jobs.sort();
        nav('/');
      })
      .catch((err) => console.log(err));
  });

  const compression_types = useCompressionTypesLoader();
  const retention_policies = useRetentionPoliciesLoader();

  return (
    <Form onSubmit$={handleSubmit$} class="grid gap-5 grid-cols-1">
      <ContentCard classes="lg:grid-cols-2">
        <Field name="description">
          {(field, props) => <TextInput {...props} type="text" label="Description" value={field.value} error={field.error} required example="Gitea" />}
        </Field>
        <Field name="local_directory">
          {(field, props) => (
            <TextInput {...props} type="text" label="Local Directory" value={field.value} error={field.error} required example="/opt/docker/gitea" />
          )}
        </Field>
        <Field name="restic_remote">
          {(field, props) => (
            <TextInput {...props} type="text" label="Restic remote" value={field.value} error={field.error} required example="rclone:pcloud:Backups/gitea" />
          )}
        </Field>
        <Field name="password_file_path">
          {(field, props) => (
            <TextInput {...props} type="text" label="Password file" value={field.value} error={field.error} required example="/secrets/.restipw" />
          )}
        </Field>
      </ContentCard>
      <ContentCard classes="lg:grid-cols-2">
        <Field name="svg_icon">
          {(field, props) => (
            <TextInput {...props} type="text" label="SVG-Icon" value={field.value} error={field.error} example={`<i class="fa-solid fa-circle-nodes"></i>`} />
          )}
        </Field>
        <Field name="routine_check">
          {(field, props) => <TextInput {...props} type="text" label="Routine check" value={field.value} error={field.error} example="15" suffix="%" />}
        </Field>
        <Field name="compression_type_id">
          {(field, props) => <SelectInput {...props} label="Compression" value={field.value} options={compression_types.value} error={field.error} required />}
        </Field>
        <Field name="retention_policy_id">
          {(field, props) => (
            <SelectInput {...props} label="Retention policy" value={field.value} options={retention_policies.value} error={field.error} required />
          )}
        </Field>
      </ContentCard>
      <ContentCard classes="lg:grid-cols-2">
        <FieldArray name="pre_commands">
          {(fieldArray) => (
            <>
              <div class="flex flex-col">
                <label class="label">
                  <span class="label-text">
                    Commands before backup
                    <br />
                    <span class="label-text-alt opacity-50">Example: "docker-compose down"</span>
                  </span>
                </label>
                {fieldArray.items.map((item, index) => (
                  <div key={item} class="flex space-x-5">
                    <Field name={`${fieldArray.name}.${index}.command`}>
                      {(field, props) => (
                        <TextInput {...props} type="text" value={field.value} error={field.error}>
                          <div
                            class="btn btn-error join-item"
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
                <div class="btn gap-2" onClick$={() => insert(jobForm, fieldArray.name, { value: { command: '' } })}>
                  <i class="fa-solid fa-plus"></i>
                  New command
                </div>
              </div>
            </>
          )}
        </FieldArray>
        <FieldArray name="post_commands">
          {(fieldArray) => (
            <>
              <div class="flex flex-col">
                <label class="label">
                  <span class="label-text">
                    Commands after backup
                    <br />
                    <span class="label-text-alt opacity-50">Example: "docker-compose up -d"</span>
                  </span>
                </label>
                {fieldArray.items.map((item, index) => (
                  <div key={item}>
                    <Field name={`${fieldArray.name}.${index}.command`}>
                      {(field, props) => (
                        <TextInput {...props} type="text" value={field.value} error={field.error}>
                          <div
                            class="btn btn-error join-item"
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
                <div class="btn gap-2" onClick$={() => insert(jobForm, fieldArray.name, { value: { command: '' } })}>
                  <i class="fa-solid fa-plus"></i>
                  New command
                </div>
              </div>
            </>
          )}
        </FieldArray>
      </ContentCard>
      <div class="mt-5 flex gap-4 justify-start flex-row-reverse">
        <button class="btn btn-primary" type="submit">
          <i class="fa-solid fa-check"></i>
          Save
        </button>
        <Link href="/" class="btn" type="reset">
          <i class="fa-solid fa-times"></i>
          Cancel
        </Link>
      </div>
    </Form>
  );
});
