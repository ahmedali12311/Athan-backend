CREATE TABLE special_topics (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    category_id  UUID  REFERENCES categories(id),
    created_by_id uuid REFERENCES users(id),
    img           TEXT,
    thumb         TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TRIGGER app_trigger_update_special_topics_updated_at
BEFORE UPDATE ON special_topics
FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();

