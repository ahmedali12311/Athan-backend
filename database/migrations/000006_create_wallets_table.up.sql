BEGIN;

CREATE TABLE wallets
(
    id               UUID PRIMARY KEY NOT NULL UNIQUE,
    credit           NUMERIC(1000, 2) DEFAULT 0,
    trx_total_debit  NUMERIC(1000, 2) DEFAULT 0,
    trx_total_credit NUMERIC(1000, 2) DEFAULT 0,
    trx_count_debit  INTEGER DEFAULT 0,
    trx_count_credit INTEGER DEFAULT 0,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- Foreign keys
    CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE,

    -- Checks
    CONSTRAINT credit_gte_0 CHECK (credit >= 0),
    CONSTRAINT trx_total_debit_gte_0 CHECK (trx_total_debit >= 0),
    CONSTRAINT trx_total_credit_gte_0 CHECK (trx_total_credit >= 0),
    CONSTRAINT trx_count_debit_gte_0 CHECK (trx_count_debit >= 0),
    CONSTRAINT trx_count_credit_gte_0 CHECK (trx_count_credit >= 0)
);


--trigger: update_update_at
CREATE TRIGGER app_trigger_update_wallets_updated_at
BEFORE UPDATE ON wallets FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();

COMMIT;
