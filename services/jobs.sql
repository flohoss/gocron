CREATE TABLE
    IF NOT EXISTS jobs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        backup_schedule_id INTEGER NOT NULL,
        FOREIGN KEY (backup_schedule_id) REFERENCES backup_schedules (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS runs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_id INTEGER NOT NULL,
        start_time DATETIME NOT NULL,
        end_time DATETIME,
        status TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        run_id INTEGER NOT NULL,
        log_severity TEXT NOT NULL,
        message TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        FOREIGN KEY (run_id) REFERENCES runs (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS envs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        value TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS commands (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_id INTEGER NOT NULL,
        command_order INTEGER NOT NULL,
        command_text TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS backup_schedules (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        cron_expression TEXT NOT NULL,
        description TEXT
    );