Project estimate_room_v1 {
  database_type: 'PostgreSQL'
}

// ---------- Enums ----------
Enum pkce_challenge_method {
  PLAIN
  S256
}

Enum team_member_role {
  OWNER
  MEMBER
}

Enum room_status {
  ACTIVE
  FINISHED
  EXPIRED
}

Enum room_participant_role {
  ADMIN
  MEMBER
  GUEST
}

Enum task_status {
  PENDING
  VOTING
  ESTIMATED
  SKIPPED
}

Enum round_status {
  ACTIVE
  REVEALED
}

// ---------- Core ----------
Table users {
  user_id        text        [pk]
  email          text        [unique]
  password_hash  text
  github_id      text
  display_name   text        [not null, default: '']
  avatar_url     text
  created_at     timestamptz [not null, default: `now()`]
  updated_at     timestamptz [not null, default: `now()`]
  last_login_at  timestamptz
  deleted_at     timestamptz
}

Table user_settings {
  user_id              text  [pk, ref: > users.user_id]
  theme                text
  timezone             text
  locale               text
  default_deck_id      text [ref: > decks.deck_id]
  default_room_options jsonb
}

Table teams {
  team_id             text        [pk]
  name           text        [not null]
  owner_user_id  text        [not null, ref: > users.user_id]
  created_at     timestamptz [not null, default: `now()`]
}

Table team_members {
  team_id   text             [not null, ref: > teams.team_id]
  user_id   text             [not null, ref: > users.user_id]
  role      team_member_role [not null, default: 'MEMBER']
  joined_at timestamptz      [not null, default: `now()`]

  Indexes {
    (team_id, user_id) [pk] // composite PK :contentReference[oaicite:0]{index=0}
  }
}

Table decks {
  deck_id    text        [pk]
  name       text        [not null]
  kind       text        [not null]
  values     jsonb       [not null, note: 'JSONB array of strings']
  created_at timestamptz [not null, default: `now()`]
}

Table rooms {
  room_id                 text        [pk, default: `gen_random_uuid()`]
  code               text        [not null, unique, note: 'short unique string used in URLs']
  name               text        [not null]
  admin_user_id      text        [not null, ref: > users.user_id]
  team_id            text        [ref: > teams.team_id]
  deck               jsonb       [not null, note: 'Object with name, kind, values[]']
  status             room_status [not null, default: 'ACTIVE']
  created_at         timestamptz [not null, default: `now()`]
  last_activity_at   timestamptz [not null, default: `now()`]
  finished_at        timestamptz

  Indexes {
    (last_activity_at) [name: 'rooms_active_last_activity_idx', note: 'partial index where status = ACTIVE']
  }
}

Table room_participants {
  room_participants_id         text                 [pk]
  room_id    text                 [not null, ref: > rooms.room_id]
  user_id    text                 [ref: > users.user_id, note: 'nullable for guests']
  guest_name text
  role       room_participant_role [not null, default: 'MEMBER']
  joined_at  timestamptz          [not null, default: `now()`]
  left_at    timestamptz
}

Table tasks {
  task_id                  text        [pk]
  room_id             text        [not null, ref: > rooms.room_id]
  title               text        [not null]
  description         text
  external_key        text
  status              task_status [not null, default: 'PENDING']
  is_active           boolean     [not null, default: false]
  final_estimate_value text
  created_at          timestamptz [not null, default: `now()`]
  updated_at          timestamptz [not null, default: `now()`]

  Indexes {
    (room_id) [name: 'tasks_one_active_per_room_idx', unique, note: 'partial unique index where is_active = true']
  }
}

Table task_rounds {
  task_id       text        [not null, ref: > tasks.task_id]
  round_number  int         [not null]
  eligible_participant_ids jsonb [not null, default: '[]', note: 'JSONB array of participant ids eligible for the round']
  status        round_status [not null, default: 'ACTIVE']
  created_at    timestamptz [not null, default: `now()`]
  updated_at    timestamptz [not null, default: `now()`]

  Indexes {
    (task_id, round_number) [pk]
  }
}

Table votes {
  votes_id             text        [pk]
  task_id        text        [not null, ref: > tasks.task_id]
  participant_id text        [not null, ref: > room_participants.room_participants_id]
  value          text        [not null, note: 'must be in deck values']
  round_number   int         [not null, default: 1]
  created_at     timestamptz [not null, default: `now()`]

  Indexes {
    (task_id, participant_id, round_number) [unique]
  }
}

