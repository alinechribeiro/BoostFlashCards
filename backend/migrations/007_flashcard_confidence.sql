ALTER TABLE flashcards
  ADD COLUMN status ENUM('not_yet', 'confident') NOT NULL DEFAULT 'not_yet',
  ADD COLUMN last_reviewed_at TIMESTAMP NULL,
  ADD COLUMN next_due_at TIMESTAMP NULL;

ALTER TABLE subjects
  ADD COLUMN last_reviewed_at TIMESTAMP NULL;

