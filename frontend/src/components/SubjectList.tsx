import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, Subject } from '../api/client'

export default function SubjectList() {
  const [subjects, setSubjects] = useState<Subject[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.getSubjects()
      .then(setSubjects)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading subjects…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  return (
    <section>
      <h1 style={{ marginBottom: 24 }}>Choose a GCSE subject</h1>
      <ul style={{ listStyle: 'none', padding: 0, display: 'flex', flexDirection: 'column', gap: 12 }}>
        {subjects.map((s) => {
          const last = s.last_reviewed_at ? new Date(s.last_reviewed_at) : null
          const daysSince =
            last != null ? (Date.now() - last.getTime()) / (1000 * 60 * 60 * 24) : null
          const showWarning = daysSince != null && daysSince > 2
          return (
            <li key={s.id}>
              <div
                style={{
                  position: 'relative',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: 8,
                  padding: 16,
                  background: 'var(--card)',
                  borderRadius: 12,
                }}
              >
                {showWarning && (
                  <span
                    style={{
                      position: 'absolute',
                      top: 10,
                      right: 12,
                      padding: '4px 10px',
                      borderRadius: 999,
                      fontSize: '0.75rem',
                      border: '1px solid #7c2d12',
                      background: '#451a03',
                      color: '#fed7aa',
                    }}
                  >
                    Not reviewed for 2+ days — keep revising, you’re building strong GCSE memory!
                  </span>
                )}
                <Link
                  to={`/subjects/${s.id}/topics`}
                  style={{
                    color: 'inherit',
                    textDecoration: 'none',
                    fontWeight: 600,
                    marginBottom: 4,
                  }}
                >
                  {s.name}
                </Link>
                <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
                  <Link to={`/subjects/${s.id}/topics`} style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
                    Browse topics
                  </Link>
                  <Link to={`/subjects/${s.id}/practice`} style={{ color: 'var(--accent)', fontSize: '0.9rem' }}>
                    AI-marked practice
                  </Link>
                  <Link to={`/subjects/${s.id}/progress`} style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
                    Progress & predicted grade
                  </Link>
                </div>
              </div>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
