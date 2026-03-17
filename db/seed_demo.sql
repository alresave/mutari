-- ============================================
-- COMPANY
-- ============================================

INSERT INTO companies (name, website)
VALUES ('Mutari LLC','https://mutari.dev')
ON CONFLICT DO NOTHING;


-- ============================================
-- MCP BLOCKS
-- ============================================

INSERT INTO mcp_blocks (id,name,description) VALUES
(0,'Core Values','Fundamental life principles'),
(1,'Lifestyle Rhythm','Daily habits and life pace'),
(2,'Emotional Stability','Emotional regulation and stress'),
(3,'Relational Architecture','How relationships are structured'),
(4,'Ambition','Drive and long term goals'),
(5,'Growth','Personal development orientation'),
(6,'Trust & Collaboration','Emotional bonding and closeness'),
(7,'Cognitive Style','Thinking and decision making'),
(8,'Control','Consistency validation questions')
ON CONFLICT DO NOTHING;


-- ============================================
-- MCP QUESTIONS
-- ============================================

INSERT INTO mcp_questions (id,block,question_text,scale) VALUES
(1,0,'I prefer stability over constant change','likert_1_5'),
(2,0,'Long-term commitment is important to me','likert_1_5'),
(3,0,'I believe loyalty is essential in relationships','likert_1_5'),
(4,0,'I value honesty above convenience','likert_1_5'),
(5,0,'Shared values matter more than attraction','likert_1_5'),
(6,0,'Keeping promises is very important to me','likert_1_5'),
(7,0,'I believe trust should be built slowly','likert_1_5'),
(8,0,'I admire people with strong principles','likert_1_5'),

(9,1,'I enjoy having a structured daily routine','likert_1_5'),
(10,1,'I prefer calm evenings over nightlife','likert_1_5'),
(11,1,'I enjoy planning my week ahead','likert_1_5'),
(12,1,'I feel comfortable with predictable schedules','likert_1_5'),
(13,1,'I prefer quiet environments over busy places','likert_1_5'),
(14,1,'I enjoy peaceful weekends more than parties','likert_1_5'),
(15,1,'I like consistency in daily habits','likert_1_5'),
(16,1,'I feel more productive with routines','likert_1_5'),

(17,2,'I tend to stay calm under stress','likert_1_5'),
(18,2,'I handle disagreements calmly','likert_1_5'),
(19,2,'I rarely react impulsively','likert_1_5'),
(20,2,'I try to understand emotions before reacting','likert_1_5'),
(21,2,'I recover quickly after difficult situations','likert_1_5'),
(22,2,'I avoid escalating conflicts','likert_1_5'),
(23,2,'I remain balanced during emotional discussions','likert_1_5'),
(24,2,'Emotional maturity is important to me','likert_1_5'),

(25,3,'I believe relationships should function as a team','likert_1_5'),
(26,3,'Important decisions should be made together','likert_1_5'),
(27,3,'Mutual support is essential in relationships','likert_1_5'),
(28,3,'Communication solves most relationship issues','likert_1_5'),
(29,3,'Shared goals strengthen relationships','likert_1_5'),
(30,3,'Both partners should participate in problem solving','likert_1_5'),
(31,3,'Healthy relationships require effort','likert_1_5'),
(32,3,'Open dialogue is key to strong partnerships','likert_1_5'),

(33,4,'I have clear long-term goals','likert_1_5'),
(34,4,'I constantly work to improve my situation','likert_1_5'),
(35,4,'Professional growth is important to me','likert_1_5'),
(36,4,'I admire ambitious people','likert_1_5'),
(37,4,'I enjoy pursuing difficult goals','likert_1_5'),
(38,4,'I strive to build a meaningful life','likert_1_5'),
(39,4,'I feel motivated by long-term progress','likert_1_5'),
(40,4,'I respect people with strong determination','likert_1_5'),

(41,5,'I actively try to learn from mistakes','likert_1_5'),
(42,5,'Personal development is important to me','likert_1_5'),
(43,5,'I enjoy learning new perspectives','likert_1_5'),
(44,5,'I believe people should evolve over time','likert_1_5'),
(45,5,'Self improvement is part of my life','likert_1_5'),
(46,5,'I enjoy gaining new knowledge','likert_1_5'),
(47,5,'I reflect on my behavior regularly','likert_1_5'),
(48,5,'I welcome opportunities for personal growth','likert_1_5'),

(49,6,'I feel comfortable expressing affection','likert_1_5'),
(50,6,'Emotional closeness is important to me','likert_1_5'),
(51,6,'I enjoy deep conversations','likert_1_5'),
(52,6,'Trust grows through emotional openness','likert_1_5'),
(53,6,'I appreciate meaningful time together','likert_1_5'),
(54,6,'I value emotional transparency','likert_1_5'),
(55,6,'I prefer deep bonds over casual connections','likert_1_5'),
(56,6,'Strong emotional connection matters to me','likert_1_5'),

(57,7,'I prefer analyzing problems logically','likert_1_5'),
(58,7,'I enjoy discussing complex ideas','likert_1_5'),
(59,7,'I like understanding how things work','likert_1_5'),
(60,7,'I think carefully before making decisions','likert_1_5'),
(61,7,'I enjoy intellectual discussions','likert_1_5'),
(62,7,'I try to understand multiple perspectives','likert_1_5'),
(63,7,'I reflect before reacting','likert_1_5'),
(64,7,'I enjoy solving complex problems','likert_1_5')
ON CONFLICT DO NOTHING;


