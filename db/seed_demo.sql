-- ============================================
-- RESET
-- ============================================

TRUNCATE team_members, team_profiles, teams, user_mcp_profiles, user_mcp_answers, users CASCADE;

INSERT INTO mcp_blocks (id,name,description) VALUES

(0,'Core Values','Fundamental life principles'),
(1,'Lifestyle Rhythm','Daily habits and life pace'),
(2,'Emotional Stability','Emotional regulation and stress'),
(3,'Relational Architecture','How relationships are structured'),
(4,'Ambition','Drive and long term goals'),
(5,'Growth','Personal development orientation'),
(6,'Trust & Collaboration ','Emotional bonding and closeness'),
(7,'Cognitive Style','Thinking and decision making'),
(8,'Control','Consistency validation questions');


INSERT INTO mcp_questions (id,block,question_text,scale) VALUES

-- Core Values
(1,0,'I prefer stability over constant change','likert_1_5'),
(2,0,'Long-term commitment is important to me','likert_1_5'),
(3,0,'I believe loyalty is essential in relationships','likert_1_5'),
(4,0,'I value honesty above convenience','likert_1_5'),
(5,0,'Shared values matter more than attraction','likert_1_5'),
(6,0,'Keeping promises is very important to me','likert_1_5'),
(7,0,'I believe trust should be built slowly','likert_1_5'),
(8,0,'I admire people with strong principles','likert_1_5'),

-- Lifestyle Rhythm
(9,1,'I enjoy having a structured daily routine','likert_1_5'),
(10,1,'I prefer calm evenings over nightlife','likert_1_5'),
(11,1,'I enjoy planning my week ahead','likert_1_5'),
(12,1,'I feel comfortable with predictable schedules','likert_1_5'),
(13,1,'I prefer quiet environments over busy places','likert_1_5'),
(14,1,'I enjoy peaceful weekends more than parties','likert_1_5'),
(15,1,'I like consistency in daily habits','likert_1_5'),
(16,1,'I feel more productive with routines','likert_1_5'),

-- Emotional Stability
(17,2,'I tend to stay calm under stress','likert_1_5'),
(18,2,'I handle disagreements calmly','likert_1_5'),
(19,2,'I rarely react impulsively','likert_1_5'),
(20,2,'I try to understand emotions before reacting','likert_1_5'),
(21,2,'I recover quickly after difficult situations','likert_1_5'),
(22,2,'I avoid escalating conflicts','likert_1_5'),
(23,2,'I remain balanced during emotional discussions','likert_1_5'),
(24,2,'Emotional maturity is important to me','likert_1_5'),

-- Relational Architecture
(25,3,'I believe relationships should function as a team','likert_1_5'),
(26,3,'Important decisions should be made together','likert_1_5'),
(27,3,'Mutual support is essential in relationships','likert_1_5'),
(28,3,'Communication solves most relationship issues','likert_1_5'),
(29,3,'Shared goals strengthen relationships','likert_1_5'),
(30,3,'Both partners should participate in problem solving','likert_1_5'),
(31,3,'Healthy relationships require effort','likert_1_5'),
(32,3,'Open dialogue is key to strong partnerships','likert_1_5'),

-- Ambition
(33,4,'I have clear long-term goals','likert_1_5'),
(34,4,'I constantly work to improve my situation','likert_1_5'),
(35,4,'Professional growth is important to me','likert_1_5'),
(36,4,'I admire ambitious people','likert_1_5'),
(37,4,'I enjoy pursuing difficult goals','likert_1_5'),
(38,4,'I strive to build a meaningful life','likert_1_5'),
(39,4,'I feel motivated by long-term progress','likert_1_5'),
(40,4,'I respect people with strong determination','likert_1_5'),

-- Growth
(41,5,'I actively try to learn from mistakes','likert_1_5'),
(42,5,'Personal development is important to me','likert_1_5'),
(43,5,'I enjoy learning new perspectives','likert_1_5'),
(44,5,'I believe people should evolve over time','likert_1_5'),
(45,5,'Self improvement is part of my life','likert_1_5'),
(46,5,'I enjoy gaining new knowledge','likert_1_5'),
(47,5,'I reflect on my behavior regularly','likert_1_5'),
(48,5,'I welcome opportunities for personal growth','likert_1_5'),

-- Trust & Collaboration 
(49,6,'I feel comfortable expressing affection','likert_1_5'),
(50,6,'Emotional closeness is important to me','likert_1_5'),
(51,6,'I enjoy deep conversations','likert_1_5'),
(52,6,'Trust grows through emotional openness','likert_1_5'),
(53,6,'I appreciate meaningful time together','likert_1_5'),
(54,6,'I value emotional transparency','likert_1_5'),
(55,6,'I prefer deep bonds over casual connections','likert_1_5'),
(56,6,'Strong emotional connection matters to me','likert_1_5'),

