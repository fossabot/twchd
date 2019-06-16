/* Create database and user */
CREATE USER twchd WITH ENCRYPTED PASSWORD 'P@ssw0rd';
CREATE DATABASE twchd WITH OWNER twchd;

/* Substitute database and user */
\c twchd
SET ROLE twchd;

CREATE TABLE IF NOT EXISTS twchd.public.user (
  display_name TEXT UNIQUE NOT NULL,
  user_id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS twchd.public.message (
  msg_text TEXT NOT NULL,
  room_id BIGINT REFERENCES public.user(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  msg_id UUID PRIMARY KEY,
  ts TIMESTAMP NOT NULL,
  user_id BIGINT REFERENCES public.user ON DELETE CASCADE ON UPDATE CASCADE,
  user_role BIT(3) DEFAULT B'000'
);
