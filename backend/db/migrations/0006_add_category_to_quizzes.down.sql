-- Remove category column from quizzes table
DROP INDEX IF EXISTS idx_quizzes_category_tier_active;
ALTER TABLE public.quizzes DROP COLUMN IF EXISTS category;
