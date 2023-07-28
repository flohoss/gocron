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
        'Keep the most recent 7 daily, 4 last-day-of-the-week, 12 or 11 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 11'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        6,
        'Keep the most recent 31 daily, 8 last-day-of-the-week, 24 or 23 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 31 --keep-weekly 9 --keep-monthly 24 --keep-yearly 11'
    );

INSERT OR IGNORE INTO
    "retention_policies"
VALUES
    (
        7,
        'Keep daily for 5 Years, 520 last-day-of-the-week, 121 or 120 last-day-of-the-month & 11 or 10 last-day-of-the-year snapshots',
        '--keep-daily 1095 --keep-weekly 521 --keep-monthly 121 --keep-yearly 11'
    );