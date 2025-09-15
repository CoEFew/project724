-- Feedbacks table
CREATE TABLE IF NOT EXISTS public.feedbacks (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL DEFAULT '' CHECK (length(name) <= 64),
  contact    TEXT NOT NULL DEFAULT '' CHECK (length(contact) <= 128),
  message    TEXT NOT NULL CHECK (length(message) BETWEEN 1 AND 4000),
  source     TEXT NOT NULL DEFAULT '' CHECK (length(source) <= 64),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index สำหรับดึงรายการล่าสุด
CREATE INDEX IF NOT EXISTS idx_feedbacks_created_desc
  ON public.feedbacks (created_at DESC);
