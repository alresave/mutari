
-- =========================
-- COMPANIES
-- =========================

INSERT INTO companies (id, name) VALUES
(1, 'NovaTech'),
(2, 'Atlas Systems'),
(3, 'Horizon Labs');

-- =========================
-- TEAMS
-- =========================

INSERT INTO teams (id, company_id, name) VALUES
(1, 1, 'NovaTech Engineering'),
(2, 1, 'NovaTech Product'),
(3, 2, 'Atlas Sales'),
(4, 2, 'Atlas Operations'),
(5, 3, 'Horizon AI'),
(6, 3, 'Horizon Design'),
(7, 3, 'Horizon Marketing');

-- =========================
-- USERS (50)
-- =========================

INSERT INTO users (id, name, photo_url, email) VALUES
(1,'Alice','https://randomuser.me/api/portraits/women/1.jpg','alice@novatech.com'),
(2,'Bob','https://randomuser.me/api/portraits/men/2.jpg','bob@novatech.com'),
(3,'Carlos','https://randomuser.me/api/portraits/men/3.jpg','carlos@novatech.com'),
(4,'Diana','https://randomuser.me/api/portraits/women/4.jpg','diana@novatech.com'),
(5,'Ethan','https://randomuser.me/api/portraits/men/5.jpg','ethan@novatech.com'),
(6,'Fernanda','https://randomuser.me/api/portraits/women/6.jpg','fernanda@novatech.com'),
(7,'Gabriel','https://randomuser.me/api/portraits/men/7.jpg','gabriel@novatech.com'),
(8,'Hannah','https://randomuser.me/api/portraits/women/8.jpg','hannah@novatech.com'),
(9,'Ivan','https://randomuser.me/api/portraits/men/9.jpg','ivan@novatech.com'),
(10,'Julia','https://randomuser.me/api/portraits/women/10.jpg','julia@novatech.com'),

(11,'Kevin','https://randomuser.me/api/portraits/men/11.jpg','kevin@atlas.com'),
(12,'Laura','https://randomuser.me/api/portraits/women/12.jpg','laura@atlas.com'),
(13,'Mateo','https://randomuser.me/api/portraits/men/13.jpg','mateo@atlas.com'),
(14,'Nina','https://randomuser.me/api/portraits/women/14.jpg','nina@atlas.com'),
(15,'Oscar','https://randomuser.me/api/portraits/men/15.jpg','oscar@atlas.com'),
(16,'Paula','https://randomuser.me/api/portraits/women/16.jpg','paula@atlas.com'),
(17,'Quentin','https://randomuser.me/api/portraits/men/17.jpg','quentin@atlas.com'),
(18,'Rosa','https://randomuser.me/api/portraits/women/18.jpg','rosa@atlas.com'),
(19,'Sergio','https://randomuser.me/api/portraits/men/19.jpg','sergio@atlas.com'),
(20,'Tania','https://randomuser.me/api/portraits/women/20.jpg','tania@atlas.com'),

(21,'Ulises','https://randomuser.me/api/portraits/men/21.jpg','ulises@horizon.com'),
(22,'Valeria','https://randomuser.me/api/portraits/women/22.jpg','valeria@horizon.com'),
(23,'William','https://randomuser.me/api/portraits/men/23.jpg','william@horizon.com'),
(24,'Ximena','https://randomuser.me/api/portraits/women/24.jpg','ximena@horizon.com'),
(25,'Yuri','https://randomuser.me/api/portraits/men/25.jpg','yuri@horizon.com'),
(26,'Zara','https://randomuser.me/api/portraits/women/26.jpg','zara@horizon.com'),
(27,'Aaron','https://randomuser.me/api/portraits/men/27.jpg','aaron@horizon.com'),
(28,'Bianca','https://randomuser.me/api/portraits/women/28.jpg','bianca@horizon.com'),
(29,'Cesar','https://randomuser.me/api/portraits/men/29.jpg','cesar@horizon.com'),
(30,'Daniela','https://randomuser.me/api/portraits/women/30.jpg','daniela@horizon.com'),

(31,'Eduardo','https://randomuser.me/api/portraits/men/31.jpg','eduardo@candidates.com'),
(32,'Fatima','https://randomuser.me/api/portraits/women/32.jpg','fatima@candidates.com'),
(33,'Gustavo','https://randomuser.me/api/portraits/men/33.jpg','gustavo@candidates.com'),
(34,'Helena','https://randomuser.me/api/portraits/women/34.jpg','helena@candidates.com'),
(35,'Ismael','https://randomuser.me/api/portraits/men/35.jpg','ismael@candidates.com'),
(36,'Julieta','https://randomuser.me/api/portraits/women/36.jpg','julieta@candidates.com'),
(37,'Karla','https://randomuser.me/api/portraits/women/37.jpg','karla@candidates.com'),
(38,'Luis','https://randomuser.me/api/portraits/men/38.jpg','luis@candidates.com'),
(39,'Monica','https://randomuser.me/api/portraits/women/39.jpg','monica@candidates.com'),
(40,'Nestor','https://randomuser.me/api/portraits/men/40.jpg','nestor@candidates.com'),
(41,'Olga','https://randomuser.me/api/portraits/women/41.jpg','olga@candidates.com'),
(42,'Pablo','https://randomuser.me/api/portraits/men/42.jpg','pablo@candidates.com'),
(43,'Raul','https://randomuser.me/api/portraits/men/43.jpg','raul@candidates.com'),
(44,'Sara','https://randomuser.me/api/portraits/women/44.jpg','sara@candidates.com'),
(45,'Tom','https://randomuser.me/api/portraits/men/45.jpg','tom@candidates.com'),
(46,'Ursula','https://randomuser.me/api/portraits/women/46.jpg','ursula@candidates.com'),
(47,'Victor','https://randomuser.me/api/portraits/men/47.jpg','victor@candidates.com'),
(48,'Wendy','https://randomuser.me/api/portraits/women/48.jpg','wendy@candidates.com'),
(49,'Xavier','https://randomuser.me/api/portraits/men/49.jpg','xavier@candidates.com'),
(50,'Yasmin','https://randomuser.me/api/portraits/women/50.jpg','yasmin@candidates.com');

-- =========================
-- TEAM MEMBERS
-- =========================

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

-- Users 36-50 remain available as candidate pool
