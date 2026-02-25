const API = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(API + path, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    credentials: 'include',
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export interface Subject {
  id: number
  name: string
  slug: string
  created_at: string
  last_reviewed_at?: string | null
}

export interface Topic {
  id: number
  subject_id: number
  name: string
  slug: string
  created_at: string
}

export interface Flashcard {
  id: number
  topic_id: number
  front: string
  back: string
   status: 'not_yet' | 'confident'
  created_at: string
  updated_at: string
  last_reviewed_at?: string | null
  next_due_at?: string | null
}

export interface AIGenerationResult {
  subject_id: number
  topic_ids: number[]
  flashcards_created: number
  message: string
}

export interface AICreateSubjectResult {
  subject_id: number
  topic_ids: number[]
  flashcards_created: number
  message: string
}

export interface PracticeQuestion {
  subject_id: number
  question: string
}

export interface PracticeGradeResult {
  attempt_id: number
  subject_id: number
  topic_id?: number
  question: string
  student_answer: string
  score: number
  max_score: number
  grade: string
  grade_band?: string
  feedback: string
  strengths?: string
  improvements?: string
  score_percentage: number
}

export interface SubjectProgressAttempt {
  id: number
  created_at: string
  score_percentage: number
  grade: string
}

export interface SubjectProgress {
  subject_id: number
  subject_name: string
  attempts: SubjectProgressAttempt[]
  latest_grade: string
  average_score: number
  attempts_count: number
  encouragement_blurb: string
}

export interface Insight {
  text: string
}

export interface ExtractInsightsResult {
  insights: Insight[]
}

export interface User {
  id: number
  email: string
  name: string
  role: 'student' | 'tutor' | 'admin'
  avatar_url: string
  created_at: string
}

export interface TutorSummary {
  id: number
  name: string
  avatar_url: string
  headline: string
  bio: string
  hourly_rate_cents: number
  subjects: string[]
}

export interface TutorDetailSubject {
  name: string
  level: string
}

export interface TutorDetail {
  id: number
  name: string
  avatar_url: string
  headline: string
  bio: string
  hourly_rate_cents: number
  subjects: TutorDetailSubject[]
}

export type EssayTier = '1_quick' | '2_standard' | '3_premium'

export interface EssayBundle {
  id: number
  student_id: number
  tutor_id: number | null
  total_essays: number
  used_essays: number
  price_cents: number
  status: 'pending_payment' | 'active' | 'exhausted' | 'expired'
  stripe_payment_intent_id: string
  created_at: string
  expires_at?: string | null
}

export interface EssayRequest {
  id: number
  student_id: number
  tutor_id: number
  tier: EssayTier
  subject: string
  bundle_id?: number | null
  status:
    | 'pending_payment'
    | 'awaiting_tutor'
    | 'in_review'
    | 'submitted'
    | 'viewed_by_student'
    | 'completed'
  question_prompt: string
  student_answer: string
  answer_file_url: string
  stripe_payment_intent_id: string
  created_at: string
  updated_at: string
}

export interface EssayReview {
  id: number
  essay_request_id: number
  grade: string
  quick_comments: string
  mark_scheme_ref: string
  strengths: string
  improvements: string
  improved_paragraph: string
  audio_video_url: string
  improvement_plan_url: string
  submitted_at?: string | null
  viewed_at?: string | null
}

export interface EssayDetail {
  request: EssayRequest
  review?: EssayReview | null
}

export const api = {
  getSubjects: () => request<Subject[]>('/subjects'),
  getSubject: (id: number) => request<Subject>(`/subjects/${id}`),
  getTopics: (subjectId: number) => request<Topic[]>(`/subjects/${subjectId}/topics`),
  getTopic: (id: number) => request<Topic>(`/topics/${id}`),
  getFlashcards: (topicId: number) => request<Flashcard[]>(`/topics/${topicId}/flashcards`),
  getFlashcard: (id: number) => request<Flashcard>(`/flashcards/${id}`),
  createFlashcard: (topicId: number, front: string, back: string) =>
    request<Flashcard>('/flashcards', {
      method: 'POST',
      body: JSON.stringify({ topic_id: topicId, front, back }),
    }),
  updateFlashcard: (id: number, front?: string, back?: string) =>
    request<Flashcard>(`/flashcards/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ front, back }),
    }),
  deleteFlashcard: (id: number) =>
    request<void>(`/flashcards/${id}`, { method: 'DELETE' }),
  setFlashcardStatus: (id: number, status: 'not_yet' | 'confident') =>
    request<Flashcard>(`/flashcards/${id}/status`, {
      method: 'POST',
      body: JSON.stringify({ status }),
    }),
  generateAIFlashcardsForSubject: (subjectId: number, numCards?: number) =>
    request<AIGenerationResult>(`/subjects/${subjectId}/ai/flashcards`, {
      method: 'POST',
      body: JSON.stringify({ num_cards: numCards ?? 10 }),
    }),
  createSubjectWithAI: (prompt: string) =>
    request<AICreateSubjectResult>('/ai/subjects', {
      method: 'POST',
      body: JSON.stringify({ prompt }),
    }),
  getPracticeQuestion: (subjectId: number) =>
    request<PracticeQuestion>(`/subjects/${subjectId}/practice/question`, {
      method: 'POST',
    }),
  gradePracticeAnswer: (subjectId: number, question: string, answer: string, topicId?: number) =>
    request<PracticeGradeResult>(`/subjects/${subjectId}/practice/answer`, {
      method: 'POST',
      body: JSON.stringify({ question, answer, topic_id: topicId }),
    }),
  getSubjectProgress: (subjectId: number) =>
    request<SubjectProgress>(`/subjects/${subjectId}/progress`),
  extractInsights: (topicId: number, text: string, maxInsights?: number) =>
    request<ExtractInsightsResult>(`/topics/${topicId}/ai/insights`, {
      method: 'POST',
      body: JSON.stringify({ text, max_insights: maxInsights ?? 12 }),
    }),
  createFlashcardsFromInsights: (
    topicId: number,
    insights: { text: string; style: 'qa' | 'true_false' }[],
  ) =>
    request<{ created: number }>(`/topics/${topicId}/ai/flashcards-from-insights`, {
      method: 'POST',
      body: JSON.stringify({ insights }),
    }),
  signup: (data: { email: string; password: string; name: string; role: 'student' | 'tutor' }) =>
    request<User>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  login: (data: { email: string; password: string }) =>
    request<User>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  logout: () =>
    request<{ ok: string }>('/auth/logout', {
      method: 'POST',
    }),
  me: () => request<User>('/auth/me'),
  completeSocialSignup: (pendingToken: string, role: 'student' | 'tutor') =>
    request<User>('/auth/complete-signup', {
      method: 'POST',
      body: JSON.stringify({ pending_token: pendingToken, role }),
    }),
  getTutors: () => request<TutorSummary[]>('/tutors'),
  getTutor: (id: number) => request<TutorDetail>(`/tutors/${id}`),
  createBundle: (data: { tutor_id: number; total_essays: number }) =>
    request<{ bundle: EssayBundle; client_secret: string }>('/billing/bundles', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  createEssayRequest: (data: {
    tutor_id: number
    tier: EssayTier
    bundle_id?: number
    subject: string
    question_prompt: string
    student_answer: string
    answer_file_url?: string
  }) =>
    request<EssayRequest>('/essays/requests', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  markEssayViewed: (id: number) =>
    request<{ status: string }>(`/essays/${id}/mark_viewed`, {
      method: 'POST',
    }),
  getStudentEssays: () => request<any>('/students/me/essays'),
  getTutorEssays: () => request<any>('/tutors/me/essays'),
  getEssayDetail: (id: number) => request<EssayDetail>(`/essays/${id}`),
  submitEssayReview: (
    id: number,
    body: {
      grade: string
      quick_comments: string
      mark_scheme_ref: string
      strengths: string
      improvements: string
      improved_paragraph: string
      audio_video_url: string
      improvement_plan_url: string
    },
  ) =>
    request<EssayReview>(`/essays/${id}/review`, {
      method: 'POST',
      body: JSON.stringify(body),
    }),
}
