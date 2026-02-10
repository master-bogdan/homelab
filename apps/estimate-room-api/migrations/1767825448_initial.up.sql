
CREATE TYPE "pkce_challenge_method" AS ENUM (
  'PLAIN',
  'S256'
);

CREATE TYPE "team_member_role" AS ENUM (
  'OWNER',
  'MEMBER'
);

CREATE TYPE "invitation_type" AS ENUM (
  'TEAM',
  'ROOM'
);

CREATE TYPE "invitation_status" AS ENUM (
  'PENDING',
  'ACCEPTED',
  'DECLINED',
  'EXPIRED'
);

CREATE TYPE "deck_type" AS ENUM (
  'FIBONACCI',
  'TSHIRT',
  'CUSTOM'
);

CREATE TYPE "room_status" AS ENUM (
  'ACTIVE',
  'FINISHED',
  'EXPIRED'
);

CREATE TYPE "room_participant_role" AS ENUM (
  'ADMIN',
  'VOTER',
  'SPECTATOR'
);

CREATE TYPE "task_status" AS ENUM (
  'PENDING',
  'VOTING',
  'ESTIMATED',
  'SKIPPED'
);

CREATE TABLE "users" (
  "user_id" text PRIMARY KEY,
  "email" text UNIQUE,
  "password_hash" text,
  "github_id" text,
  "display_name" text NOT NULL DEFAULT '',
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
  "role" team_member_role NOT NULL DEFAULT 'member',
  "joined_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("team_id", "user_id")
);

CREATE TABLE "invitations" (
  "invintation_id" text PRIMARY KEY,
  "type" invitation_type NOT NULL,
  "email" text NOT NULL,
  "to_user_id" text,
  "team_id" text,
  "room_id" text,
  "token" text UNIQUE NOT NULL,
  "status" invitation_status NOT NULL DEFAULT 'pending',
  "expires_at" timestamptz NOT NULL,
  "created_by_user_id" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "decks" (
  "deck_id" text PRIMARY KEY,
  "name" text NOT NULL,
  "type" deck_type NOT NULL,
  "values" jsonb NOT NULL
);

CREATE TABLE "rooms" (
  "room_id" text PRIMARY KEY,
  "code" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "admin_user_id" text NOT NULL,
  "team_id" text,
  "deck_id" text NOT NULL,
  "status" room_status NOT NULL DEFAULT 'active',
  "allow_guests" boolean NOT NULL DEFAULT false,
  "allow_spectators" boolean NOT NULL DEFAULT false,
  "round_timer_seconds" int NOT NULL DEFAULT 120,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_activity_at" timestamptz NOT NULL DEFAULT (now()),
  "finished_at" timestamptz
);

CREATE TABLE "room_participants" (
  "room_participants_id" text PRIMARY KEY,
  "room_id" text NOT NULL,
  "user_id" text,
  "guest_name" text,
  "role" room_participant_role NOT NULL DEFAULT 'voter',
  "joined_at" timestamptz NOT NULL DEFAULT (now()),
  "left_at" timestamptz
);

CREATE TABLE "tasks" (
  "task_id" text PRIMARY KEY,
  "room_id" text NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "external_key" text,
  "status" task_status NOT NULL DEFAULT 'pending',
  "final_estimate_value" text,
  "order_index" int NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
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
  "created_at" timestamptz NOT NULL DEFAULT (now())
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

CREATE UNIQUE INDEX ON "votes" ("task_id", "participant_id", "round_number");

CREATE INDEX "idx_oauth2_auth_codes_code" ON "oauth2_auth_codes" ("code");

CREATE INDEX "idx_oauth2_refresh_tokens_token" ON "oauth2_refresh_tokens" ("token");

CREATE INDEX "idx_oauth2_access_tokens_token" ON "oauth2_access_tokens" ("token");

COMMENT ON COLUMN "decks"."values" IS 'JSONB array of strings';

COMMENT ON COLUMN "rooms"."code" IS 'short unique string used in URLs';

COMMENT ON COLUMN "room_participants"."user_id" IS 'nullable for guests';

COMMENT ON COLUMN "votes"."value" IS 'must be in deck values';

COMMENT ON COLUMN "oauth2_clients"."redirect_uris" IS 'PostgreSQL text[]';

COMMENT ON COLUMN "oauth2_access_tokens"."refresh_token_id" IS 'ON DELETE SET NULL';

ALTER TABLE "user_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "user_settings" ADD FOREIGN KEY ("default_deck_id") REFERENCES "decks" ("deck_id");

ALTER TABLE "teams" ADD FOREIGN KEY ("owner_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "team_members" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "team_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "invitations" ADD FOREIGN KEY ("created_by_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("admin_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id");

ALTER TABLE "rooms" ADD FOREIGN KEY ("deck_id") REFERENCES "decks" ("deck_id");

ALTER TABLE "room_participants" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "room_participants" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("room_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("task_id");

ALTER TABLE "votes" ADD FOREIGN KEY ("participant_id") REFERENCES "room_participants" ("room_participants_id");

ALTER TABLE "user_stats" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "user_achievements" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

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

ALTER TABLE "oauth2_access_tokens" ADD FOREIGN KEY ("refresh_token_id") REFERENCES "oauth2_refresh_tokens" ("refresh_token_id");
