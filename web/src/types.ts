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
