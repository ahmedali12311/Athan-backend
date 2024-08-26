CREATE TABLE settings (
    id                   SERIAL PRIMARY KEY,
    key                  TEXT NOT NULL UNIQUE,
    value                TEXT,
    is_disabled          BOOLEAN DEFAULT TRUE,
    is_readonly          BOOLEAN DEFAULT FALSE,
    field_type           TEXT DEFAULT 'text',
    data_type            TEXT DEFAULT 'string',
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

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

-- trigger: update_update_at
CREATE TRIGGER app_trigger_update_settings_updated_at
BEFORE UPDATE ON settings FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();
