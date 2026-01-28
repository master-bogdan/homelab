CREATE TABLE IF NOT EXISTS oauth2_clients (
	client_id TEXT PRIMARY KEY,
	client_secret TEXT NOT NULL DEFAULT '',
	redirect_uris TEXT[] NOT NULL,
	grant_types TEXT[] NOT NULL,
	response_types TEXT[] NOT NULL,
	scopes TEXT[] NOT NULL,
	client_name TEXT NOT NULL DEFAULT '',
	client_type TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth2_users (
	user_id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth2_oidc_sessions (
	oidc_session_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES oauth2_users(user_id) ON DELETE CASCADE,
	client_id TEXT NOT NULL REFERENCES oauth2_clients(client_id) ON DELETE CASCADE,
	nonce TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth2_auth_codes (
	auth_code_id TEXT PRIMARY KEY,
	client_id TEXT NOT NULL REFERENCES oauth2_clients(client_id) ON DELETE CASCADE,
	user_id TEXT NOT NULL REFERENCES oauth2_users(user_id) ON DELETE CASCADE,
	oidc_session_id TEXT NOT NULL REFERENCES oauth2_oidc_sessions(oidc_session_id) ON DELETE CASCADE,
	code TEXT NOT NULL UNIQUE,
	redirect_uri TEXT NOT NULL,
	scopes TEXT[] NOT NULL,
	code_challenge TEXT NOT NULL,
	code_challenge_method TEXT NOT NULL,
	is_used BOOLEAN NOT NULL DEFAULT false,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth2_refresh_tokens (
	refresh_token_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES oauth2_users(user_id) ON DELETE CASCADE,
	client_id TEXT NOT NULL REFERENCES oauth2_clients(client_id) ON DELETE CASCADE,
	oidc_session_id TEXT NOT NULL REFERENCES oauth2_oidc_sessions(oidc_session_id) ON DELETE CASCADE,
	scopes TEXT[] NOT NULL,
	token TEXT NOT NULL UNIQUE,
	issued_at TIMESTAMPTZ NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL,
	is_revoked BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth2_access_tokens (
	access_token_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES oauth2_users(user_id) ON DELETE CASCADE,
	client_id TEXT NOT NULL REFERENCES oauth2_clients(client_id) ON DELETE CASCADE,
	oidc_session_id TEXT NOT NULL REFERENCES oauth2_oidc_sessions(oidc_session_id) ON DELETE CASCADE,
	refresh_token_id TEXT NULL REFERENCES oauth2_refresh_tokens(refresh_token_id) ON DELETE SET NULL,
	scopes TEXT[] NOT NULL,
	token TEXT NOT NULL UNIQUE,
	issued_at TIMESTAMPTZ NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL,
	issuer TEXT NOT NULL,
	is_revoked BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_oauth2_auth_codes_code ON oauth2_auth_codes (code);
CREATE INDEX IF NOT EXISTS idx_oauth2_refresh_tokens_token ON oauth2_refresh_tokens (token);
CREATE INDEX IF NOT EXISTS idx_oauth2_access_tokens_token ON oauth2_access_tokens (token);
