CREATE TABLE users (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE walk_options (
    id INT PRIMARY KEY,
    category TEXT NOT NULL CHECK (category IN ('free', 'minutes', 'destination')),
    title TEXT NOT NULL
);

CREATE TABLE walks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    walk_option_id INT REFERENCES walk_options(id),
    status TEXT NOT NULL CHECK (status IN ('walking', 'completed', 'canceled')),
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE characters (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    supporter_type TEXT NOT NULL CHECK (
        supporter_type IN ('dog', 'cat')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id)
);

CREATE TABLE auth_identities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    provider TEXT NOT NULL CHECK (provider IN ('google')),
    provider_user_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (provider, provider_user_id),
    UNIQUE (user_id, provider)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);