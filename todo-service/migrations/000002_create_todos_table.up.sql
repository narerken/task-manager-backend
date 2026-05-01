CREATE TABLE IF NOT EXISTS todos (
                                     id          BIGSERIAL PRIMARY KEY,
                                     title       VARCHAR(255) NOT NULL,
    description TEXT,
    completed   BOOLEAN      NOT NULL DEFAULT FALSE,
    priority    INTEGER      NOT NULL DEFAULT 1,
    due_date    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );