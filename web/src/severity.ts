export const Severity = {
  Debug: 1,
  Info: 2,
  Warning: 3,
  Error: 4,
} as const;

export type Severity = (typeof Severity)[keyof typeof Severity];

export function GetColor(severity: number): string {
  switch (severity) {
    case Severity.Debug:
      return 'text-secondary';
    case Severity.Warning:
      return 'text-warning';
    case Severity.Error:
      return 'text-error';
    default:
      return 'text-base-content';
  }
}
