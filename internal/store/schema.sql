-- Chronicle metadata schema.
--
-- Git remains the source of truth for source code. This table records the
-- development intent behind each AI interaction and points at the patch file
-- that captures the resulting diff.
--
-- Scope is per-project: each project has its own .chronicle/metadata.db, so the
-- events table holds exactly that project's history.

CREATE TABLE IF NOT EXISTS events (
    id          TEXT PRIMARY KEY,
    prompt      TEXT NOT NULL,
    model       TEXT NOT NULL,
    timestamp   TEXT NOT NULL,
    patch_path  TEXT NOT NULL,
    commit_hash TEXT
);

CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
