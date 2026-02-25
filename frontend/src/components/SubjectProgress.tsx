import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, SubjectProgress } from '../api/client'

export default function SubjectProgressPage() {
  const { subjectId } = useParams<{ subjectId: string }>()
  const [data, setData] = useState<SubjectProgress | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!subjectId) return
    setLoading(true)
    setError(null)
    api
      .getSubjectProgress(Number(subjectId))
      .then(setData)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [subjectId])

  if (!subjectId) return <p style={{ color: '#ef4444' }}>Missing subject id.</p>
  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading progress…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>
  if (!data) return null

  const maxHeight = 120
  const attempts = data.attempts
  const highest = attempts.reduce((max, a) => (a.score_percentage > max ? a.score_percentage : max), 0)
  const scale = highest > 0 ? highest : 100

  return (
    <section>
      <Link to="/" style={{ marginBottom: 16, display: 'inline-block' }}>
        ← Subjects
      </Link>
      <h1 style={{ marginBottom: 8 }}>{data.subject_name} progress</h1>
      <p style={{ marginBottom: 12, color: 'var(--muted)' }}>
        AI-marked answers over time with predicted GCSE grade.
      </p>

      <div style={{ marginBottom: 16 }}>
        <p style={{ marginBottom: 4 }}>
          <strong>Current predicted grade:</strong> {data.latest_grade || '—'}
        </p>
        <p style={{ marginBottom: 4 }}>
          <strong>Average score:</strong> {data.average_score.toFixed(1)}%
        </p>
        <p style={{ marginBottom: 4 }}>
          <strong>Attempts:</strong> {data.attempts_count}
        </p>
      </div>

      {attempts.length === 0 ? (
        <p style={{ color: 'var(--muted)', marginBottom: 16 }}>
          No AI-marked answers yet. Do an AI-marked practice session to start your progress graph.
        </p>
      ) : (
        <div style={{ marginBottom: 16 }}>
          <div
            style={{
              display: 'flex',
              alignItems: 'flex-end',
              gap: 8,
              padding: 12,
              borderRadius: 12,
              background: 'var(--card)',
              border: '1px solid #333',
              minHeight: maxHeight + 32,
            }}
          >
            {attempts.map((a, idx) => {
              const height = (a.score_percentage / scale) * maxHeight
              return (
                <div
                  key={a.id}
                  title={`Attempt ${idx + 1}: ${a.score_percentage.toFixed(1)}% (grade ${a.grade})`}
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    flex: 1,
                    minWidth: 10,
                  }}
                >
                  <div
                    style={{
                      width: '60%',
                      height: Math.max(4, height),
                      borderRadius: 999,
                      background:
                        a.score_percentage >= 70
                          ? 'linear-gradient(to top, #22c55e, #4ade80)'
                          : a.score_percentage >= 40
                          ? 'linear-gradient(to top, #f59e0b, #fbbf24)'
                          : 'linear-gradient(to top, #ef4444, #f97373)',
                    }}
                  />
                  <span style={{ marginTop: 6, fontSize: '0.7rem', color: 'var(--muted)' }}>{idx + 1}</span>
                </div>
              )
            })}
          </div>
          <p style={{ marginTop: 8, fontSize: '0.8rem', color: 'var(--muted)' }}>
            Each bar is one AI-marked answer (left = earliest, right = latest).
          </p>
        </div>
      )}

      <p style={{ marginTop: 8, color: 'var(--accent)' }}>{data.encouragement_blurb}</p>

      <p style={{ marginTop: 24 }}>
        <Link to={`/subjects/${subjectId}/practice`} style={{ color: 'var(--accent)' }}>
          Do another AI-marked question
        </Link>
      </p>
    </section>
  )
}

