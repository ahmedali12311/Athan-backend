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

    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

-- trigger: update_update_at
CREATE TRIGGER app_trigger_update_user_notifications_updated_at
BEFORE UPDATE ON user_notifications FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
