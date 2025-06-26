CREATE TABLE tokens (
    id          SERIAL PRIMARY KEY,
    user_id     UUID NOT NULL,
    token_type  TEXT,
    token_value TEXT,
    data        JSONB,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);
