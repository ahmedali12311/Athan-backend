CREATE TABLE user_notifications (
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    is_read     BOOLEAN NOT NULL DEFAULT FALSE,
    is_notified BOOLEAN NOT NULL DEFAULT FALSE,
    title       TEXT,
    body        TEXT,
    response    TEXT,
    data        JSONB,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);
