import type { database_Job } from './openapi';

export const emptyJob: database_Job = {
  compression_type_id: 0,
  description: '',
  local_directory: '',
  password_file_path: '',
  restic_remote: '',
  retention_policy_id: 0,
};
