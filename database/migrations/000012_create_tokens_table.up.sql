CREATE TABLE tokens (
    id            SERIAL PRIMARY KEY,
    user_id       UUID,
    city_id       UUID REFERENCES Cities(id) ON DELETE SET NULL,
    token_type    TEXT,
    token_value   TEXT NOT NULL UNIQUE,
    data          JSONB,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);