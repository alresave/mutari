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

    block INT REFERENCES mcp_blocks(id),

    question_text TEXT NOT NULL,

    scale TEXT DEFAULT 'likert_1_5'

);

CREATE TABLE user_mcp_answers (

    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT REFERENCES users(id),

    question_id INT REFERENCES mcp_questions(id),

    block INT,

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

CREATE INDEX idx_user_email ON users(email);

CREATE INDEX idx_answers_user
ON user_mcp_answers(user_id);

CREATE INDEX idx_answers_question
ON user_mcp_answers(question_id);


users
mcp_blocks
mcp_questions
user_mcp_answers
user_mcp_profiles




INSERT INTO users (name,email,password_hash,photo_url) VALUES
('Lucas Silva','lucas@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/11.jpg'),
('Mateo Torres','mateo@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/12.jpg'),
('Santiago Ruiz','santiago@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/13.jpg'),
('Sebastian Vega','sebastian@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/14.jpg'),
('Daniel Herrera','daniel@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/15.jpg'),
('Alejandro Cruz','alejandro@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/16.jpg'),
('Gabriel Medina','gabriel@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/17.jpg'),
('Nicolas Rojas','nicolas@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/18.jpg'),
('David Navarro','david@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/19.jpg'),
('Adrian Campos','adrian@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/men/20.jpg'),

('Emma Collins','emma@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/11.jpg'),
('Sophia Turner','sophia@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/12.jpg'),
('Olivia Carter','olivia@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/13.jpg'),
('Isabella Hughes','isabella@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/14.jpg'),
('Mia Bennett','mia@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/15.jpg'),
('Charlotte Reed','charlotte@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/16.jpg'),
('Amelia Foster','amelia@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/17.jpg'),
('Harper Brooks','harper@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/18.jpg'),
('Evelyn Hayes','evelyn@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/19.jpg'),
('Abigail Stone','abigail@mutari.dev','$2a$10$testhash','https://randomuser.me/api/portraits/women/20.jpg');

INSERT INTO user_mcp_answers
SELECT 13,id,block,5 FROM mcp_questions WHERE block IN (5,7);

INSERT INTO user_mcp_answers
SELECT 14,id,block,4 FROM mcp_questions WHERE block=1;

INSERT INTO user_mcp_answers
SELECT 13,id,block,3 FROM mcp_questions WHERE block NOT IN (5,7,1);

CREATE TABLE team_profiles (
    team_id INT PRIMARY KEY,
    effective_vector VECTOR(64)
);

CREATE TABLE team_members (
    team_id INT REFERENCES teams(id),
    user_id INT REFERENCES users(id)
);

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);
INSERT INTO teams (name,description)
VALUES ('Backend Team','Core platform engineers');

INSERT INTO team_members (team_id,user_id) VALUES
(1,1),
(1,2),
(1,3),
(1,4);
