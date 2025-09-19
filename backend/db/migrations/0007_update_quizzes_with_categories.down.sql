-- Revert category updates
UPDATE public.quizzes 
SET category = 'สัตว์' 
WHERE category = 'สัตว์';
