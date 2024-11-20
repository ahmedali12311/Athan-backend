CREATE TABLE fcm_notifications
(
    id         uuid PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    title      character varying(255)   NOT NULL,
    body       character varying(255)   NOT NULL,
    topic      character varying(255)   NOT NULL,
    is_sent    boolean                  NOT NULL DEFAULT false,
    send_at    timestamp with time zone,
    data       JSONB,
    response   text,
    sender_id  uuid,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),

    CONSTRAINT fk_sender_id
        FOREIGN KEY (sender_id)
            REFERENCES users (id) ON DELETE SET NULL
);

--trigger: update_update_at
CREATE TRIGGER app_trigger_update_fcm_notifications_updated_at
    BEFORE UPDATE
    ON fcm_notifications
    FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
