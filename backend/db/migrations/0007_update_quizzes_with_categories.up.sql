-- Update existing quizzes with categories
-- All existing data will be categorized as 'สัตว์' (Animals)

UPDATE public.quizzes 
SET category = 'สัตว์' 
WHERE category = 'สัตว์' OR category IS NULL;
