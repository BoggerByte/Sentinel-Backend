CREATE TABLE IF NOT EXISTS admin
(
    "id"         bigserial PRIMARY KEY,
    "discord_id" bigint UNIQUE,
    "nickname"   varchar(255),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS guild
(
    "id"         bigserial PRIMARY KEY,
    "discord_id" bigint UNIQUE,
    "owner_id"   bigint
);

CREATE TABLE IF NOT EXISTS config
(
    "id"         bigserial PRIMARY KEY,
    "version"    numeric     NOT NULL default 0,
    "filename"   varchar(255) UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "guild_id"   bigserial REFERENCES guild (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS "guild_admin"
(
    "guilds_id" bigserial REFERENCES guild (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    "admins_id" bigserial REFERENCES admin (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    CONSTRAINT guild_admin_pk PRIMARY KEY ("guilds_id", "admins_id")
);

CREATE INDEX ON admin ("discord_id");

CREATE INDEX ON guild ("owner_id");

CREATE INDEX ON guild ("discord_id");

CREATE INDEX ON config ("guild_id");

COMMENT ON COLUMN admin."nickname" IS 'discord name with tag';
