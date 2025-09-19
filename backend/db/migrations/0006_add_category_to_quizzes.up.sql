-- Add category column to quizzes table
ALTER TABLE public.quizzes 
ADD COLUMN category VARCHAR(50) NOT NULL DEFAULT 'สัตว์';

-- Create index for category filtering
CREATE INDEX IF NOT EXISTS idx_quizzes_category_tier_active
  ON public.quizzes (category, tier, active);

-- Update existing data to have proper categories
-- This will be handled in the seed data migration
