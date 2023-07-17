CREATE TABLE IF NOT EXISTS
    jobs (
        job_id INTEGER PRIMARY KEY,
        description TEXT NOT NULL,
        local_directory TEXT NOT NULL,
        restic_remote TEXT NOT NULL,
        restart_option INTEGER NOT NULL,
        password_file_path TEXT NOT NULL,
        compression_type TEXT NOT NULL,
        svg_icon TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        retention_policy_id INTEGER NOT NULL,
        FOREIGN KEY (retention_policy_id) REFERENCES retention_policies (retention_policy_id) ON DELETE SET NULL
    );

CREATE TABLE IF NOT EXISTS
    retention_policies (
        retention_policy_id INTEGER PRIMARY KEY,
        retention_policy TEXT NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS
    commands (
        command_id INTEGER PRIMARY KEY,
        job_id INTEGER NOT NULL,
        command_type TEXT NOT NULL,
        command TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (job_id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    logs (
        log_id INTEGER PRIMARY KEY,
        job_id INTEGER NOT NULL,
        log_type TEXT NOT NULL,
        message TEXT NOT NULL,
        FOREIGN KEY (job_id) REFERENCES jobs (job_id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS
    runs (
        run_id INTEGER PRIMARY KEY,
        job_id INTEGER NOT NULL,
        start_time DATETIME NOT NULL,
        end_time DATETIME NOT NULL,
        log_id INTEGER,
        FOREIGN KEY (job_id) REFERENCES jobs (job_id) ON DELETE CASCADE,
        FOREIGN KEY (log_id) REFERENCES logs (log_id) ON DELETE SET NULL
    );