-- Cognitive Style
(57,7,'I prefer analyzing problems logically','likert_1_5'),
(58,7,'I enjoy discussing complex ideas','likert_1_5'),
(59,7,'I like understanding how things work','likert_1_5'),
(60,7,'I think carefully before making decisions','likert_1_5'),
(61,7,'I enjoy intellectual discussions','likert_1_5'),
(62,7,'I try to understand multiple perspectives','likert_1_5'),
(63,7,'I reflect before reacting','likert_1_5'),
(64,7,'I enjoy solving complex problems','likert_1_5'),

-- Control Questions
(65,8,'I rarely reflect on my own mistakes','likert_1_5'),
(66,8,'I almost never feel stressed','likert_1_5'),
(67,8,'I make important decisions impulsively','likert_1_5'),
(68,8,'I avoid discussing difficult topics','likert_1_5'),
(69,8,'I change my long-term plans constantly','likert_1_5'),
(70,8,'I dislike planning ahead','likert_1_5'),
(71,8,'I rarely question my own opinions','likert_1_5'),
(72,8,'I never doubt my decisions','likert_1_5');


-- ============================================
-- USERS (30 demo candidates)
-- ============================================


INSERT INTO users (name,email,password_hash,photo_url) VALUES
('Alberto Gonzalez','alresave@gmail.com','123456','https://randomuser.me/api/portraits/men/11.jpg'),
('Alicia Gonzalez','alicia.gonzalez@gmail.com','12345678','https://randomuser.me/api/portraits/women/11.jpg'),
('Daniel Schmidt','daniel.schmidt@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/13.jpg'),
('Marco Rossi','marco.rossi@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/14.jpg'),
('Carlos Alvarez','carlos.alvarez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/15.jpg'),
('Juan Torres','juan.torres@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/16.jpg'),
('Mateo Gonzalez','mateo.gonzalez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/17.jpg'),
('Pedro Carvalho','pedro.carvalho@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/18.jpg'),
('Nicolas Dubois','nicolas.dubois@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/19.jpg'),
('Liam Murphy','liam.murphy@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/20.jpg'),
('Oliver Johnson','oliver.johnson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/21.jpg'),
('Noah Brown','noah.brown@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/22.jpg'),
('Ethan Williams','ethan.williams@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/23.jpg'),
('Alexander Ivanov','alex.ivanov@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/24.jpg'),
('Yuki Tanaka','yuki.tanaka@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/25.jpg'),
('Min-Jun Park','minjun.park@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/26.jpg'),
('Arjun Patel','arjun.patel@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/27.jpg'),
('Ahmed Hassan','ahmed.hassan@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/28.jpg'),
('Omar Khaled','omar.khaled@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/29.jpg'),
('George Thompson','george.thompson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/30.jpg'),
('William Scott','william.scott@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/31.jpg'),
('Benjamin Lee','benjamin.lee@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/32.jpg'),
('Logan White','logan.white@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/33.jpg'),
('Henry Wilson','henry.wilson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/34.jpg'),
('Jack Anderson','jack.anderson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/men/35.jpg'),

('Emma Wilson','emma.wilson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/11.jpg'),
('Olivia Taylor','olivia.taylor@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/12.jpg'),
('Sophia Brown','sophia.brown@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/13.jpg'),
('Isabella Garcia','isabella.garcia@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/14.jpg'),
('Mia Martinez','mia.martinez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/15.jpg'),
('Charlotte Dubois','charlotte.dubois@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/16.jpg'),
('Amelia Rossi','amelia.rossi@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/17.jpg'),
('Harper Johnson','harper.johnson@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/18.jpg'),
('Evelyn Davis','evelyn.davis@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/19.jpg'),
('Abigail Miller','abigail.miller@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/20.jpg'),
('Emily Clark','emily.clark@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/21.jpg'),
('Sofia Hernandez','sofia.hernandez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/22.jpg'),
('Camila Torres','camila.torres@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/23.jpg'),
('Valentina Lopez','valentina.lopez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/24.jpg'),
('Lucia Fernandez','lucia.fernandez@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/25.jpg'),
('Anna Schmidt','anna.schmidt@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/26.jpg'),
('Laura Müller','laura.mueller@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/27.jpg'),
('Giulia Bianchi','giulia.bianchi@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/28.jpg'),
('Yuna Kim','yuna.kim@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/29.jpg'),
('Aiko Sato','aiko.sato@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/30.jpg'),
('Priya Sharma','priya.sharma@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/31.jpg'),
('Fatima Ali','fatima.ali@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/32.jpg'),
('Layla Hassan','layla.hassan@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/33.jpg'),
('Nora Jensen','nora.jensen@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/34.jpg'),
('Chloe Martin','chloe.martin@gmail.com','$2y$10$hashdemo','https://randomuser.me/api/portraits/women/35.jpg');
SELECT setval(
    pg_get_serial_sequence('users','id'),
    COALESCE((SELECT MAX(id) FROM users),1),
    true
);

