-- ============================================
-- USERS
-- ============================================

INSERT INTO users (name, email, password_hash)
VALUES
('Alberto Gonzalez', 'alberto@mutari.dev',crypt('123456', gen_salt('bf')));

-- ============================================
-- COMPANIES
-- ============================================

INSERT INTO companies (name, website, owner)
VALUES (
  'Mutari',
  'https://mutari.ai',
  (SELECT id FROM users WHERE email = 'alberto@mutari.dev')
);

INSERT INTO company_users (company_id, user_id, role, is_default)
VALUES (
  (SELECT id FROM companies WHERE name = 'Mutari'),
  (SELECT id FROM users WHERE email = 'alberto@mutari.dev'),
  'owner',
  TRUE
);

-- ============================================
-- COMPANY MEMBERS
-- ============================================

INSERT INTO company_members (company_id, name, email, role)
VALUES
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Alice Johnson',
  'alice@mutari.com',
  'Engineer'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Bob Smith',
  'bob@mutari.com',
  'Designer'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Carol White',
  'carol@mutari.com',
  'Product Manager'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'David Brown',
  'david@mutari.com',
  'Engineer'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Emma Davis',
  'emma@mutari.com',
  'QA Engineer'
);

-- ============================================
-- TEAMS
-- ============================================

INSERT INTO teams (name, description, company_id)
VALUES
(
  'Core Product Team',
  'Main product builders',
  (SELECT id FROM companies WHERE name = 'Mutari')
),
(
  'Design Team',
  'UX/UI and design',
  (SELECT id FROM companies WHERE name = 'Mutari')
);

-- ============================================
-- TEAM MEMBERS
-- ============================================

-- Core Team
INSERT INTO team_members (team_id, company_member_id)
VALUES
(
  (SELECT id FROM teams WHERE name = 'Core Product Team'),
  (SELECT id FROM company_members WHERE email = 'alice@mutari.com')
),
(
  (SELECT id FROM teams WHERE name = 'Core Product Team'),
  (SELECT id FROM company_members WHERE email = 'carol@mutari.com')
),
(
  (SELECT id FROM teams WHERE name = 'Core Product Team'),
  (SELECT id FROM company_members WHERE email = 'david@mutari.com')
);

-- Design Team
INSERT INTO team_members (team_id, company_member_id)
VALUES
(
  (SELECT id FROM teams WHERE name = 'Design Team'),
  (SELECT id FROM company_members WHERE email = 'bob@mutari.com')
),
(
  (SELECT id FROM teams WHERE name = 'Design Team'),
  (SELECT id FROM company_members WHERE email = 'emma@mutari.com')
);

-- ============================================
-- CANDIDATES
-- ============================================

INSERT INTO candidates (company_id, name, email, status)
VALUES
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Ethan Hunt',
  'ethan@gmail.com',
  'active'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Sophia Turner',
  'sophia@gmail.com',
  'active'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Liam Walker',
  'liam@gmail.com',
  'active'
),
(
  (SELECT id FROM companies WHERE name = 'Mutari'),
  'Olivia Harris',
  'olivia@gmail.com',
  'archived'
);

-- ============================================
-- TOKENS (para probar quiz después)
-- ============================================

INSERT INTO mcp_tokens (entity_type, entity_id, token, expires_at)
VALUES
(
  'candidate',
  (SELECT id FROM candidates WHERE email = 'ethan@gmail.com'),
  'demo-token-1',
  NOW() + INTERVAL '1 day'
),
(
  'candidate',
  (SELECT id FROM candidates WHERE email = 'sophia@gmail.com'),
  'demo-token-2',
  NOW() + INTERVAL '1 day'
),
(
  'candidate',
  (SELECT id FROM candidates WHERE email = 'liam@gmail.com'),
  'demo-token-3',
  NOW() + INTERVAL '1 day'
);
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
