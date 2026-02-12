DROP TABLE IF EXISTS oauth2_access_tokens;
DROP TABLE IF EXISTS oauth2_refresh_tokens;
DROP TABLE IF EXISTS oauth2_auth_codes;
DROP TABLE IF EXISTS oauth2_oidc_sessions;
DROP TABLE IF EXISTS oauth2_clients;

DROP TABLE IF EXISTS team_achievements;
DROP TABLE IF EXISTS team_stats;
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS room_participants;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS decks;
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS room_participant_role;
DROP TYPE IF EXISTS room_status;
DROP TYPE IF EXISTS deck_type;
DROP TYPE IF EXISTS invitation_status;
DROP TYPE IF EXISTS invitation_type;
DROP TYPE IF EXISTS team_member_role;
DROP TYPE IF EXISTS pkce_challenge_method;

DROP EXTENSION IF EXISTS "pgcrypto";
