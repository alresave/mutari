-- ============================================
-- COMPANY
-- ============================================

INSERT INTO companies (id, name, website) VALUES
(1, 'NovaTech', 'https://novatech.dev'),
(2, 'Atlas Systems', 'https://atlas-systems.dev'),
(3, 'Horizon Labs', 'https://horizon-labs.dev'),
(4, 'Pinnacle Solutions', 'https://pinnacle-solutions.dev'),
(5, 'Vertex Innovations', 'https://vertex-innovations.dev');



-- ============================================
-- MCP BLOCKS
-- ============================================

INSERT INTO mcp_blocks (name, description) VALUES
( 'Execution Drive', 'level of initiative, accountability, and focus on delivering results'),
( 'Cognitive Processing', 'preferred way of thinking, analyzing information, and making decisions'),
( 'Professional Values', 'core beliefs and principles that guide professional priorities and motivation'),
( 'Emotional Resilience', 'ability to manage stress, pressure, and emotional responses effectively'),
( 'Learning Velocity', 'capacity to learn quickly, adapt, and apply new knowledge in changing contexts'),
( 'Collaboration Depth', 'quality of engagement in teamwork, communication, and shared problem-solving'),
( 'Operational Rhythm', 'preferred pace, structure, and consistency in executing tasks and responsibilities'),
( 'Team Positioning', 'natural role and contribution within team dynamics and organizational structures');
-- ============================================
-- MCP QUESTIONS
-- ============================================

INSERT INTO mcp_questions (block, question_text, scale) VALUES

-- 1: Execution Drive
(1,'I take ownership of problems beyond my formal responsibilities.','0-1'),
(1,'I consistently push tasks to completion without external pressure.','0-1'),
(1,'I feel responsible for the overall success of my team.','0-1'),
(1,'I proactively identify and solve inefficiencies.','0-1'),
(1,'I perform well in high-performance environments.','0-1'),
(1,'I hold myself accountable for outcomes, not just effort.','0-1'),
(1,'I maintain high standards even under pressure.','0-1'),
(1,'I naturally step up when leadership is needed.','0-1'),

-- 2: Cognitive Processing
(2,'I prefer structured approaches when solving problems.','0-1'),
(2,'I rely on data and evidence when making decisions.','0-1'),
(2,'I break complex problems into smaller parts.','0-1'),
(2,'I enjoy analyzing systems and patterns.','0-1'),
(2,'I can quickly understand abstract concepts.','0-1'),
(2,'I consider multiple perspectives before deciding.','0-1'),
(2,'I prefer logical consistency over intuition.','0-1'),
(2,'I adapt my thinking based on new information.','0-1'),

-- 3: Professional Values
(3,'Integrity is more important than short-term results.','0-1'),
(3,'I value accountability in team environments.','0-1'),
(3,'I prefer clear expectations and responsibilities.','0-1'),
(3,'I prioritize long-term impact over quick wins.','0-1'),
(3,'I value transparency in communication.','0-1'),
(3,'I believe reliability is critical in teams.','0-1'),
(3,'I respect structured processes and standards.','0-1'),
(3,'I align my work with organizational goals.','0-1'),

-- 4: Emotional Resilience
(4,'I remain calm under pressure.','0-1'),
(4,'I recover quickly from setbacks.','0-1'),
(4,'I handle criticism constructively.','0-1'),
(4,'I maintain focus during stressful situations.','0-1'),
(4,'I can manage frustration without affecting others.','0-1'),
(4,'I stay composed in uncertain environments.','0-1'),
(4,'I regulate my emotions in professional settings.','0-1'),
(4,'I maintain productivity during difficult periods.','0-1'),

-- 5: Learning Velocity
(5,'I actively seek opportunities to learn new skills.','0-1'),
(5,'I adapt quickly to new tools or processes.','0-1'),
(5,'I learn effectively from mistakes.','0-1'),
(5,'I enjoy experimenting with new approaches.','0-1'),
(5,'I integrate feedback quickly into my work.','0-1'),
(5,'I stay updated with new knowledge in my field.','0-1'),
(5,'I am comfortable working outside my expertise.','0-1'),
(5,'I continuously improve my performance.','0-1'),

-- 6: Collaboration Depth
(6,'I build trust easily with colleagues.','0-1'),
(6,'I actively support the success of others.','0-1'),
(6,'I communicate openly in team settings.','0-1'),
(6,'I value deep collaboration over individual work.','0-1'),
(6,'I help resolve conflicts within teams.','0-1'),
(6,'I share knowledge proactively.','0-1'),
(6,'I enjoy mentoring others.','0-1'),
(6,'I adapt my communication style to different people.','0-1'),

-- 7: Operational Rhythm
(7,'I work effectively under tight deadlines.','0-1'),
(7,'I manage multiple priorities efficiently.','0-1'),
(7,'I prefer structured workflows.','0-1'),
(7,'I maintain consistent productivity.','0-1'),
(7,'I am comfortable in fast-paced environments.','0-1'),
(7,'I organize my work systematically.','0-1'),
(7,'I balance speed with quality.','0-1'),
(7,'I adapt to changing priorities quickly.','0-1'),