// ---------- Gamification ----------
Table user_stats {
  user_id               text [pk, ref: > users.user_id]
  sessions_participated int  [not null, default: 0]
  sessions_admined      int  [not null, default: 0]
  tasks_estimated       int  [not null, default: 0]
  xp                    int  [not null, default: 0]
}

Table user_achievements {
  user_id         text        [not null, ref: > users.user_id]
  achievement_key text        [not null]
  level           int         [not null, default: 1]
  unlocked_at     timestamptz [not null, default: `now()`]

  Indexes {
    (user_id, achievement_key) [pk]
  }
}

Table team_stats {
  team_id         text [pk, ref: > teams.team_id]
  sessions_total  int  [not null, default: 0]
  tasks_estimated int  [not null, default: 0]
  xp              int  [not null, default: 0]
}

Table team_achievements {
  team_id         text        [not null, ref: > teams.team_id]
  achievement_key text        [not null]
  level           int         [not null, default: 1]
  unlocked_at     timestamptz [not null, default: `now()`]

  Indexes {
    (team_id, achievement_key) [pk]
  }
}

// ---------- OAuth2 / OIDC ----------
Table oauth2_clients {
  client_id      text        [pk]
  client_secret  text        [not null, default: '']
  redirect_uris  text[]      [not null, note: 'PostgreSQL text[]'] // arrays as type tokens :contentReference[oaicite:1]{index=1}
  grant_types    text[]      [not null]
  response_types text[]      [not null]
  scopes         text[]      [not null]
  client_name    text        [not null, default: '']
  client_type    text        [not null, default: '']
  created_at     timestamptz [not null, default: `now()`]
}

Table oauth2_oidc_sessions {
  oidc_session_id text        [pk]
  user_id         text        [not null, ref: > users.user_id]
  client_id       text        [not null, ref: > oauth2_clients.client_id]
  nonce           text        [not null]
  created_at      timestamptz [not null, default: `now()`]
}

Table oauth2_auth_codes {
  auth_code_id           text                  [pk]
  client_id              text                  [not null, ref: > oauth2_clients.client_id]
  user_id                text                  [not null, ref: > users.user_id]
  oidc_session_id         text                 [not null, ref: > oauth2_oidc_sessions.oidc_session_id]
  code                   text                  [not null, unique]
  redirect_uri           text                  [not null]
  scopes                 text[]                [not null]
  code_challenge         text                  [not null]
  code_challenge_method  pkce_challenge_method [not null]
  is_used                boolean               [not null, default: false]
  expires_at             timestamptz           [not null]
  created_at             timestamptz           [not null, default: `now()`]

  Indexes {
    (code) [name: 'idx_oauth2_auth_codes_code'] // index syntax :contentReference[oaicite:2]{index=2}
  }
}

Table oauth2_refresh_tokens {
  refresh_token_id text        [pk]
  user_id          text        [not null, ref: > users.user_id]
  client_id        text        [not null, ref: > oauth2_clients.client_id]
  oidc_session_id  text        [not null, ref: > oauth2_oidc_sessions.oidc_session_id]
  scopes           text[]      [not null]
  token            text        [not null, unique]
  issued_at        timestamptz [not null]
  expires_at       timestamptz [not null]
  is_revoked       boolean     [not null, default: false]
  created_at       timestamptz [not null, default: `now()`]

  Indexes {
    (token) [name: 'idx_oauth2_refresh_tokens_token']
  }
}

Table oauth2_access_tokens {
  access_token_id  text        [pk]
  user_id          text        [not null, ref: > users.user_id]
  client_id        text        [not null, ref: > oauth2_clients.client_id]
  oidc_session_id  text        [not null, ref: > oauth2_oidc_sessions.oidc_session_id]
  refresh_token_id text        [ref: > oauth2_refresh_tokens.refresh_token_id, note: 'ON DELETE SET NULL']
  scopes           text[]      [not null]
  token            text        [not null, unique]
  issued_at        timestamptz [not null]
  expires_at       timestamptz [not null]
  issuer           text        [not null]
  is_revoked       boolean     [not null, default: false]
  created_at       timestamptz [not null, default: `now()`]

  Indexes {
    (token) [name: 'idx_oauth2_access_tokens_token']
  }
}
