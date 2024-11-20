CREATE TABLE wallet_transactions
(
    id                UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    wallet_id         UUID NOT NULL,
    amount            NUMERIC(1000, 2) NOT NULL DEFAULT 0,
    type              TEXT,
    notes             TEXT,
    payment_method    TEXT,
    payment_reference TEXT,
    recharged_by_id   UUID,
    is_confirmed      BOOLEAN NOT NULL DEFAULT FALSE,
    tlync_url         TEXT,
    tlync_response    JSONB,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CHECK ((type = any(ARRAY[
        'credit'::TEXT,
        'debit'::TEXT
    ]))),

    CONSTRAINT fk_wallet_id FOREIGN KEY (wallet_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_recharged_by_id FOREIGN KEY (recharged_by_id) REFERENCES users (id) ON DELETE RESTRICT
);

CREATE TRIGGER app_trigger_update_wallet_transactions_updated_at
BEFORE UPDATE ON wallet_transactions FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
