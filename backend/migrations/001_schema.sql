-- BoostFlashCards GCSE UK - MySQL schema

CREATE DATABASE IF NOT EXISTS boostflashcards;
USE boostflashcards;

CREATE TABLE IF NOT EXISTS subjects (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS topics (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  subject_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
  UNIQUE KEY (subject_id, slug)
);

CREATE TABLE IF NOT EXISTS flashcards (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  topic_id BIGINT NOT NULL,
  front TEXT NOT NULL,
  back TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
);

-- Seed GCSE subjects and sample topics (UK GCSE)
INSERT INTO subjects (name, slug) VALUES
  ('Mathematics', 'mathematics'),
  ('Biology', 'biology'),
  ('Chemistry', 'chemistry'),
  ('Physics', 'physics'),
  ('English Literature', 'english-literature'),
  ('History', 'history')
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO topics (subject_id, name, slug) VALUES
  (1, 'Algebra', 'algebra'),
  (1, 'Number', 'number'),
  (1, 'Geometry', 'geometry'),
  (2, 'Cell Biology', 'cell-biology'),
  (2, 'Organisation', 'organisation'),
  (3, 'Atomic Structure', 'atomic-structure'),
  (4, 'Energy', 'energy'),
  (5, 'Shakespeare', 'shakespeare'),
  (6, 'Norman England', 'norman-england')
ON DUPLICATE KEY UPDATE name = VALUES(name);