-- 8: Team Positioning
(8,'I naturally take leadership in group settings.','0-1'),
(8,'I help coordinate team efforts.','0-1'),
(8,'I see how individual work connects to team goals.','0-1'),
(8,'I support alignment within teams.','0-1'),
(8,'I clarify roles and responsibilities when needed.','0-1'),
(8,'I contribute to team cohesion.','0-1'),
(8,'I help others stay aligned with objectives.','0-1'),
(8,'I understand my role within a broader system.','0-1');
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


-- MCP Profile Vector Seed
DELETE FROM user_mcp_profiles;

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (1, '[0.989,0.02,0.218,0.143,0.123,0.033,0.2,0.08,0.199,0.211,0.024,0.249,0.047,0.963,0.143,0.176,0.063,0.22,0.246,0.057,0.158,0.019,0.006,0.05,0.242,0.125,0.047,0.169,0.041,0.994,0.138,0.231,0.117,0.174,0.186,0.032,0.163,0.116,0.187,0.061,0.028,0.211,0.058,0.199,0.069,0.125,0.011,0.989,0.113,0.027,0.228,0.138,0.043,0.057,0.199,0.052,0.129,0.07,0.061,0.148,0.215,0.223,0.249,0.009]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (2, '[0.131,0.192,0.142,0.214,0.246,0.014,0.095,0.966,0.192,0.013,0.155,0.23,0.011,0.115,0.197,0.231,0.164,0.063,0.02,0.204,0.955,0.123,0.036,0.202,0.101,0.113,0.01,0.244,0.154,0.027,0.077,0.113,0.248,0.214,0.145,0.242,0.98,0.056,0.008,0.122,0.215,0.098,0.113,0.086,0.153,0.087,0.197,0.111,0.088,0.002,0.18,0.231,0.2,0.123,0.987,0.07,0.1,0.18,0.084,0.01,0.094,0.066,0.188,0.081]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (3, '[0.21,0.169,0.17,0.007,0.193,0.182,0.15,0.143,0.073,0.128,0.242,0.182,0.042,0.181,0.998,0.216,0.19,0.027,0.067,0.08,0.097,0.155,0.219,0.033,0.219,0.117,0.191,0.959,0.193,0.066,0.033,0.189,0.106,0.189,0.04,0.228,0.099,0.201,0.108,0.155,0.018,0.037,0.153,0.969,0.226,0.142,0.066,0.052,0.215,0.176,0.214,0.09,0.033,0.098,0.041,0.225,0.037,0.128,0.162,0.015,0.134,0.971,0.238,0.09]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (4, '[0.105,0.199,0.225,0.059,0.961,0.097,0.233,0.178,0.12,0.249,0.166,0.041,0.002,0.07,0.17,0.004,0.14,0.125,0.078,0.155,0.178,0.965,0.102,0.248,0.16,0.075,0.027,0.071,0.073,0.237,0.162,0.115,0.213,0.018,0.987,0.21,0.062,0.043,0.043,0.155,0.052,0.238,0.111,0.019,0.062,0.022,0.061,0.002,0.232,0.207,0.961,0.071,0.11,0.164,0.152,0.129,0.093,0.196,0.129,0.148,0.187,0.07,0.089,0.224]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (5, '[0.186,0.026,0.039,0.228,0.156,0.101,0.019,0.236,0.109,0.221,0.242,0.988,0.077,0.111,0.129,0.055,0.221,0.177,0.128,0.151,0.028,0.238,0.125,0.188,0.05,0.026,0.064,0.149,0.973,0.047,0.077,0.053,0.057,0.074,0.022,0.134,0.223,0.142,0.147,0.142,0.009,0.977,0.052,0.05,0.047,0.203,0.041,0.228,0.195,0.11,0.06,0.225,0.093,0.191,0.003,0.131,0.026,0.971,0.118,0.011,0.045,0.114,0.217,0.149]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (6, '[0.992,0.023,0.13,0.144,0.097,0.098,0.089,0.178,0.041,0.248,0.132,0.168,0.085,0.058,0.178,0.21,0.024,0.009,0.984,0.14,0.017,0.112,0.064,0.142,0.224,0.083,0.17,0.059,0.015,0.187,0.044,0.201,0.192,0.121,0.218,0.983,0.122,0.13,0.242,0.077,0.153,0.095,0.023,0.107,0.001,0.053,0.004,0.211,0.991,0.148,0.178,0.173,0.18,0.018,0.051,0.027,0.142,0.056,0.033,0.078,0.17,0.001,0.171,0.101]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (7, '[0.139,0.215,0.0,0.025,0.12,0.004,0.052,1.0,0.029,0.111,0.115,0.106,0.234,0.065,0.151,0.194,0.009,0.093,0.012,0.2,0.097,0.016,0.246,0.083,0.042,1.0,0.038,0.132,0.07,0.051,0.083,0.02,0.034,0.235,0.141,0.1,0.124,0.186,0.101,0.213,0.191,0.001,0.991,0.114,0.199,0.236,0.012,0.09,0.033,0.224,0.089,0.122,0.084,0.027,0.094,0.987,0.175,0.173,0.024,0.161,0.191,0.074,0.027,0.072]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (8, '[0.018,0.101,0.104,0.16,0.055,0.169,0.002,0.046,0.124,0.048,0.168,0.037,0.054,0.122,0.959,0.142,0.205,0.033,0.023,0.025,0.148,0.08,0.12,0.129,0.189,0.014,0.116,0.16,0.165,0.234,0.167,0.242,0.971,0.103,0.197,0.082,0.214,0.017,0.212,0.182,0.058,0.175,0.073,0.139,0.119,0.026,0.234,0.192,0.073,0.969,0.071,0.162,0.109,0.214,0.083,0.046,0.208,0.141,0.225,0.019,0.005,0.055,0.957,0.033]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (9, '[0.081,0.1,0.126,0.204,0.134,0.98,0.112,0.005,0.055,0.205,0.112,0.058,0.244,0.042,0.13,0.152,0.164,0.127,0.079,0.231,0.083,0.974,0.052,0.006,0.015,0.233,0.225,0.209,0.145,0.028,0.075,0.009,0.165,0.012,0.142,0.103,0.024,0.125,0.085,0.977,0.111,0.234,0.109,0.173,0.041,0.174,0.234,0.038,0.12,0.006,0.242,0.117,0.238,0.023,0.026,0.071,0.98,0.134,0.238,0.073,0.047,0.217,0.028,0.059]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (10, '[0.016,0.008,0.222,0.078,0.12,0.023,0.176,0.065,0.179,0.13,0.223,0.072,0.954,0.032,0.018,0.017,0.093,0.102,0.01,0.175,0.223,0.008,0.12,0.221,0.166,0.03,0.094,0.038,0.999,0.038,0.125,0.139,0.143,0.025,0.037,0.062,0.196,0.179,0.239,0.204,0.071,0.203,0.131,0.134,0.204,0.232,0.959,0.169,0.159,0.189,0.113,0.132,0.228,0.223,0.063,0.083,0.223,0.004,0.229,0.025,0.113,0.004,0.026,0.993]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (11, '[0.212,0.15,0.108,0.081,0.039,0.237,0.954,0.199,0.077,0.138,0.115,0.241,0.02,0.083,0.081,0.049,0.161,0.036,0.032,0.981,0.125,0.184,0.043,0.209,0.115,0.071,0.135,0.103,0.193,0.035,0.078,0.101,0.068,0.003,0.022,0.995,0.215,0.208,0.056,0.24,0.219,0.064,0.102,0.055,0.049,0.066,0.235,0.13,0.057,0.069,0.048,0.234,0.031,0.965,0.093,0.205,0.025,0.025,0.141,0.249,0.062,0.242,0.127,0.24]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (12, '[0.087,0.166,0.13,0.113,0.034,0.05,0.196,0.043,0.168,0.081,0.078,0.069,0.048,0.985,0.023,0.125,0.021,0.241,0.129,0.05,0.049,0.062,0.2,0.101,0.118,0.125,0.994,0.007,0.027,0.046,0.178,0.109,0.065,0.121,0.068,0.202,0.157,0.036,0.165,0.201,0.004,0.194,0.989,0.051,0.113,0.23,0.068,0.039,0.102,0.077,0.006,0.175,0.224,0.184,0.228,0.214,0.119,0.102,0.006,0.215,0.957,0.178,0.083,0.239]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (13, '[0.174,0.019,0.05,0.97,0.219,0.014,0.131,0.207,0.13,0.229,0.157,0.231,0.064,0.218,0.105,0.041,0.234,0.236,0.145,0.186,0.966,0.151,0.007,0.06,0.109,0.244,0.19,0.007,0.001,0.028,0.089,0.043,0.184,0.977,0.223,0.177,0.136,0.114,0.228,0.198,0.148,0.02,0.023,0.112,0.19,0.243,0.074,0.231,0.168,0.997,0.171,0.044,0.152,0.167,0.227,0.025,0.016,0.088,0.059,0.086,0.102,0.202,0.182,0.1]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (14, '[0.236,0.11,0.178,0.204,0.09,0.055,0.24,0.013,0.006,0.136,0.99,0.029,0.036,0.119,0.049,0.01,0.014,0.13,0.114,0.142,0.172,0.228,0.099,0.122,0.047,0.041,0.087,0.965,0.134,0.078,0.123,0.092,0.001,0.079,0.181,0.201,0.212,0.246,0.042,0.238,0.978,0.119,0.218,0.123,0.044,0.127,0.088,0.232,0.108,0.173,0.013,0.176,0.06,0.162,0.221,0.09,0.956,0.018,0.209,0.019,0.222,0.053,0.117,0.069]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (15, '[0.189,0.072,0.025,0.196,0.179,0.229,0.171,0.224,0.245,0.117,0.153,0.029,0.224,0.24,0.194,0.134,0.193,0.977,0.217,0.206,0.041,0.156,0.096,0.07,0.17,0.006,0.099,0.022,0.017,0.182,0.105,0.013,0.241,0.138,0.966,0.248,0.212,0.065,0.055,0.115,0.186,0.194,0.18,0.15,0.016,0.006,0.232,0.977,0.223,0.248,0.033,0.073,0.165,0.138,0.147,0.067,0.143,0.082,0.104,0.011,0.048,0.051,0.141,0.996]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (16, '[0.076,0.15,0.09,0.023,0.149,0.04,0.987,0.133,0.098,0.068,0.055,0.096,0.027,0.216,0.041,0.249,0.022,0.169,0.124,0.064,0.203,0.077,0.096,0.016,0.998,0.175,0.085,0.158,0.117,0.173,0.108,0.163,0.17,0.205,0.221,0.243,0.221,0.138,0.022,0.016,0.143,0.965,0.021,0.211,0.171,0.092,0.206,0.036,0.2,0.174,0.14,0.097,0.239,0.218,0.96,0.039,0.018,0.109,0.013,0.052,0.232,0.058,0.007,0.152]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (17, '[0.218,0.057,0.177,0.219,0.245,0.222,0.103,0.025,0.248,0.079,0.059,0.241,0.161,0.988,0.038,0.246,0.078,0.124,0.046,0.153,0.236,0.151,0.008,0.073,0.084,0.065,0.247,0.113,0.093,0.053,0.161,0.981,0.113,0.189,0.058,0.179,0.237,0.095,0.184,0.13,0.01,0.024,0.133,0.142,0.206,0.166,0.012,0.2,0.995,0.14,0.099,0.012,0.165,0.044,0.027,0.205,0.085,0.045,0.133,0.231,0.122,0.96,0.161,0.123]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (18, '[0.108,0.016,0.131,0.229,0.957,0.103,0.062,0.015,0.061,0.073,0.042,0.04,0.077,0.13,0.115,0.094,0.033,0.161,0.225,0.056,0.976,0.167,0.043,0.024,0.016,0.116,0.139,0.105,0.122,0.051,0.013,0.011,0.072,0.138,0.22,0.174,0.008,0.033,0.953,0.2,0.045,0.023,0.126,0.048,0.045,0.013,0.221,0.069,0.063,0.204,0.16,0.136,0.22,0.108,0.149,0.95,0.02,0.178,0.223,0.065,0.07,0.013,0.145,0.157]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (19, '[0.077,0.213,0.019,0.053,0.065,0.196,0.071,0.089,0.032,0.177,0.071,0.995,0.024,0.099,0.085,0.142,0.175,0.197,0.183,0.225,0.046,0.144,0.144,0.171,0.116,0.165,0.012,0.96,0.183,0.204,0.19,0.149,0.086,0.108,0.011,0.156,0.144,0.03,0.09,0.242,0.113,0.15,0.156,0.187,0.012,0.971,0.102,0.151,0.215,0.012,0.062,0.224,0.05,0.133,0.096,0.115,0.013,0.065,0.093,0.241,0.235,0.137,0.991,0.177]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (20, '[0.056,0.008,0.171,0.22,0.169,0.985,0.145,0.063,0.199,0.184,0.157,0.199,0.249,0.033,0.174,0.2,0.233,0.162,0.954,0.076,0.189,0.016,0.158,0.244,0.249,0.046,0.1,0.142,0.132,0.137,0.107,0.246,0.014,0.019,0.97,0.05,0.232,0.192,0.194,0.137,0.018,0.203,0.202,0.099,0.031,0.143,0.064,0.165,0.175,0.126,0.112,0.013,0.973,0.184,0.005,0.201,0.237,0.178,0.109,0.166,0.052,0.022,0.198,0.137]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (21, '[0.181,0.22,0.024,0.095,0.036,0.215,0.188,0.012,0.18,0.009,0.199,0.133,0.978,0.073,0.198,0.113,0.083,0.024,0.227,0.194,0.077,0.144,0.09,0.13,0.132,0.956,0.045,0.212,0.234,0.047,0.065,0.045,0.026,0.206,0.006,0.158,0.09,0.213,0.014,0.047,0.14,0.975,0.225,0.07,0.046,0.05,0.141,0.135,0.14,0.053,0.103,0.057,0.058,0.161,0.097,0.247,0.192,0.188,0.23,0.994,0.233,0.158,0.221,0.052]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (22, '[0.11,0.247,0.976,0.236,0.028,0.04,0.158,0.085,0.173,0.046,0.036,0.028,0.238,0.237,0.136,0.036,0.099,0.184,0.197,0.961,0.049,0.057,0.049,0.059,0.224,0.121,0.245,0.204,0.242,0.016,0.182,0.001,0.956,0.2,0.206,0.178,0.165,0.121,0.202,0.034,0.101,0.002,0.091,0.167,0.16,0.106,0.139,0.118,0.996,0.143,0.031,0.202,0.169,0.219,0.143,0.063,0.183,0.085,0.1,0.239,0.09,0.212,0.01,0.179]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (23, '[0.02,0.068,0.073,0.093,0.249,0.15,0.18,0.143,0.227,0.968,0.232,0.034,0.092,0.04,0.192,0.158,0.19,0.123,0.056,0.218,0.182,0.172,0.079,0.067,0.045,0.229,0.971,0.006,0.22,0.219,0.148,0.058,0.244,0.226,0.0,0.011,0.063,0.025,0.11,0.983,0.135,0.035,0.158,0.19,0.172,0.171,0.004,0.027,0.013,0.162,0.005,0.213,0.166,0.095,0.249,0.98,0.055,0.044,0.054,0.219,0.24,0.122,0.247,0.14]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (24, '[0.135,0.143,0.003,0.107,0.089,0.001,0.036,0.13,0.217,0.13,0.063,0.234,0.119,0.087,0.114,0.168,0.976,0.217,0.158,0.223,0.138,0.18,0.227,0.222,0.232,0.104,0.046,0.072,0.237,0.005,0.02,0.044,0.183,0.954,0.088,0.098,0.001,0.018,0.245,0.235,0.143,0.09,0.043,0.052,0.103,0.123,0.953,0.117,0.049,0.124,0.023,0.008,0.111,0.156,0.248,0.023,0.064,0.195,0.01,0.002,0.057,0.235,0.956,0.1]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (25, '[0.1,0.098,0.11,0.221,0.098,0.977,0.214,0.157,0.227,0.203,0.133,0.158,0.035,0.125,0.116,0.191,0.162,0.198,0.055,0.095,0.008,0.152,0.128,0.966,0.079,0.08,0.043,0.102,0.01,0.166,0.069,0.246,0.11,0.079,0.082,0.228,0.047,0.179,0.02,0.049,0.959,0.197,0.246,0.008,0.055,0.05,0.113,0.25,0.074,0.005,0.001,0.014,0.164,0.996,0.247,0.001,0.13,0.013,0.157,0.179,0.203,0.2,0.118,0.22]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (26, '[0.127,0.032,0.208,0.187,0.23,0.104,0.168,0.029,0.105,0.122,0.219,0.021,0.961,0.01,0.101,0.218,0.033,0.013,0.02,0.082,0.013,0.144,0.215,0.107,0.237,0.19,0.067,0.197,0.068,0.015,0.982,0.032,0.158,0.138,0.092,0.152,0.151,0.116,0.151,0.002,0.146,0.205,0.04,0.043,0.049,0.115,0.048,0.967,0.242,0.027,0.166,0.1,0.125,0.079,0.21,0.048,0.205,0.056,0.232,0.206,0.975,0.182,0.17,0.117]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (27, '[0.029,0.237,0.133,0.969,0.145,0.056,0.003,0.154,0.169,0.004,0.217,0.119,0.016,0.221,0.184,0.22,0.039,0.077,0.02,0.98,0.123,0.011,0.097,0.052,0.193,0.177,0.047,0.021,0.19,0.19,0.145,0.185,0.063,0.222,0.17,0.15,0.107,0.971,0.021,0.206,0.173,0.166,0.136,0.112,0.213,0.084,0.109,0.147,0.226,0.032,0.083,0.163,0.224,0.103,0.952,0.011,0.225,0.192,0.16,0.113,0.2,0.142,0.247,0.04]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (28, '[0.003,0.109,0.071,0.004,0.113,0.044,0.191,0.085,0.208,0.235,0.954,0.216,0.113,0.075,0.232,0.023,0.049,0.133,0.114,0.214,0.133,0.193,0.134,0.191,0.072,0.068,0.966,0.244,0.14,0.247,0.096,0.007,0.122,0.163,0.197,0.006,0.063,0.165,0.009,0.058,0.206,0.233,0.212,0.106,0.951,0.153,0.155,0.129,0.135,0.09,0.103,0.1,0.163,0.122,0.153,0.06,0.074,0.001,0.117,0.193,0.028,0.95,0.228,0.174]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (29, '[0.179,0.043,0.089,0.041,0.992,0.194,0.23,0.218,0.034,0.069,0.11,0.243,0.051,0.115,0.157,0.002,0.086,0.966,0.052,0.135,0.167,0.071,0.223,0.058,0.064,0.083,0.012,0.027,0.139,0.09,0.144,0.239,0.068,0.957,0.165,0.108,0.038,0.027,0.186,0.137,0.177,0.171,0.046,0.092,0.089,0.194,0.063,0.144,0.071,0.085,0.142,0.989,0.098,0.158,0.124,0.147,0.221,0.104,0.195,0.245,0.082,0.029,0.061,0.043]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (30, '[0.148,0.148,0.078,0.103,0.199,0.043,0.19,0.184,0.066,0.201,0.164,0.988,0.018,0.149,0.205,0.012,0.138,0.215,0.149,0.067,0.149,0.217,0.068,0.135,0.989,0.201,0.243,0.125,0.248,0.023,0.024,0.074,0.009,0.058,0.163,0.166,0.005,0.036,0.19,0.022,0.963,0.229,0.106,0.139,0.168,0.159,0.08,0.097,0.118,0.012,0.053,0.17,0.073,0.002,0.155,0.168,0.173,0.083,0.962,0.176,0.083,0.035,0.194,0.172]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (31, '[0.246,0.951,0.047,0.189,0.02,0.224,0.06,0.024,0.202,0.173,0.206,0.246,0.184,0.127,0.156,0.088,0.068,0.043,0.995,0.12,0.01,0.064,0.218,0.091,0.112,0.031,0.2,0.162,0.196,0.14,0.045,0.966,0.131,0.121,0.016,0.063,0.218,0.184,0.124,0.014,0.01,0.233,0.027,0.037,0.149,0.002,0.164,0.994,0.029,0.162,0.147,0.207,0.112,0.064,0.168,0.181,0.156,0.028,0.134,0.082,0.21,0.214,0.233,0.133]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (32, '[0.049,0.032,0.098,0.138,0.006,0.103,0.193,0.232,0.989,0.04,0.009,0.199,0.051,0.153,0.174,0.188,0.022,0.042,0.248,0.098,0.066,0.156,0.239,0.172,0.064,0.95,0.165,0.244,0.24,0.22,0.21,0.046,0.142,0.155,0.234,0.128,0.047,0.188,0.999,0.219,0.073,0.122,0.151,0.058,0.148,0.111,0.012,0.215,0.038,0.018,0.184,0.158,0.141,0.194,0.971,0.003,0.231,0.176,0.123,0.001,0.061,0.241,0.12,0.062]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (33, '[0.158,0.052,0.165,0.227,0.06,0.212,0.006,0.127,0.042,0.239,0.098,0.097,0.112,0.171,0.184,0.97,0.022,0.144,0.22,0.013,0.153,0.091,0.072,0.027,0.141,0.15,0.185,0.227,0.09,0.219,0.238,0.022,0.962,0.068,0.136,0.001,0.111,0.031,0.064,0.029,0.014,0.172,0.172,0.097,0.18,0.959,0.045,0.068,0.098,0.088,0.208,0.028,0.119,0.148,0.06,0.006,0.183,0.084,0.003,0.008,0.07,0.951,0.124,0.113]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (34, '[0.072,0.088,0.11,0.059,0.984,0.121,0.118,0.078,0.003,0.168,0.117,0.136,0.132,0.006,0.124,0.096,0.072,0.226,0.124,0.029,0.006,0.212,0.993,0.085,0.178,0.165,0.035,0.106,0.208,0.115,0.154,0.069,0.024,0.053,0.184,0.152,0.201,0.222,0.091,0.96,0.24,0.088,0.146,0.133,0.157,0.091,0.083,0.207,0.158,0.099,0.176,0.164,0.961,0.145,0.001,0.221,0.176,0.041,0.076,0.114,0.084,0.182,0.216,0.148]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (35, '[0.113,0.088,0.18,0.219,0.003,0.213,0.129,0.04,0.011,0.13,0.198,0.991,0.195,0.009,0.119,0.181,0.231,0.009,0.096,0.068,0.154,0.115,0.051,0.228,0.221,0.143,0.247,0.238,0.094,0.969,0.105,0.081,0.01,0.149,0.115,0.011,0.187,0.08,0.045,0.072,0.128,0.212,0.19,0.184,0.113,0.157,0.956,0.154,0.098,0.108,0.199,0.059,0.08,0.202,0.069,0.091,0.243,0.063,0.137,0.957,0.166,0.168,0.188,0.01]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (36, '[0.004,0.129,0.952,0.1,0.156,0.225,0.101,0.139,0.081,0.006,0.215,0.148,0.044,0.067,0.165,0.092,0.086,0.22,0.959,0.087,0.12,0.032,0.159,0.168,0.003,0.174,0.051,0.067,0.174,0.073,0.093,0.119,0.01,0.03,0.202,0.064,0.954,0.034,0.186,0.018,0.123,0.128,0.095,0.144,0.211,0.19,0.104,0.059,0.072,0.002,0.171,0.047,0.246,0.988,0.052,0.244,0.21,0.153,0.224,0.015,0.155,0.211,0.093,0.037]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (37, '[0.123,0.097,0.105,0.239,0.097,0.206,0.075,0.021,0.137,0.988,0.203,0.017,0.128,0.048,0.067,0.221,0.041,0.046,0.166,0.016,0.053,0.053,0.009,0.221,0.138,0.987,0.221,0.122,0.004,0.216,0.213,0.167,0.184,0.141,0.158,0.058,0.199,0.147,0.197,0.227,0.105,0.201,0.121,0.954,0.104,0.089,0.133,0.196,0.183,0.103,0.032,0.078,0.128,0.173,0.031,0.047,0.052,0.242,0.062,0.054,0.996,0.245,0.06,0.054]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (38, '[0.159,0.076,0.019,0.955,0.108,0.141,0.22,0.003,0.005,0.01,0.054,0.191,0.207,0.057,0.183,0.154,0.984,0.085,0.238,0.164,0.129,0.227,0.219,0.116,0.117,0.192,0.065,0.226,0.134,0.216,0.154,0.213,0.974,0.061,0.081,0.18,0.205,0.098,0.102,0.004,0.204,0.123,0.088,0.081,0.133,0.093,0.089,0.031,0.164,0.1,0.968,0.028,0.203,0.146,0.177,0.143,0.222,0.162,0.191,0.132,0.241,0.171,0.13,0.237]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (39, '[0.201,0.159,0.187,0.033,0.156,0.215,0.083,0.191,0.122,0.055,0.957,0.116,0.037,0.048,0.116,0.099,0.035,0.116,0.16,0.121,0.205,0.098,0.123,0.967,0.194,0.125,0.079,0.111,0.032,0.009,0.075,0.074,0.162,0.087,0.135,0.094,0.026,0.007,0.155,0.996,0.156,0.205,0.026,0.185,0.129,0.066,0.117,0.05,0.186,0.021,0.18,0.041,0.092,0.229,0.09,0.087,0.151,0.96,0.099,0.237,0.161,0.236,0.001,0.131]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (40, '[0.952,0.207,0.204,0.143,0.091,0.001,0.072,0.169,0.144,0.052,0.019,0.111,0.235,0.019,0.052,0.223,0.021,0.978,0.168,0.092,0.236,0.225,0.027,0.036,0.166,0.221,0.144,0.074,0.054,0.138,0.968,0.061,0.211,0.043,0.226,0.019,0.071,0.031,0.025,0.232,0.118,0.153,0.188,0.246,0.128,0.164,0.99,0.159,0.246,0.101,0.031,0.035,0.048,0.199,0.089,0.233,0.177,0.143,0.095,0.131,0.088,0.006,0.067,0.065]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (41, '[0.118,0.023,0.223,0.091,0.053,0.19,0.233,0.996,0.01,0.078,0.234,0.159,0.202,0.23,0.158,0.125,0.167,0.082,0.112,0.219,0.089,0.131,0.129,0.126,0.978,0.006,0.083,0.121,0.094,0.111,0.061,0.214,0.014,0.2,0.244,0.014,0.057,0.979,0.185,0.035,0.036,0.109,0.167,0.113,0.015,0.043,0.006,0.048,0.127,0.033,0.168,0.025,0.248,0.963,0.199,0.167,0.11,0.158,0.083,0.117,0.137,0.024,0.176,0.211]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (42, '[0.078,0.043,0.018,0.233,0.025,0.096,0.246,0.234,0.042,0.179,0.245,0.061,0.037,0.168,0.974,0.236,0.014,0.179,0.138,0.108,0.153,0.027,0.184,0.113,0.197,0.074,0.071,0.205,0.008,0.193,0.019,0.994,0.216,0.078,0.137,0.019,0.156,0.149,0.165,0.048,0.016,0.135,0.146,0.054,0.994,0.19,0.018,0.046,0.08,0.101,0.128,0.203,0.172,0.137,0.023,0.05,0.112,0.116,0.069,0.065,0.96,0.111,0.095,0.13]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (43, '[0.177,0.128,0.248,0.955,0.167,0.22,0.168,0.132,0.037,0.153,0.122,0.058,0.157,0.232,0.038,0.198,0.053,0.154,0.029,0.084,0.08,0.966,0.154,0.079,0.002,0.126,0.184,0.128,0.13,0.155,0.212,0.188,0.095,0.094,0.147,0.153,0.21,0.165,0.979,0.026,0.209,0.156,0.079,0.184,0.042,0.023,0.001,0.029,0.221,0.093,0.176,0.987,0.147,0.221,0.237,0.121,0.092,0.081,0.155,0.133,0.069,0.028,0.249,0.072]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (44, '[0.081,0.184,0.101,0.018,0.167,0.132,0.222,0.115,0.198,0.195,0.994,0.159,0.229,0.188,0.04,0.122,0.071,0.127,0.095,0.21,0.032,0.215,0.21,0.125,0.093,0.026,0.187,0.129,0.957,0.138,0.183,0.204,0.181,0.202,0.029,0.15,0.037,0.094,0.133,0.163,0.222,0.054,0.046,0.177,0.2,0.999,0.195,0.199,0.156,0.096,0.141,0.1,0.227,0.071,0.069,0.172,0.125,0.142,0.967,0.178,0.202,0.04,0.076,0.086]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (45, '[0.148,0.961,0.243,0.078,0.184,0.199,0.153,0.092,0.063,0.062,0.185,0.077,0.038,0.027,0.117,0.147,0.036,0.967,0.234,0.135,0.19,0.235,0.149,0.112,0.249,0.18,0.044,0.238,0.035,0.074,0.083,0.069,0.201,0.229,0.165,0.963,0.191,0.186,0.082,0.214,0.159,0.031,0.181,0.243,0.086,0.083,0.008,0.159,0.191,0.209,0.201,0.157,0.981,0.225,0.16,0.006,0.15,0.132,0.143,0.056,0.084,0.023,0.072,0.1]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (46, '[0.074,0.116,0.145,0.007,0.234,0.161,0.187,0.107,0.98,0.182,0.134,0.234,0.076,0.028,0.208,0.154,0.22,0.147,0.171,0.085,0.118,0.145,0.016,0.06,0.961,0.077,0.244,0.075,0.001,0.02,0.235,0.161,0.188,0.141,0.175,0.039,0.113,0.178,0.055,0.173,0.04,0.013,0.961,0.148,0.059,0.046,0.156,0.134,0.215,0.146,0.157,0.048,0.238,0.149,0.019,0.106,0.037,0.0,0.207,0.952,0.18,0.196,0.023,0.021]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (47, '[0.082,0.109,0.966,0.083,0.239,0.12,0.222,0.228,0.019,0.067,0.079,0.112,0.131,0.156,0.205,0.952,0.085,0.016,0.002,0.013,0.171,0.213,0.162,0.153,0.043,0.026,0.051,0.037,0.028,0.087,0.054,0.99,0.015,0.021,0.048,0.003,0.126,0.144,0.057,0.088,0.172,0.146,0.062,0.244,0.099,0.187,0.148,0.194,0.011,0.962,0.104,0.208,0.04,0.206,0.236,0.145,0.092,0.096,0.028,0.028,0.031,0.125,0.227,0.224]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (48, '[0.015,0.022,0.081,0.23,0.135,0.163,0.192,0.127,0.231,0.996,0.209,0.169,0.008,0.109,0.068,0.159,0.11,0.11,0.217,0.154,0.117,0.146,0.954,0.226,0.119,0.25,0.026,0.01,0.114,0.221,0.131,0.116,0.104,0.2,0.233,0.047,0.167,0.062,0.972,0.021,0.022,0.09,0.211,0.199,0.219,0.21,0.041,0.18,0.142,0.054,0.143,0.236,0.13,0.071,0.057,0.023,0.975,0.051,0.075,0.034,0.147,0.248,0.017,0.058]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (49, '[0.176,0.219,0.237,0.128,0.084,0.125,0.129,0.178,0.082,0.062,0.194,0.225,0.031,0.16,0.198,0.207,0.989,0.134,0.173,0.217,0.098,0.239,0.244,0.204,0.25,0.177,0.175,0.007,0.157,0.985,0.154,0.076,0.01,0.089,0.055,0.131,0.008,0.12,0.114,0.237,0.179,0.195,0.1,0.124,0.052,0.964,0.208,0.151,0.033,0.016,0.027,0.097,0.131,0.128,0.048,0.08,0.173,0.049,0.204,0.168,0.163,0.132,0.146,0.993]'::vector);

INSERT INTO user_mcp_profiles (user_id, effective_vector)
VALUES (50, '[0.194,0.236,0.111,0.172,0.162,0.225,0.951,0.216,0.104,0.087,0.236,0.071,0.042,0.165,0.017,0.097,0.154,0.219,0.123,0.005,0.049,0.135,0.037,0.956,0.026,0.036,0.037,0.155,0.124,0.15,0.031,0.062,0.024,0.18,0.122,0.016,0.97,0.006,0.198,0.216,0.125,0.017,0.235,0.147,0.208,0.147,0.042,0.195,0.12,0.152,0.094,0.19,0.989,0.099,0.082,0.127,0.08,0.142,0.088,0.021,0.197,0.188,0.106,0.222]'::vector);







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

INSERT INTO team_members (team_id, user_id, role) VALUES
(1,1,'Leader'),
(1,2,'Specialist'),
(1,3,'Specialist'),
(1,4,'Operator'),
(1,5,'Coordinator'),

(2,6,'Leader'),
(2,7,'Visionary'),
(2,8,'Specialist'),
(2,9,'Operator'),
(2,10,'Coordinator'),

(3,11,'Leader'),
(3,12,'Specialist'),
(3,13,'Operator'),
(3,14,'Coordinator'),
(3,15,'Executor'),

(4,16,'Leader'),
(4,17,'Operator'),
(4,18,'Specialist'),
(4,19,'Coordinator'),
(4,20,'Executor'),

(5,21,'Leader'),
(5,22,'Visionary'),
(5,23,'Specialist'),
(5,24,'Operator'),
(5,25,'Coordinator'),

(6,26,'Leader'),
(6,27,'Creative'),
(6,28,'Specialist'),
(6,29,'Operator'),
(6,30,'Coordinator'),

(7,31,'Leader'),
(7,32,'Strategist'),
(7,33,'Specialist'),
(7,34,'Operator'),
(7,35,'Coordinator');