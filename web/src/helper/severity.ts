import { database_LogSeverity } from '@/openapi';

export const severityColor = (severity: number | undefined) => {
  switch (severity) {
    case database_LogSeverity.LogWarning:
      return 'text-warning';
    case database_LogSeverity.LogError:
      return 'text-error';
    default:
      return '';
  }
};

export const severityIcons = (severity: number | undefined) => {
  switch (severity) {
    case database_LogSeverity.LogInfo:
      return '<i class="fa-solid fa-check"></i>';
    case database_LogSeverity.LogWarning:
      return `<i class="fa-solid fa-triangle-exclamation"></i>`;
    case database_LogSeverity.LogError:
      return `<i class="fa-solid fa-exclamation"></i>`;
    default:
      return `<span class="loading loading-spinner"></span>`;
  }
};