-- ============================================
-- USERS
-- ============================================

INSERT INTO users (name,email,password_hash,photo_url, company_id) VALUES
('Alberto Gonzalez','alberto@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/11.jpg',1),
('Mateo Torres','mateo@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/12.jpg',1),
('Santiago Ruiz','santiago@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/13.jpg',1),
('Sebastian Vega','sebastian@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/14.jpg',1),
('Daniel Herrera','daniel@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/15.jpg',1),
('Alejandro Cruz','alejandro@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/16.jpg',1),
('Gabriel Medina','gabriel@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/17.jpg',1),
('Nicolas Rojas','nicolas@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/18.jpg',1),
('David Navarro','david@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/19.jpg',1),
('Adrian Campos','adrian@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/20.jpg',1),

('Emma Collins','emma@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/11.jpg',1),
('Sophia Turner','sophia@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/12.jpg',1),
('Olivia Carter','olivia@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/13.jpg',1),
('Isabella Hughes','isabella@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/14.jpg',1),
('Mia Bennett','mia@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/15.jpg',1),
('Charlotte Reed','charlotte@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/16.jpg',1),
('Amelia Foster','amelia@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/17.jpg',1),
('Harper Brooks','harper@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/18.jpg',1),
('Evelyn Hayes','evelyn@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/19.jpg',1),
('Abigail Stone','abigail@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/20.jpg',1),

('Lucas Mendoza','lucas@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/21.jpg',1),
('Diego Castillo','diego@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/22.jpg',1),
('Juan Morales','juan@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/23.jpg',1),
('Andres Silva','andres@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/24.jpg',1),
('Carlos Ortega','carlos@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/25.jpg',1),
('Miguel Vargas','miguel@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/26.jpg',1),
('Fernando Soto','fernando@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/27.jpg',1),
('Ricardo Paredes','ricardo@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/28.jpg',1),
('Pablo Guerrero','pablo@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/29.jpg',1),
('Esteban Fuentes','esteban@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/30.jpg',1),
('Hector Aguilar','hector@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/31.jpg',1),
('Jorge Ibarra','jorge@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/32.jpg',1),
('Ivan Salazar','ivan@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/33.jpg',1),
('Oscar Valdez','oscar@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/34.jpg',1),
('Rafael Lozano','rafael@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/men/35.jpg',1),

('Emily Parker','emily@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/21.jpg',1),
('Lily Watson','lily@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/22.jpg',1),
('Chloe Bennett','chloe@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/23.jpg',1),
('Ella Richardson','ella@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/24.jpg',1),
('Grace Coleman','grace@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/25.jpg',1),
('Scarlett Simmons','scarlett@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/26.jpg',1),
('Victoria Sanders','victoria@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/27.jpg',1),
('Aria Matthews','aria@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/28.jpg',1),
('Hannah Powell','hannah@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/29.jpg',1),
('Zoe Long','zoe@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/30.jpg',1),
('Layla Patterson','layla@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/31.jpg',1),
('Nora Hughes','nora@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/32.jpg',1),
('Lucy Flores','lucy@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/33.jpg',1),
('Stella Kim','stella@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/34.jpg',1),
('Aurora Diaz','aurora@mutari.dev', crypt('123456', gen_salt('bf')), 'https://randomuser.me/api/portraits/women/35.jpg',1);


-- ============================================
-- USER MCP PROFILES
-- ============================================

INSERT INTO user_mcp_profiles (user_id,effective_vector)
SELECT
    u.id,
    (
        ARRAY(
            SELECT
                round(
                    (
                        CASE
                            WHEN random() < 0.45 THEN random()*0.25
                            WHEN random() > 0.55 THEN 0.75 + random()*0.25
                            ELSE 0.45 + random()*0.1
                        END
                        +
                        (user_bias * 0.15)
                        +
                        ((random()-0.5)*0.05)
                    )::numeric
                ,3)
            FROM generate_series(1,64)
        )
    )::vector
FROM (
    SELECT id,(random()-0.5) AS user_bias
    FROM users
) u
ON CONFLICT (user_id)
DO UPDATE SET effective_vector = EXCLUDED.effective_vector;


-- ============================================
-- ASSIGN COMPANY OWNER
-- ============================================

UPDATE companies
SET user_id = 1
WHERE id = 1;


-- ============================================
-- TEAMS
-- ============================================

INSERT INTO teams (name,description,company_id) VALUES
('Backend Team','Core platform engineers',1),
('Frontend Team','User interface and experience',1),
('Data Science Team','Data analysis and machine learning',1),
('Product Team','Product management and design',1),
('Growth Team','Marketing and growth strategies',1),
('Customer Success Team','Customer support and success',1),
('Experimental Team','Internal research team',1)
ON CONFLICT DO NOTHING;


-- ============================================
-- TEAM MEMBERS
-- ============================================

INSERT INTO team_members (team_id,user_id) VALUES
(1,1),(1,2),(1,3),(1,4),(1,5),
(2,6),(2,7),(2,8),(2,9),(2,10),
(3,11),(3,12),(3,13),(3,14),(3,15),
(4,16),(4,17),(4,18),(4,19),(4,20),
(5,21),(5,22),(5,23),(5,24),(5,25),
(5,26),(5,27),(5,28),(5,29),(5,30),
(6,31),(6,32),(6,33),(6,34),(6,35),
(7,36),(7,37),(7,38),(7,39),(7,40)
ON CONFLICT DO NOTHING;