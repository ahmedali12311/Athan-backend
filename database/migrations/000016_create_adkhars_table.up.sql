CREATE TABLE adhkars (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    source VARCHAR(255) NOT NULL,
    repeat INTEGER NOT NULL DEFAULT 1,
    category_id  UUID  NOT NULL REFERENCES categories(id),
        created_by_id uuid REFERENCES users(id),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TRIGGER app_trigger_update_adhkars_updated_at
BEFORE UPDATE ON adhkars
FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();