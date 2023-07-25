/**
 * Compression Types
 */
INSERT OR IGNORE INTO
    "compression_types"
VALUES
    (1, 'Automatic', 'auto');

INSERT OR IGNORE INTO
    "compression_types"
VALUES
    (2, 'Maximum', 'max');

INSERT OR IGNORE INTO
    "compression_types"
VALUES
    (3, 'No compression', 'off');

/**
 * Log Severity
 */
INSERT OR IGNORE INTO
    "log_severities"
VALUES
    (1, 'Info');

INSERT OR IGNORE INTO
    "log_severities"
VALUES
    (2, 'Warning');

INSERT OR IGNORE INTO
    "log_severities"
VALUES
    (3, 'Error');

/**
 * Log Types
 */
INSERT OR IGNORE INTO
    "log_types"
VALUES
    (1, 'General');

INSERT OR IGNORE INTO
    "log_types"
VALUES
    (2, 'Restic');

INSERT OR IGNORE INTO
    "log_types"
VALUES
    (3, 'CustomCommand');

INSERT OR IGNORE INTO
    "log_types"
VALUES
    (4, 'Prune');

INSERT OR IGNORE INTO
    "log_types"
VALUES
    (5, 'Check');

/**
 * Retention Policies
 */
INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (1, 'Keep all snapshots', '');

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        2,
        'Keep daily snapshots for the last 2 days',
        '--keep-daily 2'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        3,
        'Keep daily snapshots for the last 7 days',
        '--keep-daily 7'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        4,
        'Keep daily snapshots for the last 31 days',
        '--keep-daily 31'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        5,
        'Keep the most recent 7 daily, 4 last-day-of-the-week, 12 or 11 last-day-of-the-month & 75 or 74 last-day-of-the-year snapshots',
        '--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        6,
        'Keep the most recent 31 daily, 8 last-day-of-the-week, 24 or 23 last-day-of-the-month & 75 or 74 last-day-of-the-year snapshots',
        '--keep-daily 31 --keep-weekly 9 --keep-monthly 24 --keep-yearly 75'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        7,
        'Keep daily for 5 Years, 520 last-day-of-the-week, 121 or 120 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 1095 --keep-weekly 521 --keep-monthly 121 --keep-yearly 11'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        8,
        'Keep the most recent 7 daily, 4 last-day-of-the-week, 12 or 11 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 11'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        9,
        'Keep the most recent 31 daily, 8 last-day-of-the-week, 24 or 23 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 31 --keep-weekly 9 --keep-monthly 24 --keep-yearly 11'
    );