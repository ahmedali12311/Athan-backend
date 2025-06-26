CREATE TABLE settings (
    id                   SERIAL PRIMARY KEY,
    name                 JSONB NOT NULL,
    key                  TEXT NOT NULL UNIQUE,
    value                TEXT,
    sort                 INTEGER NOT NULL DEFAULT 0,
    is_disabled          BOOLEAN DEFAULT TRUE,
    is_readonly          BOOLEAN DEFAULT FALSE,
    field_type           TEXT DEFAULT 'text',
    data_type            TEXT DEFAULT 'string',
    category_id          UUID,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_category_id FOREIGN KEY (category_id)
    REFERENCES categories (id) ON DELETE CASCADE,

    -- reflective of html form field types
    CONSTRAINT settings_field_type_check CHECK (field_type = any(
        ARRAY[
            'text'::TEXT,
            'number'::TEXT,
            'email'::TEXT,
            'password'::TEXT,
            'tel'::TEXT,
            'date'::TEXT,
            'color'::TEXT,
            'url'::TEXT,
            'textarea'::TEXT,
            'toggle'::TEXT,
            'file'::TEXT
        ]
    )),

    -- reflective of typescript variable types
    CONSTRAINT settings_datatype_check CHECK (data_type = any(
        ARRAY[
            'boolean'::TEXT,
            'number'::TEXT,
            'string'::TEXT
        ]
    ))
);
