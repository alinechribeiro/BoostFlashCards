-- Stores AI-marked student answers and predicted grades

USE boostflashcards;

CREATE TABLE IF NOT EXISTS answer_attempts (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  subject_id BIGINT NOT NULL,
  topic_id BIGINT NULL,
  question TEXT NOT NULL,
  student_answer TEXT NOT NULL,
  predicted_score DECIMAL(5,2) NOT NULL,
  max_score DECIMAL(5,2) NOT NULL,
  predicted_grade VARCHAR(16) NOT NULL,
  feedback TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
  FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE SET NULL
);

