CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "username" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "songs" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "artist" varchar(255) NOT NULL,
    "album" varchar(255) NOT NULL,
    "genre" varchar(255) NOT NULL,
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "playlists" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL DEFAULT 'playlist #',
    "owner_id" INTEGER NOT NULL REFERENCES "users"("id"),
    "modifyable" boolean NOT NULL DEFAULT true,
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "songs_playlists" (
    "playlist_id" INTEGER NOT NULL REFERENCES "playlists"("id"),
    "song_id" INTEGER NOT NULL REFERENCES "songs"("id"),
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "interactions" (
    "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
    "song_id" INTEGER NOT NULL REFERENCES "songs"("id"),
    "type" varchar(255) NOT NULL,
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

INSERT INTO users(users, password) VALUES('admin', );
