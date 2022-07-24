CREATE TABLE IF NOT EXISTS "user"
(
    "id"            bigserial PRIMARY KEY,
    "discord_id"    varchar(64) UNIQUE NOT NULL,
    "username"      varchar(64)        NOT NULL,
    "discriminator" varchar(4)         NOT NULL,
    "verified"      bool               NOT NULL,
    "email"         varchar(255)       NOT NULL,
    "avatar"        varchar(255)       NOT NULL DEFAULT '#',
    "banner"        varchar(255)       NOT NULL DEFAULT '#',
    "accent_color"  bigint             NOT NULL DEFAULT 0,
    "created_at"    timestamptz        NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS guild
(
    "id"               bigserial PRIMARY KEY,
    "discord_id"       varchar(64) UNIQUE NOT NULL,
    "owner_discord_id" varchar(64)        NOT NULL,
    "name"             varchar(255)       NOT NULL,
    "icon"             varchar(255)       NOT NULL DEFAULT '#',
    CONSTRAINT "guild_config_fk" FOREIGN KEY ("id") REFERENCES "guild" ("id") DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE IF NOT EXISTS "guild_config"
(
    "id"         bigserial PRIMARY KEY,
    "json"       jsonb       NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT "guild_fk" FOREIGN KEY ("id") REFERENCES "guild_config" ("id") DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE IF NOT EXISTS "user_guild"
(
    "account_discord_id" varchar(64) NOT NULL,
    "guild_discord_id"   varchar(64) NOT NULL,
    "permissions"        bigint      NOT NULL,
    CONSTRAINT "account_guild_pk" PRIMARY KEY ("guild_discord_id", "account_discord_id")
);

CREATE INDEX ON "user" ("discord_id");

CREATE INDEX ON "guild" ("discord_id");

COMMENT ON COLUMN "user"."accent_color" IS 'color encoded as an integer representation of hexadecimal color code';