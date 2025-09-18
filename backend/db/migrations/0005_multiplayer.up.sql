-- ห้องเล่นหลายคน (สูงสุด 4 คน)
CREATE TABLE IF NOT EXISTS public.rooms (
  id           BIGSERIAL PRIMARY KEY,
  code         TEXT NOT NULL UNIQUE CHECK (length(code) BETWEEN 4 AND 12),
  owner_name   TEXT NOT NULL CHECK (length(owner_name) BETWEEN 1 AND 64),
  status       TEXT NOT NULL CHECK (status IN ('waiting', 'playing', 'finished')),
  max_players  INTEGER NOT NULL DEFAULT 4 CHECK (max_players BETWEEN 2 AND 4),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.room_players (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT NOT NULL REFERENCES public.rooms(id) ON DELETE CASCADE,
  name         TEXT NOT NULL CHECK (length(name) BETWEEN 1 AND 64),
  is_owner     BOOLEAN NOT NULL DEFAULT FALSE,
  is_ready     BOOLEAN NOT NULL DEFAULT FALSE,
  score        INTEGER NOT NULL DEFAULT 0 CHECK (score >= 0),
  is_out       BOOLEAN NOT NULL DEFAULT FALSE,
  joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (room_id, name)
);

-- รอบเกม แบ่งด้วยรอบและข้อ (ทุกคนได้ "คำเดียวกัน")
CREATE TABLE IF NOT EXISTS public.room_rounds (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT NOT NULL REFERENCES public.rooms(id) ON DELETE CASCADE,
  round_no     INTEGER NOT NULL CHECK (round_no >= 1),
  quiz_id      TEXT NOT NULL,
  quiz_token   TEXT NOT NULL,
  quiz_exp     BIGINT NOT NULL,
  started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  ended_at     TIMESTAMPTZ,
  UNIQUE (room_id, round_no)
);

-- การเดาในรอบนั้น ๆ (ใช้เก็บผลลัพธ์/ลำดับเวลา)
CREATE TABLE IF NOT EXISTS public.room_guesses (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT NOT NULL REFERENCES public.rooms(id) ON DELETE CASCADE,
  round_no     INTEGER NOT NULL CHECK (round_no >= 1),
  player_id    BIGINT NOT NULL REFERENCES public.room_players(id) ON DELETE CASCADE,
  guess        TEXT NOT NULL,
  correct      BOOLEAN NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for queries
CREATE INDEX IF NOT EXISTS idx_room_players_room ON public.room_players(room_id);
CREATE INDEX IF NOT EXISTS idx_room_rounds_room ON public.room_rounds(room_id);
CREATE INDEX IF NOT EXISTS idx_room_guesses_room_round ON public.room_guesses(room_id, round_no);
