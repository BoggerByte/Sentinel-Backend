CREATE TABLE IF NOT EXISTS admin
(
    "id"            bigserial PRIMARY KEY,
    "discord_id"    varchar(64) UNIQUE NOT NULL,
    "username"      varchar(64)        NOT NULL,
    "discriminator" varchar(4)         NOT NULL,
    "verified"      bool               NOT NULL,
    "email"         varchar(255)       NOT NULL,
    "avatar"        varchar(255)       NOT NULL DEFAULT '#',
    "flags"         bigint             NOT NULL DEFAULT 0,
    "banner"        varchar(255)       NOT NULL DEFAULT '#',
    "accent_color"  bigint             NOT NULL DEFAULT 0,
    "public_flags"  bigint             NOT NULL DEFAULT 0,
    "created_at"    timestamptz        NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS guild
(
    "id"         bigserial PRIMARY KEY,
    "discord_id" varchar(64) UNIQUE NOT NULL,
    "name"       varchar(255)       NOT NULL,
    "icon"       varchar(255)       NOT NULL DEFAULT '#',
    "owner_id"   varchar(64)        NOT NULL
);

CREATE TABLE IF NOT EXISTS config
(
    "id"         bigserial PRIMARY KEY,
    "version"    numeric UNIQUE                                                      NOT NULL DEFAULT 0,
    "filename"   varchar(255) UNIQUE                                                 NOT NULL,
    "created_at" timestamptz                                                         NOT NULL DEFAULT now(),
    "guild_id"   bigserial REFERENCES guild (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS guild_admin
(
    "guild_id" bigserial REFERENCES guild (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    "admin_id" bigserial REFERENCES admin (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    CONSTRAINT guild_admin_pk PRIMARY KEY (guild_id, admin_id)
);

CREATE INDEX ON admin (discord_id);

CREATE INDEX ON guild (owner_id);

CREATE INDEX ON guild (discord_id);

CREATE INDEX ON config (guild_id);