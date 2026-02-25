import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { api, EssayDetail } from '../api/client'
import { useAuth } from '../auth/AuthContext'

export default function EssayDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()
  const [detail, setDetail] = useState<EssayDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [grade, setGrade] = useState('')
  const [quickComments, setQuickComments] = useState('')
  const [markSchemeRef, setMarkSchemeRef] = useState('')
  const [strengths, setStrengths] = useState('')
  const [improvements, setImprovements] = useState('')
  const [improvedParagraph, setImprovedParagraph] = useState('')
  const [audioVideoURL, setAudioVideoURL] = useState('')
  const [improvementPlanURL, setImprovementPlanURL] = useState('')
  const [saving, setSaving] = useState(false)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [saveSuccess, setSaveSuccess] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    setError(null)
    api
      .getEssayDetail(Number(id))
      .then(async (d) => {
        setDetail(d)
        if (d.review) {
          setGrade(d.review.grade || '')
          setQuickComments(d.review.quick_comments || '')
          setMarkSchemeRef(d.review.mark_scheme_ref || '')
          setStrengths(d.review.strengths || '')
          setImprovements(d.review.improvements || '')
          setImprovedParagraph(d.review.improved_paragraph || '')
          setAudioVideoURL(d.review.audio_video_url || '')
          setImprovementPlanURL(d.review.improvement_plan_url || '')
        }
        // If student opens a reviewed essay, mark as viewed.
        if (user?.role === 'student' && d.review) {
          try {
            await api.markEssayViewed(Number(id))
          } catch {
            // ignore
          }
        }
      })
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load essay'))
      .finally(() => setLoading(false))
  }, [id, user?.role])

  if (!id) {
    return <p style={{ color: '#ef4444' }}>Missing essay id.</p>
  }

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading essay…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>
  if (!detail) return <p style={{ color: 'var(--muted)' }}>Essay not found.</p>

  const { request, review } = detail
  const isTutor = user?.role === 'tutor'

  return (
    <section>
      <button
        type="button"
        onClick={() => navigate(-1)}
        style={{ marginBottom: 12, padding: '4px 10px', borderRadius: 999, border: '1px solid #333', background: 'var(--card)', color: 'inherit', fontSize: '0.8rem' }}
      >
        ← Back
      </button>
      <h1 style={{ marginBottom: 4 }}>Essay request #{request.id}</h1>
      <p style={{ marginBottom: 12, color: 'var(--muted)', fontSize: '0.9rem' }}>
        Subject: <strong>{request.subject || '—'}</strong> • Tier: <strong>{request.tier}</strong> • Status:{' '}
        <strong>{request.status}</strong>
      </p>

      <div
        style={{
          marginBottom: 16,
          padding: 12,
          borderRadius: 12,
          border: '1px solid #222',
          background: 'var(--card)',
        }}
      >
        <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Question</h2>
        <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>{request.question_prompt}</p>
      </div>

      <div
        style={{
          marginBottom: 16,
          padding: 12,
          borderRadius: 12,
          border: '1px solid #222',
          background: 'var(--card)',
        }}
      >
        <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Student answer</h2>
        {request.answer_file_url && (
          <p style={{ marginBottom: 8, fontSize: '0.85rem' }}>
            Uploaded file:{' '}
            <a href={request.answer_file_url} target="_blank" rel="noreferrer">
              open
            </a>
          </p>
        )}
        <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
          {request.student_answer || <span style={{ color: 'var(--muted)' }}>No typed answer provided.</span>}
        </p>
      </div>

      {isTutor ? (
        <div
          style={{
            padding: 12,
            borderRadius: 12,
            border: '1px solid #333',
            background: '#050505',
          }}
        >
          <h2 style={{ fontSize: '1rem', marginBottom: 4 }}>Submit / update your review</h2>
          <p style={{ marginBottom: 12, color: 'var(--muted)', fontSize: '0.85rem' }}>
            Tier 1 focuses on grade + 2 short comments. Tier 2 adds mark scheme, strengths, improvements and an improved
            paragraph. Tier 3 also includes audio/video explanation and a personalised improvement plan PDF.
          </p>
          <form
            onSubmit={async (e) => {
              e.preventDefault()
              setSaving(true)
              setSaveError(null)
              setSaveSuccess(null)
              try {
                await api.submitEssayReview(Number(id), {
                  grade,
                  quick_comments: quickComments,
                  mark_scheme_ref: markSchemeRef,
                  strengths,
                  improvements,
                  improved_paragraph: improvedParagraph,
                  audio_video_url: audioVideoURL,
                  improvement_plan_url: improvementPlanURL,
                })
                setSaveSuccess('Review saved.')
              } catch (err) {
                setSaveError(err instanceof Error ? err.message : 'Failed to save review')
              } finally {
                setSaving(false)
              }
            }}
            style={{ display: 'flex', flexDirection: 'column', gap: 8 }}
          >
            <label style={{ fontSize: '0.85rem' }}>
              Grade
              <input
                value={grade}
                onChange={(e) => setGrade(e.target.value)}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Quick comments (Tier 1)
              <textarea
                value={quickComments}
                onChange={(e) => setQuickComments(e.target.value)}
                rows={2}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Mark scheme reference (Tier 2)
              <textarea
                value={markSchemeRef}
                onChange={(e) => setMarkSchemeRef(e.target.value)}
                rows={2}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Strengths (Tier 2)
              <textarea
                value={strengths}
                onChange={(e) => setStrengths(e.target.value)}
                rows={2}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Improvements (Tier 2)
              <textarea
                value={improvements}
                onChange={(e) => setImprovements(e.target.value)}
                rows={3}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Improved paragraph example (Tier 2)
              <textarea
                value={improvedParagraph}
                onChange={(e) => setImprovedParagraph(e.target.value)}
                rows={3}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Audio / video URL (Tier 3)
              <input
                value={audioVideoURL}
                onChange={(e) => setAudioVideoURL(e.target.value)}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>
            <label style={{ fontSize: '0.85rem' }}>
              Improvement plan PDF URL (Tier 3)
              <input
                value={improvementPlanURL}
                onChange={(e) => setImprovementPlanURL(e.target.value)}
                style={{
                  width: '100%',
                  marginTop: 4,
                  borderRadius: 8,
                  border: '1px solid #333',
                  background: '#000',
                  color: 'inherit',
                  padding: 6,
                  fontSize: '0.9rem',
                }}
              />
            </label>

            {saveError && (
              <p style={{ color: '#ef4444', fontSize: '0.85rem', margin: 0 }}>{saveError}</p>
            )}
            {saveSuccess && (
              <p style={{ color: '#22c55e', fontSize: '0.85rem', margin: 0 }}>{saveSuccess}</p>
            )}

            <button
              type="submit"
              disabled={saving}
              style={{
                marginTop: 8,
                alignSelf: 'flex-start',
                padding: '8px 14px',
                borderRadius: 8,
                border: 'none',
                background: saving ? '#555' : 'var(--accent)',
                color: '#fff',
                fontWeight: 600,
                fontSize: '0.9rem',
              }}
            >
              {saving ? 'Saving…' : 'Save review'}
            </button>
          </form>
        </div>
      ) : review ? (
        <div
          style={{
            padding: 12,
            borderRadius: 12,
            border: '1px solid #333',
            background: '#050505',
          }}
        >
          <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Tutor feedback</h2>
          <p style={{ marginBottom: 8, fontSize: '0.9rem' }}>
            <strong>Grade:</strong> {review.grade || '—'}
          </p>
          {review.quick_comments && (
            <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
              <strong>Quick comments:</strong> {review.quick_comments}
            </p>
          )}
          {review.mark_scheme_ref && (
            <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
              <strong>Mark scheme reference:</strong> {review.mark_scheme_ref}
            </p>
          )}
          {review.strengths && (
            <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
              <strong>Strengths:</strong> {review.strengths}
            </p>
          )}
          {review.improvements && (
            <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
              <strong>Improvements:</strong> {review.improvements}
            </p>
          )}
          {review.improved_paragraph && (
            <p style={{ whiteSpace: 'pre-line', fontSize: '0.9rem' }}>
              <strong>Improved paragraph example:</strong> {review.improved_paragraph}
            </p>
          )}
          {review.audio_video_url && (
            <p style={{ fontSize: '0.9rem' }}>
              <strong>Audio / video:</strong>{' '}
              <a href={review.audio_video_url} target="_blank" rel="noreferrer">
                open
              </a>
            </p>
          )}
          {review.improvement_plan_url && (
            <p style={{ fontSize: '0.9rem' }}>
              <strong>Improvement plan PDF:</strong>{' '}
              <a href={review.improvement_plan_url} target="_blank" rel="noreferrer">
                open
              </a>
            </p>
          )}
        </div>
      ) : (
        <p style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
          Your tutor is still preparing feedback for this essay.
        </p>
      )}
    </section>
  )
}

