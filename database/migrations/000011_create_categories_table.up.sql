CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            JSONB NOT NULL,
    depth           INTEGER NOT NULL DEFAULT 0,
    sort            INTEGER NOT NULL DEFAULT 0,
    is_disabled     BOOLEAN NOT NULL DEFAULT FALSE,
    is_featured     BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id       UUID,
    super_parent_id UUID,
    img             TEXT,
    thumb           TEXT,
    path            UUID [] DEFAULT '{}',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_parent_id
    FOREIGN KEY (parent_id)
    REFERENCES categories (id)
    ON DELETE CASCADE,

    CONSTRAINT fk_super_parent_id
    FOREIGN KEY (super_parent_id)
    REFERENCES categories (id)
    ON DELETE CASCADE,

    CONSTRAINT no_parent_without_super_parent CHECK (
        NOT (
            parent_id IS NOT NULL AND
            super_parent_id IS NULL
        )
    )
);

CREATE UNIQUE INDEX unique_name_parent_idx
ON categories (name, parent_id)
WHERE parent_id IS NOT NULL;

CREATE UNIQUE INDEX unique_name_null_parent_idx
ON categories (name)
WHERE parent_id IS NULL;

CREATE INDEX idx_categories_path ON categories USING gin (path);

-- create auto update function
CREATE OR REPLACE FUNCTION update_category_path() RETURNS TRIGGER AS $$
BEGIN
    NEW.path = ARRAY[NEW.id];

    IF NEW.parent_id IS NOT NULL THEN
        SELECT path || NEW.id INTO NEW.path 
        FROM categories 
        WHERE id = NEW.parent_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger auto update for path
CREATE TRIGGER trigger_update_category_path
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION update_category_path();
