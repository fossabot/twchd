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
  ts TIMESTAMPTZ NOT NULL,
  user_id BIGINT REFERENCES public.user ON DELETE CASCADE ON UPDATE CASCADE,
  user_role BIT(3) DEFAULT B'000'
);

SET timezone = 'Asia/Yekaterinburg';

CREATE OR REPLACE PROCEDURE add_data(msg_text TEXT, msg_id UUID, ts TIMESTAMPTZ, chan_name TEXT, room_id BIGINT, display_name TEXT, user_id BIGINT, user_role BIT(3))
LANGUAGE SQL
AS $$
  INSERT INTO twchd.public.user VALUES (display_name,user_id), (chan_name,room_id) ON CONFLICT DO NOTHING;
  INSERT INTO twchd.public.message VALUES (msg_text, room_id, msg_id, ts, user_id, user_role);
$$;
