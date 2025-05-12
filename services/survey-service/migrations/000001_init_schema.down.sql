-- Drop indexes
DROP INDEX IF EXISTS idx_survey_options_question_id;
DROP INDEX IF EXISTS idx_survey_questions_survey_id;
DROP INDEX IF EXISTS idx_surveys_created_by;

-- Drop tables
DROP TABLE IF EXISTS survey_options;
DROP TABLE IF EXISTS survey_questions;
DROP TABLE IF EXISTS surveys; 