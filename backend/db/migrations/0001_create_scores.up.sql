-- Create table
CREATE TABLE IF NOT EXISTS public.scores (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(name) BETWEEN 1 AND 64),
  score      INTEGER NOT NULL CHECK (score >= 0),
  gamename   TEXT NOT NULL CHECK (length(gamename) BETWEEN 1 AND 64),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes (แนะนำ)
-- เรียกดู Top score ของเกมใดเกมหนึ่งบ่อย ๆ
CREATE INDEX IF NOT EXISTS idx_scores_gamename_score_desc
  ON public.scores (gamename, score DESC, created_at ASC);

-- (ทางเลือก) เรียกดู Top score รวมทุกเกม
CREATE INDEX IF NOT EXISTS idx_scores_score_desc
  ON public.scores (score DESC, created_at ASC);