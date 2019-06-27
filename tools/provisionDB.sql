/* Create database and user */
CREATE USER twchd WITH ENCRYPTED PASSWORD 'P@ssw0rd';
CREATE DATABASE twchd WITH OWNER twchd;

/* Substitute database and user */
\c twchd
SET ROLE twchd;

CREATE TABLE IF NOT EXISTS twchd.public.user (
  user_id BIGINT PRIMARY KEY,
  display_name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS twchd.public.message (
  msg_id UUID PRIMARY KEY,
  ts TIMESTAMPTZ NOT NULL,
  room_id BIGINT REFERENCES public.user(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  user_id BIGINT REFERENCES public.user ON DELETE CASCADE ON UPDATE CASCADE,
  user_role BIT(3) DEFAULT B'000',
  msg_text TEXT NOT NULL
);

CREATE OR REPLACE PROCEDURE add_data(msg_id UUID, ts TIMESTAMPTZ, room_id BIGINT,
  chan_name TEXT, user_id BIGINT, display_name TEXT, user_role BIT(3), msg_text TEXT)
LANGUAGE SQL
AS $$
  INSERT INTO twchd.public.user VALUES (user_id, display_name), (room_id, chan_name) ON CONFLICT DO NOTHING;
  INSERT INTO twchd.public.message VALUES (msg_id, ts, room_id, user_id, user_role, msg_text);
$$;
