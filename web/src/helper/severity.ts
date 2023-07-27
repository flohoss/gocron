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
      return 'check';
    case database_LogSeverity.LogWarning:
      return 'triangle-exclamation';
    case database_LogSeverity.LogError:
      return 'exclamation';
    default:
      return '';
  }
};
