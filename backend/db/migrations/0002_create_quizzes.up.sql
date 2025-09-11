CREATE TABLE IF NOT EXISTS public.quizzes (
  id         BIGSERIAL PRIMARY KEY,
  answer     TEXT NOT NULL,
  hint1      TEXT NOT NULL,
  hint2      TEXT NOT NULL,
  tier       INT  NOT NULL DEFAULT 1,     -- 1 = easy, 2 = hard
  active     BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- index ใช้ค้นหา
CREATE INDEX IF NOT EXISTS idx_quizzes_tier_active
  ON public.quizzes (tier, active);

-- ป้องกันข้อมูลซ้ำ (คำตอบเดียวกันใน tier เดียวกัน)
CREATE UNIQUE INDEX IF NOT EXISTS uq_quizzes_answer_tier
  ON public.quizzes (answer, tier);
