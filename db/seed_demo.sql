
-- ============================================
-- RESET
-- ============================================

TRUNCATE team_members CASCADE;
TRUNCATE teams CASCADE;
TRUNCATE user_mcp_profiles CASCADE;
TRUNCATE user_mcp_answers CASCADE;
TRUNCATE users CASCADE;

-- ============================================
-- USERS (30 demo candidates)
-- ============================================

INSERT INTO users (id,name,email,password_hash,photo_url) VALUES
(1,'Visionary Founder','demo1@mutari.io','demo','https://randomuser.me/api/portraits/men/1.jpg'),
(2,'Product Strategist','demo2@mutari.io','demo','https://randomuser.me/api/portraits/men/2.jpg'),
(3,'Engineering Architect','demo3@mutari.io','demo','https://randomuser.me/api/portraits/men/3.jpg'),
(4,'Creative Director','demo4@mutari.io','demo','https://randomuser.me/api/portraits/women/4.jpg'),
(5,'Operations Leader','demo5@mutari.io','demo','https://randomuser.me/api/portraits/men/5.jpg'),
(6,'UX Researcher','demo6@mutari.io','demo','https://randomuser.me/api/portraits/women/6.jpg'),
(7,'Growth Hacker','demo7@mutari.io','demo','https://randomuser.me/api/portraits/men/7.jpg'),
(8,'Financial Analyst','demo8@mutari.io','demo','https://randomuser.me/api/portraits/women/8.jpg'),
(9,'Customer Success','demo9@mutari.io','demo','https://randomuser.me/api/portraits/women/9.jpg'),
(10,'Data Scientist','demo10@mutari.io','demo','https://randomuser.me/api/portraits/men/10.jpg'),

(11,'Strategic Thinker','demo11@mutari.io','demo','https://randomuser.me/api/portraits/men/11.jpg'),
(12,'Creative Builder','demo12@mutari.io','demo','https://randomuser.me/api/portraits/women/12.jpg'),
(13,'Systems Engineer','demo13@mutari.io','demo','https://randomuser.me/api/portraits/men/13.jpg'),
(14,'Design Visionary','demo14@mutari.io','demo','https://randomuser.me/api/portraits/women/14.jpg'),
(15,'Startup Operator','demo15@mutari.io','demo','https://randomuser.me/api/portraits/men/15.jpg'),
(16,'People Connector','demo16@mutari.io','demo','https://randomuser.me/api/portraits/women/16.jpg'),
(17,'Execution Specialist','demo17@mutari.io','demo','https://randomuser.me/api/portraits/men/17.jpg'),
(18,'Market Analyst','demo18@mutari.io','demo','https://randomuser.me/api/portraits/women/18.jpg'),
(19,'Culture Builder','demo19@mutari.io','demo','https://randomuser.me/api/portraits/women/19.jpg'),
(20,'Technical Innovator','demo20@mutari.io','demo','https://randomuser.me/api/portraits/men/20.jpg'),

(21,'Product Visionary','demo21@mutari.io','demo','https://randomuser.me/api/portraits/men/21.jpg'),
(22,'Growth Strategist','demo22@mutari.io','demo','https://randomuser.me/api/portraits/women/22.jpg'),
(23,'Platform Architect','demo23@mutari.io','demo','https://randomuser.me/api/portraits/men/23.jpg'),
(24,'Creative Technologist','demo24@mutari.io','demo','https://randomuser.me/api/portraits/women/24.jpg'),
(25,'Operations Strategist','demo25@mutari.io','demo','https://randomuser.me/api/portraits/men/25.jpg'),
(26,'Community Builder','demo26@mutari.io','demo','https://randomuser.me/api/portraits/women/26.jpg'),
(27,'Execution Leader','demo27@mutari.io','demo','https://randomuser.me/api/portraits/men/27.jpg'),
(28,'Insight Analyst','demo28@mutari.io','demo','https://randomuser.me/api/portraits/women/28.jpg'),
(29,'Innovation Catalyst','demo29@mutari.io','demo','https://randomuser.me/api/portraits/women/29.jpg'),
(30,'Engineering Leader','demo30@mutari.io','demo','https://randomuser.me/api/portraits/men/30.jpg');

-- ============================================
-- TEAMS
-- ============================================

INSERT INTO teams (id,name) VALUES
(1,'Startup Founders'),
(2,'Product Team'),
(3,'Engineering Team');

-- ============================================
-- TEAM MEMBERS
-- ============================================

INSERT INTO team_members (team_id,user_id,role) VALUES
(1,1,'Leader'),
(1,2,'Strategist'),
(1,5,'Operator'),

(2,2,'Leader'),
(2,4,'Designer'),
(2,6,'Research'),

(3,3,'Architect'),
(3,10,'Data'),
(3,20,'Engineer');

-- ============================================
-- ARCHETYPE VECTOR GENERATION (64 dimensions)
-- ============================================

INSERT INTO user_mcp_profiles(user_id,effective_vector)
SELECT id,
ARRAY(
0.6 + random()*0.2,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,

0.6 + random()*0.3,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,

0.6 + random()*0.3,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,

0.6 + random()*0.3,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,

0.6 + random()*0.3,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,

0.6 + random()*0.3,
0.5 + random()*0.3,
0.6 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2,
0.7 + random()*0.2,
0.6 + random()*0.2,
0.7 + random()*0.2
)::vector
FROM users;
