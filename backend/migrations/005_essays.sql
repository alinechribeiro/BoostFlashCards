-- Essay review requests and reviews

CREATE TABLE IF NOT EXISTS essay_bundles (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  student_id BIGINT NOT NULL,
  tutor_id BIGINT NULL,
  total_essays INT NOT NULL,
  used_essays INT NOT NULL DEFAULT 0,
  price_cents INT NOT NULL,
  status ENUM('pending_payment', 'active', 'exhausted', 'expired') NOT NULL DEFAULT 'pending_payment',
  stripe_payment_intent_id VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL,
  CONSTRAINT fk_essay_bundles_student FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_essay_bundles_tutor FOREIGN KEY (tutor_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS essay_requests (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  student_id BIGINT NOT NULL,
  tutor_id BIGINT NOT NULL,
  tier ENUM('1_quick', '2_standard', '3_premium') NOT NULL,
  bundle_id BIGINT NULL,
  status ENUM(
    'pending_payment',
    'awaiting_tutor',
    'in_review',
    'submitted',
    'viewed_by_student',
    'completed'
  ) NOT NULL DEFAULT 'pending_payment',
  question_prompt TEXT NOT NULL,
  student_answer TEXT,
  answer_file_url VARCHAR(512) NOT NULL DEFAULT '',
  stripe_payment_intent_id VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_essay_requests_student FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_essay_requests_tutor FOREIGN KEY (tutor_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_essay_requests_bundle FOREIGN KEY (bundle_id) REFERENCES essay_bundles(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS essay_reviews (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  essay_request_id BIGINT NOT NULL,
  grade VARCHAR(32) NOT NULL DEFAULT '',
  quick_comments TEXT,
  mark_scheme_ref TEXT,
  strengths TEXT,
  improvements TEXT,
  improved_paragraph TEXT,
  audio_video_url VARCHAR(512) NOT NULL DEFAULT '',
  improvement_plan_url VARCHAR(512) NOT NULL DEFAULT '',
  submitted_at TIMESTAMP NULL,
  viewed_at TIMESTAMP NULL,
  CONSTRAINT fk_essay_reviews_request FOREIGN KEY (essay_request_id) REFERENCES essay_requests(id) ON DELETE CASCADE
);

