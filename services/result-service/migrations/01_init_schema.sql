-- Create responses_sessions table to store response sessions
CREATE TABLE IF NOT EXISTS response_sessions (
    id UUID PRIMARY KEY,
    survey_id UUID NOT NULL,
    respondent_id UUID,  -- Can be NULL for anonymous responses
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create responses table to store individual question responses
CREATE TABLE IF NOT EXISTS responses (
    id UUID PRIMARY KEY,
    response_id UUID NOT NULL REFERENCES response_sessions(id) ON DELETE CASCADE,
    survey_id UUID NOT NULL,
    question_id UUID NOT NULL,
    option_id UUID,  -- NULL for text answers
    text_answer TEXT,  -- NULL for choice answers
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create response_analytics table to store pre-computed analytics data
CREATE TABLE IF NOT EXISTS response_analytics (
    id UUID PRIMARY KEY,
    survey_id UUID NOT NULL,
    analytics_data JSONB NOT NULL,  -- Store complex analytics data as JSON
    period VARCHAR(50) NOT NULL,     -- day, week, month, year, all
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create exported_results table to store exported results
CREATE TABLE IF NOT EXISTS exported_results (
    id UUID PRIMARY KEY,
    survey_id UUID NOT NULL,
    file_path TEXT NOT NULL,
    format VARCHAR(10) NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_response_sessions_survey_id ON response_sessions(survey_id);
CREATE INDEX IF NOT EXISTS idx_response_sessions_respondent_id ON response_sessions(respondent_id);
CREATE INDEX IF NOT EXISTS idx_responses_response_id ON responses(response_id);
CREATE INDEX IF NOT EXISTS idx_responses_survey_id ON responses(survey_id);
CREATE INDEX IF NOT EXISTS idx_responses_question_id ON responses(question_id);
CREATE INDEX IF NOT EXISTS idx_responses_option_id ON responses(option_id);
CREATE INDEX IF NOT EXISTS idx_response_analytics_survey_id ON response_analytics(survey_id);
CREATE INDEX IF NOT EXISTS idx_exported_results_survey_id ON exported_results(survey_id); 