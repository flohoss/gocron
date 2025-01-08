CREATE TABLE IF NOT EXISTS
    status (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        status TEXT NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS
    severities (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        severity TEXT NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS
    jobs (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        cron TEXT NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    envs (
        id uuid PRIMARY KEY DEFAULT (LOWER(REPLACE(HEX(RANDOMBLOB(16)), '-', ''))),
        job_id TEXT NOT NULL,
        KEY TEXT NOT NULL,
        value TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    commands (
        id uuid PRIMARY KEY DEFAULT (LOWER(REPLACE(HEX(RANDOMBLOB(16)), '-', ''))),
        job_id TEXT NOT NULL,
        command TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    runs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_id TEXT NOT NULL,
        status_id INTEGER NOT NULL,
        start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        end_time DATETIME,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (status_id) REFERENCES status (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        run_id INTEGER NOT NULL,
        severity_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (run_id) REFERENCES runs (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (severity_id) REFERENCES severities (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE VIEW IF NOT EXISTS
    runs_view AS
SELECT
    id,
    job_id,
    status_id,
    DATETIME(runs.start_time, 'localtime') AS start_time,
    CASE
        WHEN runs.end_time IS NOT NULL THEN DATETIME(runs.end_time, 'localtime')
        ELSE NULL
    END AS end_time,
    CASE
        WHEN runs.end_time IS NOT NULL THEN PRINTF(
            '%02dh%02dm%02ds',
            CAST(
                (
                    JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                ) * 24 AS INTEGER
            ),
            CAST(
                (
                    (
                        JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                    ) * 24 * 60
                ) % 60 AS INTEGER
            ),
            CAST(
                (
                    (
                        JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                    ) * 24 * 60 * 60
                ) % 60 AS INTEGER
            )
        )
        ELSE 'N/A'
    END AS duration
FROM
    runs;