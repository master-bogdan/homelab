CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE "pkce_challenge_method" AS ENUM (
  'PLAIN',
  'S256'
);

CREATE TYPE "team_member_role" AS ENUM (
  'OWNER',
  'MEMBER'
);

CREATE TYPE "invitation_kind" AS ENUM (
  'TEAM_MEMBER',
  'ROOM_EMAIL',
  'ROOM_LINK'
);

CREATE TYPE "invitation_status" AS ENUM (
  'ACTIVE',
  'ACCEPTED',
  'DECLINED',
  'REVOKED'
);

CREATE TYPE "room_status" AS ENUM (
  'ACTIVE',
  'FINISHED',
  'EXPIRED'
);

CREATE TYPE "room_participant_role" AS ENUM (
  'ADMIN',
  'MEMBER',
  'GUEST'
);

CREATE TYPE "task_status" AS ENUM (
  'PENDING',
  'VOTING',
  'ESTIMATED',
  'SKIPPED'
);

CREATE TYPE "round_status" AS ENUM (
  'ACTIVE',
  'REVEALED'
);

CREATE TABLE "users" (
  "user_id" text PRIMARY KEY,
  "email" text UNIQUE,
  "password_hash" text,
  "github_id" text,
  "display_name" text NOT NULL DEFAULT '',
  "organization" text,
  "occupation" text,
  "avatar_url" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "user_settings" (
  "user_id" text PRIMARY KEY,
  "theme" text,
  "timezone" text,
  "locale" text,
  "default_deck_id" text,
  "default_room_options" jsonb
);

CREATE TABLE "teams" (
  "team_id" text PRIMARY KEY,
  "name" text NOT NULL,
  "owner_user_id" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "team_members" (
  "team_id" text NOT NULL,
  "user_id" text NOT NULL,
  "role" team_member_role NOT NULL DEFAULT 'MEMBER',
  "joined_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("team_id", "user_id")
);

CREATE TABLE "decks" (
  "deck_id" text PRIMARY KEY,
  "name" text NOT NULL,
  "kind" text NOT NULL,
  "values" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "rooms" (
  "room_id" text PRIMARY KEY DEFAULT (gen_random_uuid()::text),
  "code" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "admin_user_id" text NOT NULL,
  "team_id" text,
  "deck" jsonb NOT NULL,
  "status" room_status NOT NULL DEFAULT 'ACTIVE',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_activity_at" timestamptz NOT NULL DEFAULT (now()),
  "finished_at" timestamptz
);

CREATE TABLE "room_participants" (
  "room_participants_id" text PRIMARY KEY,
  "room_id" text NOT NULL,
  "user_id" text,
  "guest_name" text,
  "role" room_participant_role NOT NULL DEFAULT 'MEMBER',
  "joined_at" timestamptz NOT NULL DEFAULT (now()),
  "left_at" timestamptz
);

CREATE TABLE "invitations" (
  "invitation_id" text PRIMARY KEY,
  "kind" invitation_kind NOT NULL,
  "status" invitation_status NOT NULL DEFAULT 'ACTIVE',
  "team_id" text,
  "room_id" text,
  "invited_user_id" text,
  "invited_email" text,
  "created_by_user_id" text NOT NULL,
  "token_id" text UNIQUE NOT NULL,
  "accepted_at" timestamptz,
  "declined_at" timestamptz,
  "revoked_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CHECK (
    (
      "kind" = 'TEAM_MEMBER' AND
      "team_id" IS NOT NULL AND
      "room_id" IS NULL AND
      "invited_user_id" IS NOT NULL AND
      "invited_email" IS NOT NULL
    ) OR (
      "kind" = 'ROOM_EMAIL' AND
      "team_id" IS NULL AND
      "room_id" IS NOT NULL AND
      "invited_email" IS NOT NULL
    ) OR (
      "kind" = 'ROOM_LINK' AND
      "team_id" IS NULL AND
      "room_id" IS NOT NULL AND
      "invited_user_id" IS NULL AND
      "invited_email" IS NULL
    )
  )
);

CREATE TABLE "tasks" (
  "task_id" text PRIMARY KEY,
  "room_id" text NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "external_key" text,
  "status" task_status NOT NULL DEFAULT 'PENDING',
  "is_active" boolean NOT NULL DEFAULT false,
  "final_estimate_value" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "task_rounds" (
  "task_id" text NOT NULL,
  "round_number" int NOT NULL,
  "eligible_participant_ids" jsonb NOT NULL DEFAULT '[]'::jsonb,
  "status" round_status NOT NULL DEFAULT 'ACTIVE',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("task_id", "round_number")
);

CREATE TABLE "votes" (
  "votes_id" text PRIMARY KEY,
  "task_id" text NOT NULL,
  "participant_id" text NOT NULL,
  "value" text NOT NULL,
  "round_number" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_stats" (
  "user_id" text PRIMARY KEY,
  "sessions_participated" int NOT NULL DEFAULT 0,
  "sessions_admined" int NOT NULL DEFAULT 0,
  "tasks_estimated" int NOT NULL DEFAULT 0,
  "xp" int NOT NULL DEFAULT 0
);

CREATE TABLE "user_achievements" (
  "user_id" text NOT NULL,
  "achievement_key" text NOT NULL,
  "level" int NOT NULL DEFAULT 1,
  "unlocked_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "achievement_key")
);

CREATE TABLE "user_session_rewards" (
  "room_id" text NOT NULL,
  "user_id" text NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "sessions_participated_delta" int NOT NULL DEFAULT 0,
  "sessions_admined_delta" int NOT NULL DEFAULT 0,
  "tasks_estimated_delta" int NOT NULL DEFAULT 0,
  "xp_gained" int NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("room_id", "user_id")
);

CREATE TABLE "team_stats" (
  "team_id" text PRIMARY KEY,
  "sessions_total" int NOT NULL DEFAULT 0,
  "tasks_estimated" int NOT NULL DEFAULT 0,
  "xp" int NOT NULL DEFAULT 0
);

CREATE TABLE "team_achievements" (
  "team_id" text NOT NULL,
  "achievement_key" text NOT NULL,
  "level" int NOT NULL DEFAULT 1,
  "unlocked_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("team_id", "achievement_key")
);

CREATE TABLE "oauth2_clients" (
  "client_id" text PRIMARY KEY,
  "client_secret" text NOT NULL DEFAULT '',
  "redirect_uris" text[] NOT NULL,
  "grant_types" text[] NOT NULL,
  "response_types" text[] NOT NULL,
  "scopes" text[] NOT NULL,
  "client_name" text NOT NULL DEFAULT '',
  "client_type" text NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "oauth2_oidc_sessions" (
  "oidc_session_id" text PRIMARY KEY,
  "user_id" text NOT NULL,
  "client_id" text NOT NULL,
  "nonce" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "revoked_at" timestamptz
);

CREATE TABLE "oauth2_auth_codes" (
  "auth_code_id" text PRIMARY KEY,
  "client_id" text NOT NULL,
  "user_id" text NOT NULL,
  "oidc_session_id" text NOT NULL,
  "code" text UNIQUE NOT NULL,
  "redirect_uri" text NOT NULL,
  "scopes" text[] NOT NULL,
  "code_challenge" text NOT NULL,
  "code_challenge_method" pkce_challenge_method NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "oauth2_refresh_tokens" (
  "refresh_token_id" text PRIMARY KEY,
  "user_id" text NOT NULL,
  "client_id" text NOT NULL,
  "oidc_session_id" text NOT NULL,
  "scopes" text[] NOT NULL,
  "token" text UNIQUE NOT NULL,
  "issued_at" timestamptz NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_revoked" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "oauth2_access_tokens" (
  "access_token_id" text PRIMARY KEY,
  "user_id" text NOT NULL,
  "client_id" text NOT NULL,
  "oidc_session_id" text NOT NULL,
  "refresh_token_id" text,
  "scopes" text[] NOT NULL,
  "token" text UNIQUE NOT NULL,
  "issued_at" timestamptz NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "issuer" text NOT NULL,
  "is_revoked" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "auth_password_reset_tokens" (
  "password_reset_token_id" text PRIMARY KEY,
  "user_id" text NOT NULL,
  "token_hash" text UNIQUE NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

INSERT INTO "decks" ("deck_id", "name", "kind", "values")
VALUES
  (
    'FIBONACCI',
    'Fibonacci',
    'FIBONACCI',
    '["0","1","2","3","5","8","13","21","34","55","?"]'::jsonb
  ),
  (
    'TSHIRT',
    'T-Shirt',
    'TSHIRT',
    '["XS","S","M","L","XL","XXL","?"]'::jsonb
  );

CREATE UNIQUE INDEX ON "votes" ("task_id", "participant_id", "round_number");
CREATE UNIQUE INDEX "tasks_one_active_per_room_idx" ON "tasks" ("room_id") WHERE "is_active" = true;
CREATE INDEX "rooms_active_last_activity_idx" ON "rooms" ("last_activity_at") WHERE "status" = 'ACTIVE';
CREATE INDEX "rooms_team_id_idx" ON "rooms" ("team_id");
CREATE INDEX "invitations_room_id_idx" ON "invitations" ("room_id");
CREATE INDEX "invitations_team_id_idx" ON "invitations" ("team_id");
CREATE UNIQUE INDEX "invitations_active_team_member_unique_idx"
  ON "invitations" ("team_id", "invited_user_id")
  WHERE "kind" = 'TEAM_MEMBER' AND "status" = 'ACTIVE';
CREATE INDEX "invitations_active_invited_email_idx" ON "invitations" ("invited_email") WHERE "status" = 'ACTIVE' AND "invited_email" IS NOT NULL;
CREATE INDEX "user_session_rewards_user_id_idx" ON "user_session_rewards" ("user_id");

CREATE INDEX "idx_oauth2_auth_codes_code" ON "oauth2_auth_codes" ("code");

CREATE INDEX "idx_oauth2_refresh_tokens_token" ON "oauth2_refresh_tokens" ("token");

CREATE INDEX "idx_oauth2_access_tokens_token" ON "oauth2_access_tokens" ("token");
CREATE INDEX "idx_oauth2_oidc_sessions_active_user_id" ON "oauth2_oidc_sessions" ("user_id") WHERE "revoked_at" IS NULL;
CREATE INDEX "idx_auth_password_reset_tokens_user_id" ON "auth_password_reset_tokens" ("user_id");

COMMENT ON COLUMN "rooms"."code" IS 'short unique string used in URLs';

COMMENT ON COLUMN "room_participants"."user_id" IS 'nullable for guests';

COMMENT ON COLUMN "invitations"."token_id" IS 'token reference embedded in signed invite claims';

COMMENT ON COLUMN "votes"."value" IS 'must be in deck values';

COMMENT ON COLUMN "oauth2_clients"."redirect_uris" IS 'PostgreSQL text[]';

COMMENT ON COLUMN "oauth2_access_tokens"."refresh_token_id" IS 'ON DELETE SET NULL';

ALTER TABLE "user_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "user_settings" ADD FOREIGN KEY ("default_deck_id") REFERENCES "decks" ("deck_id");

ALTER TABLE "teams" ADD FOREIGN KEY ("owner_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "team_members" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "team_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("admin_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "room_participants" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "room_participants" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("invited_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("created_by_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "task_rounds" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("task_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("task_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("participant_id") REFERENCES "room_participants" ("room_participants_id");

ALTER TABLE "user_stats" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "user_achievements" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "user_session_rewards" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "user_session_rewards" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "team_stats" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "team_achievements" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "oauth2_oidc_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "oauth2_oidc_sessions" ADD FOREIGN KEY ("client_id") REFERENCES "oauth2_clients" ("client_id");

ALTER TABLE "oauth2_auth_codes" ADD FOREIGN KEY ("client_id") REFERENCES "oauth2_clients" ("client_id");

ALTER TABLE "oauth2_auth_codes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "oauth2_auth_codes" ADD FOREIGN KEY ("oidc_session_id") REFERENCES "oauth2_oidc_sessions" ("oidc_session_id");

ALTER TABLE "oauth2_refresh_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "oauth2_refresh_tokens" ADD FOREIGN KEY ("client_id") REFERENCES "oauth2_clients" ("client_id");

ALTER TABLE "oauth2_refresh_tokens" ADD FOREIGN KEY ("oidc_session_id") REFERENCES "oauth2_oidc_sessions" ("oidc_session_id");

ALTER TABLE "oauth2_access_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "oauth2_access_tokens" ADD FOREIGN KEY ("client_id") REFERENCES "oauth2_clients" ("client_id");

ALTER TABLE "oauth2_access_tokens" ADD FOREIGN KEY ("oidc_session_id") REFERENCES "oauth2_oidc_sessions" ("oidc_session_id");

ALTER TABLE "oauth2_access_tokens" ADD FOREIGN KEY ("refresh_token_id") REFERENCES "oauth2_refresh_tokens" ("refresh_token_id") ON DELETE SET NULL;

ALTER TABLE "auth_password_reset_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
