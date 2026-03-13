CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    photo_url TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE mcp_blocks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE mcp_questions (
    id INT PRIMARY KEY,
    block INT NOT NULL REFERENCES mcp_blocks(id),
    question_text TEXT NOT NULL,
    scale TEXT DEFAULT 'likert_1_5'
);

CREATE TABLE user_mcp_answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    question_id INT REFERENCES mcp_questions(id),
    block INT REFERENCES mcp_blocks(id),
    value FLOAT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_mcp_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    effective_vector VECTOR(64),
    confidence_score FLOAT DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_vectors
ON user_mcp_profiles
USING ivfflat (effective_vector vector_cosine_ops)
WITH (lists = 100);

CREATE INDEX idx_user_email
ON users(email);

CREATE INDEX idx_answers_user
ON user_mcp_answers(user_id);

CREATE INDEX idx_answers_question
ON user_mcp_answers(question_id);

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE team_profiles (
    team_id INT PRIMARY KEY REFERENCES teams(id),
    effective_vector VECTOR(64)
);

CREATE TABLE team_members (
    team_id INT REFERENCES teams(id),
    user_id BIGINT REFERENCES users(id),
    PRIMARY KEY (team_id, user_id)
);