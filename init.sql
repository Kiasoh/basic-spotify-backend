SELECT 'Initializing database...' as message;

-- Create database if it doesn't exist
SELECT 'Creating ds_db database...' as message;
CREATE DATABASE ds_db;

-- Connect to the database and set up extensions
\c ds_db

-- Create additional users if needed
-- CREATE USER another_user WITH PASSWORD 'secure_password';

-- Grant privileges
-- GRANT ALL PRIVILEGES ON DATABASE ds_db TO niflheim;

-- Create schema
CREATE SCHEMA IF NOT EXISTS public;
GRANT ALL ON SCHEMA public TO niflheim;

-- Set default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO niflheim;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO niflheim;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO niflheim;

SELECT 'Database initialization complete!' as message;
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

INSERT INTO users(users, password) VALUES('admin','\x123');
