CREATE DATABASE db;
ALTER DATABASE db SET TIMEZONE = 'UTC';

CREATE TABLE invites (
    id uuid primary key
);

CREATE TABLE users (
    id uuid primary key,
    login text NOT NULL UNIQUE,
    password text NOT NULL,
    secret bytea NOT NULL
);

-- CREATE TABLE users_teams (
--     id uuid primary key,
--     user_id uuid not null references users (id) on delete cascade,
--     team_id uuid not null references teams (id) on delete cascade
-- );

INSERT INTO invites (
    id
) VALUES (
    '95285f8f-4880-4258-8712-a622f413bd30'
),
(
    'ed75f0e9-ca28-4741-af93-a459ddbabe08'
),
(
    'd89342dc-9224-441d-a4af-bdd837a3b239'
);