-- ============================================
-- TEAMS
-- ============================================

INSERT INTO teams (name) VALUES
('Startup Founders'),
('Product Team'),
('Engineering Team');

-- ============================================
-- TEAM MEMBERS
-- ============================================

INSERT INTO team_members (team_id,user_id) VALUES
(1,1),
(1,2),
(1,5),

(2,2),
(2,4),
(2,6),

(3,3),
(3,10),
(3,20);

-- ============================================
-- USER VECTOR GENERATION (64 dimensions)
-- ============================================

-- ============================================
-- USER VECTOR GENERATION (REALISTIC CLUSTERS)
-- ============================================

-- ============================================
-- USER VECTOR GENERATION (REALISTIC PROFILES)
-- ============================================

INSERT INTO user_mcp_profiles (user_id, effective_vector)
SELECT
    id,
    (
        SELECT array_agg(
            GREATEST(
                0,
                LEAST(
                    1,
                    base + gaussian_noise
                )
            )
        )
        FROM (
            SELECT
                g,
                archetype,
                CASE archetype

                    -- Leader
                    WHEN 0 THEN
                        CASE (g-1)/8
                            WHEN 4 THEN 0.85
                            WHEN 5 THEN 0.75
                            ELSE 0.40
                        END

                    -- Strategist
                    WHEN 1 THEN
                        CASE (g-1)/8
                            WHEN 7 THEN 0.90
                            WHEN 5 THEN 0.75
                            ELSE 0.40
                        END

                    -- Builder
                    WHEN 2 THEN
                        CASE (g-1)/8
                            WHEN 5 THEN 0.85
                            WHEN 4 THEN 0.80
                            ELSE 0.40
                        END

                    -- Connector
                    WHEN 3 THEN
                        CASE (g-1)/8
                            WHEN 6 THEN 0.90
                            WHEN 3 THEN 0.85
                            ELSE 0.45
                        END

                    -- Stabilizer
                    ELSE
                        CASE (g-1)/8
                            WHEN 2 THEN 0.90
                            WHEN 0 THEN 0.85
                            ELSE 0.45
                        END
                END AS base,

                -- gaussian noise
                (random() + random() + random() + random() - 2) * 0.15 AS gaussian_noise

            FROM generate_series(1,64) g
        ) s
    )::vector
FROM (
    SELECT id, (id % 5) AS archetype
    FROM users
) u;
-- ============================================
-- TEAM VECTOR GENERATION (64 dimensions)
-- ============================================

-- ============================================
-- TEAM VECTOR GENERATION (CLEAR ARCHETYPES)
-- ============================================

-- ============================================
-- TEAM VECTOR GENERATION (REALISTIC TEAMS)
-- ============================================

INSERT INTO team_profiles (team_id, effective_vector)
SELECT
    id,
    (
        SELECT array_agg(
            GREATEST(
                0,
                LEAST(
                    1,
                    base + gaussian_noise
                )
            )
        )
        FROM (
            SELECT
                g,
                archetype,

                CASE archetype

                    -- Visionary Startup
                    WHEN 0 THEN
                        CASE (g-1)/8
                            WHEN 7 THEN 0.90
                            WHEN 4 THEN 0.85
                            ELSE 0.40
                        END

                    -- Execution Engine
                    WHEN 1 THEN
                        CASE (g-1)/8
                            WHEN 5 THEN 0.90
                            WHEN 4 THEN 0.85
                            ELSE 0.40
                        END

                    -- Collaborative Studio
                    WHEN 2 THEN
                        CASE (g-1)/8
                            WHEN 6 THEN 0.90
                            WHEN 3 THEN 0.85
                            ELSE 0.45
                        END

                    -- Stable Organization
                    WHEN 3 THEN
                        CASE (g-1)/8
                            WHEN 2 THEN 0.90
                            WHEN 0 THEN 0.85
                            ELSE 0.45
                        END

                    -- Balanced Team
                    ELSE
                        0.55
                END AS base,

                -- Gaussian noise
                (random() + random() + random() + random() - 2) * 0.12 AS gaussian_noise

            FROM generate_series(1,64) g
        ) s
    )::vector
FROM (
    SELECT id, (id % 5) AS archetype
    FROM teams
) t;