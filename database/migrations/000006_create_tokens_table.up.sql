CREATE TABLE tokens (
    id          SERIAL PRIMARY KEY,
    user_id     UUID NOT NULL,
    token_type  TEXT,
    token_value TEXT,
    data        JSONB,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

-- trigger: update_update_at
CREATE TRIGGER app_trigger_update_tokens_updated_at
BEFORE UPDATE ON tokens FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
