CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "username" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "avg_interest" jsonb NOT NULL DEFAULT '[]',
    "recomm_plylist_id" INTEGER NOT NULL,
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
    "description" varchar(255),
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


CREATE TABLE "spotify_tracks" (
    "track_id" TEXT,
    "artists" TEXT,
    "album_name" TEXT,
    "track_name" TEXT,
    "popularity" BIGINT,
    "duration_ms" BIGINT,
    "explicit" BOOLEAN,
    "danceability" DOUBLE PRECISION,
    "energy" DOUBLE PRECISION,
    "key" BIGINT,
    "loudness" DOUBLE PRECISION,
    "mode" BIGINT,
    "speechiness" DOUBLE PRECISION,
    "acousticness" DOUBLE PRECISION,
    "instrumentalness" DOUBLE PRECISION,
    "liveness" DOUBLE PRECISION,
    "valence" DOUBLE PRECISION,
    "tempo" DOUBLE PRECISION,
    "time_signature" BIGINT,
    "track_genre" TEXT
);

CREATE INDEX idx_track_name_trgm ON spotify_tracks USING gin (track_name gin_trgm_ops);
CREATE INDEX idx_artists_trgm ON spotify_tracks USING gin (artists gin_trgm_ops);
