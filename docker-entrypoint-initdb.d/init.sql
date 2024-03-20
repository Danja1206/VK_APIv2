CREATE DATABASE IF NOT EXISTS vk;
USE vk;

CREATE TABLE users (
 id int AUTO_INCREMENT,
 name varchar(255) NOT NULL,
 balance int NOT NULL,
 exp int unsigned NOT NULL DEFAULT '0',
 lvl int unsigned NOT NULL DEFAULT '1',
 PRIMARY KEY (id)
);

CREATE TABLE quests (
 id int NOT NULL AUTO_INCREMENT,
 name varchar(255) NOT NULL,
 cost int NOT NULL,
 is_unique tinyint(1) DEFAULT '0',
 max_limit int unsigned NOT NULL DEFAULT '1',
 exp int NOT NULL,
 PRIMARY KEY (id)
);

CREATE TABLE daily_quests (
 id int NOT NULL AUTO_INCREMENT,
 name varchar(255) NOT NULL,
 description text,
 cost int NOT NULL DEFAULT '0',
 issued_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 exp int NOT NULL,
 PRIMARY KEY (id)
);

CREATE TABLE big_quests (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    cost int DEFAULT NULL,
    exp int NOT NULL
);

CREATE TABLE big_quest_steps (
    id INT PRIMARY KEY AUTO_INCREMENT,
    big_quest_id INT,
    quest_id INT, 
    step_order INT, 
    FOREIGN KEY (big_quest_id) REFERENCES big_quests(id),
    FOREIGN KEY (quest_id) REFERENCES quests(id)
);

CREATE TABLE levels (
 id int unsigned NOT NULL AUTO_INCREMENT,
 lvl int unsigned NOT NULL,
 exp_total int NOT NULL,
 exp_to_lvl int NOT NULL,
 PRIMARY KEY (id)
);

CREATE TABLE user_big_quests (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    big_quest_id INT,
     completed_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (big_quest_id) REFERENCES big_quests(id)
); 

CREATE TABLE user_big_quests_steps (
    id INT PRIMARY KEY AUTO_INCREMENT, 
    user_id INT,
    big_quest_id INT,
    current_step INT, 
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (big_quest_id) REFERENCES big_quests(id)
); 

CREATE TABLE user_daily_quests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    daily_quest_id INT NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (daily_quest_id) REFERENCES daily_quests(id),
    UNIQUE KEY (user_id, daily_quest_id)
);

CREATE TABLE user_quests (
    id INT AUTO_INCREMENT,
    user_id INT,
    quest_id INT,
    completed_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quest_id) REFERENCES quests(id)
);

INSERT INTO levels (lvl, exp_total, exp_to_lvl)
VALUES
('1', '0', '40'),
('2', '40', '121'),
('3', '161', '231'),
('4', '392', '366'),
('5', '758', '523'),
('6', '1281', '700'),
('7', '1981', '896'),
('8', '2877', '1109'),
('9', '3986', '1339'),
('10', '5325', '1585'),
('11', '6910', '1846'),
('12', '8756', '2122'),
('13', '10878', '2412'),
('14', '13290', '2715'),
('15', '16005', '3032');