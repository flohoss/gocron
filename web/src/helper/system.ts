import type { system_Data } from '@/openapi';

export const emptySystem: system_Data = {
  configuration: {
    hostname: '',
    rclone_config_file: '',
  },
  job_stats: {
    check_runs: 0,
    custom_runs: 0,
    prune_runs: 0,
    restic_runs: 0,
    total_logs: 0,
    total_runs: 0,
    error_logs: 0,
    warning_logs: 0,
    info_logs: 0,
  },
  versions: {
    compose: '',
    docker: '',
    go: '',
    gobackup: '',
    rclone: '',
    restic: '',
  },
};
