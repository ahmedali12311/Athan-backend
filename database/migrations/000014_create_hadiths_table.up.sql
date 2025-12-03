CREATE TABLE hadiths (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    source VARCHAR(255) NOT NULL,
    topic VARCHAR(255) NOT NULL,
    category_id  UUID  NOT NULL REFERENCES categories(id),
    created_by_id uuid REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
      
);

CREATE INDEX idx_hadiths_topic ON hadiths(topic);

CREATE TRIGGER app_trigger_update_hadiths_updated_at
BEFORE UPDATE ON hadiths
FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();

