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
    runs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_name TEXT NOT NULL,
        job_name_normalized TEXT GENERATED ALWAYS AS (LOWER(TRIM(job_name))) STORED,
        status_id INTEGER NOT NULL,
        start_time INTEGER NOT NULL,
        end_time INTEGER,
        FOREIGN KEY (status_id) REFERENCES status (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_runs_job_name_normalized_start_time ON runs (job_name_normalized, start_time DESC);

CREATE TABLE IF NOT EXISTS
    logs (
        created_at INTEGER PRIMARY KEY,
        run_id INTEGER NOT NULL,
        severity_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        FOREIGN KEY (run_id) REFERENCES runs (id) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY (severity_id) REFERENCES severities (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );