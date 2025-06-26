CREATE TABLE users (
    id            UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    ref           TEXT NOT NULL,
    name          TEXT,
    phone         TEXT,
    email         TEXT,
    password_hash BYTEA,
    img           TEXT,
    thumb         TEXT,
    gender        TEXT,
    details       TEXT,
    birthdate     DATE,
    is_anon       BOOLEAN NOT NULL DEFAULT FALSE,
    is_notifiable BOOLEAN NOT NULL DEFAULT TRUE,
    is_disabled   BOOLEAN NOT NULL DEFAULT FALSE,
    is_confirmed  BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted    BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    last_ref      INTEGER NOT NULL DEFAULT 0,

    pin           TEXT,
    pin_expiry    TIMESTAMP WITH TIME ZONE,

    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    location      GEOMETRY (POINT, 4326)
);

ALTER TABLE users ADD CHECK ((gender = any(ARRAY[
    'male'::TEXT,
    'female'::TEXT
])));

ALTER TABLE users ADD CONSTRAINT phone_or_email_or_anon CHECK ((
    (phone IS NOT NULL)::INTEGER +
    (email IS NOT NULL)::INTEGER +
    (is_deleted IS NOT FALSE)::INTEGER +
    (is_anon IS NOT FALSE)::INTEGER
) >= 1);

CREATE UNIQUE INDEX users_email_unique_nullable
ON users (email)
WHERE email IS NOT NULL;

CREATE UNIQUE INDEX users_phone_unique_nullable
ON users (phone)
WHERE phone IS NOT NULL;

-- spatial indices
CREATE INDEX users_location_idx
ON users USING gist (location);
