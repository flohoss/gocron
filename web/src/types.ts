import { type InjectionKey, type Ref } from 'vue';
import { type database_Job, type database_JobStats, type system_SystemConfig } from './openapi';

export const sseKey = Symbol() as InjectionKey<Ref<SSEvent | undefined>>;

export const emptyJob: database_Job = {
  compression_type: 1,
  description: '',
  id: 0,
  local_directory: '',
  password_file_path: '',
  post_commands: [],
  pre_commands: [],
  restic_remote: '',
  retention_policy: 1,
  routine_check: 0,
};

export interface SelectOption {
  value: number;
  description: string;
}

export const CompressionOptions: SelectOption[] = [
  { value: 1, description: 'Automatic' },
  { value: 2, description: 'Maximum' },
  { value: 3, description: 'No compression' },
];

export const RetentionPolicyOptions: SelectOption[] = [
  { value: 1, description: 'Keep all snapshots' },
  { value: 2, description: 'Keep daily snapshots for the last 2 days' },
  { value: 3, description: 'Keep daily snapshots for the last 7 days' },
  { value: 4, description: 'Keep daily snapshots for the last 31 days' },
  {
    value: 5,
    description: 'Keep the most recent 7 daily, 4 last-day-of-the-week, 12 or 11 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
  },
  {
    value: 6,
    description: 'Keep the most recent 31 daily, 8 last-day-of-the-week, 24 or 23 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
  },
  {
    value: 7,
    description: 'Keep daily for 5 Years, 520 last-day-of-the-week, 121 or 120 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
  },
];

export enum EventType {
  EventCreateRun = 1,
  EventUpdateRun,
  EventCreateLog,
}

export interface SSEvent {
  event_type: EventType;
  content: any;
}

export const emptyJobStats: database_JobStats = {
  check_runs: 0,
  custom_runs: 0,
  error_logs: 0,
  info_logs: 0,
  prune_runs: 0,
  restic_runs: 0,
  total_logs: 0,
  total_runs: 0,
  warning_logs: 0,
  general_runs: 0,
};

export const emptySystemConfig: system_SystemConfig = {
  hostname: '',
  rclone_config_file: '',
  versions: {
    compose: '',
    docker: '',
    go: '',
    gobackup: '',
    rclone: '',
    restic: '',
  },
  config: {
    backup_cron: undefined,
    check_cron: undefined,
    cleanup_cron: undefined,
    healthcheck_url: undefined,
    healthcheck_uuid: undefined,
    identifier: undefined,
    log_level: undefined,
    notification_url: undefined,
    port: undefined,
    swagger_host: undefined,
    time_zone: undefined,
    version: undefined,
  },
};
