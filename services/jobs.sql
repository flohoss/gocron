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
        id INTEGER PRIMARY KEY,
        job_id TEXT NOT NULL,
        KEY TEXT NOT NULL,
        value TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    commands (
        id INTEGER PRIMARY KEY,
        job_id TEXT NOT NULL,
        command TEXT NOT NULL,
        file_output TEXT,
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
        created_at INTEGER PRIMARY KEY,
        run_id INTEGER NOT NULL,
        severity_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        FOREIGN KEY (run_id) REFERENCES runs (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (severity_id) REFERENCES severities (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE VIEW IF NOT EXISTS
    runs_view AS
SELECT
    id,
    job_id,
    status_id,
    CASE
        WHEN start_time IS NOT NULL THEN STRFTIME('%H:%M:%S', start_time)
        ELSE NULL
    END AS start_time,
    CASE
        WHEN start_time IS NOT NULL THEN STRFTIME('%Y-%m-%d', start_time)
        ELSE NULL
    END AS start_date,
    CASE
        WHEN end_time IS NOT NULL THEN STRFTIME('%H:%M:%S', end_time)
        ELSE NULL
    END AS end_time,
    CASE
        WHEN end_time IS NOT NULL THEN STRFTIME('%Y-%m-%d', end_time)
        ELSE NULL
    END AS end_date,
    CASE
        WHEN end_time IS NOT NULL THEN PRINTF(
            '%dh%dm%ds',
            CAST(
                (
                    JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                ) * 24 AS INTEGER
            ),
            CAST(
                (
                    JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                ) * 24 * 60 % 60 AS INTEGER
            ),
            CAST(
                (
                    JULIANDAY(runs.end_time) - JULIANDAY(runs.start_time)
                ) * 24 * 60 * 60 % 60 AS INTEGER
            )
        )
        ELSE NULL
    END AS duration,
    NULL AS logs
FROM
    runs;

CREATE VIEW IF NOT EXISTS
    jobs_view AS
SELECT
    id,
    name,
    cron,
    NULL AS runs
FROM
    jobs;