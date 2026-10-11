CREATE TABLE IF NOT EXISTS requests_log (
    id              BIGSERIAL PRIMARY KEY,
    timestamp       TIMESTAMPTZ NOT NULL,
    ip              TEXT NOT NULL,
    method          TEXT,
    path            TEXT,
    query           TEXT,
    user_agent      TEXT,
    content_type    TEXT,
    features        TEXT,
    route           TEXT,
    category        TEXT,
    score           DOUBLE PRECISION,
    action          TEXT,
    status_code     INT,
    latency_ms      BIGINT
);
