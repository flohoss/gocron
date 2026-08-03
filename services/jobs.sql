CREATE TABLE IF NOT EXISTS runs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  job_name TEXT NOT NULL,
  job_slug TEXT NOT NULL,
  status_id INTEGER NOT NULL,
  start_time INTEGER NOT NULL,
  end_time INTEGER
);

CREATE INDEX IF NOT EXISTS idx_runs_job_slug_start_time ON runs (job_slug, start_time DESC);

CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);

CREATE TABLE IF NOT EXISTS logs (
  created_at INTEGER PRIMARY KEY,
  run_id INTEGER NOT NULL,
  severity_id INTEGER NOT NULL,
  message TEXT NOT NULL,
  FOREIGN KEY (run_id) REFERENCES runs (id) ON DELETE CASCADE ON UPDATE CASCADE
);
