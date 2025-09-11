CREATE TABLE IF NOT EXISTS public.scores (
  id          BIGSERIAL PRIMARY KEY,
  name        TEXT NOT NULL CHECK (length(name) BETWEEN 1 AND 64),
  score       INTEGER NOT NULL CHECK (score >= 0),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_scores_score_desc ON public.scores (score DESC);
