CREATE TABLE fcm_notifications
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    title      TEXT NOT NULL,
    body       TEXT NOT NULL,
    topic      TEXT NOT NULL,
    is_sent    BOOLEAN NOT NULL DEFAULT FALSE,
    send_at    TIMESTAMP WITH TIME ZONE,
    data       JSONB,
    response   TEXT,
    sender_id  UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_sender_id FOREIGN KEY (sender_id)
    REFERENCES users (id) ON DELETE SET NULL
);
