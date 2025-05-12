-- Create surveys table
CREATE TABLE IF NOT EXISTS surveys (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create survey_questions table
CREATE TABLE IF NOT EXISTS survey_questions (
    id UUID PRIMARY KEY,
    survey_id UUID NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create survey_options table
CREATE TABLE IF NOT EXISTS survey_options (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES survey_questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_surveys_created_by ON surveys(created_by);
CREATE INDEX IF NOT EXISTS idx_survey_questions_survey_id ON survey_questions(survey_id);
CREATE INDEX IF NOT EXISTS idx_survey_options_question_id ON survey_options(question_id); 