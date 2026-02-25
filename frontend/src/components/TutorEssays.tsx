import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

interface TutorEssayRow {
  id: number
  student_name: string
  subject: string
  tier: string
  status: string
  created_at: string
  viewed_at?: string | null
}

export default function TutorEssays() {
  const [rows, setRows] = useState<TutorEssayRow[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .getTutorEssays()
      .then(setRows)
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load tutor essays'))
      .finally(() => setLoading(false))
  }, [])

  const awaiting = useMemo(() => rows.filter((r) => r.status === 'awaiting_tutor'), [rows])
  const inReview = useMemo(() => rows.filter((r) => r.status === 'in_review'), [rows])
  const history = useMemo(
    () => rows.filter((r) => r.status === 'submitted' || r.status === 'viewed_by_student' || r.status === 'completed'),
    [rows],
  )

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading essay dashboard…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  const renderSection = (title: string, data: TutorEssayRow[]) => (
    <section style={{ marginBottom: 24 }}>
      <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>{title}</h2>
      {data.length === 0 ? (
        <p style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>Nothing here right now.</p>
      ) : (
        <div
          style={{
            borderRadius: 12,
            border: '1px solid #222',
            overflow: 'hidden',
          }}
        >
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
            <thead style={{ background: '#090909' }}>
              <tr>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Student</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Subject</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Tier</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Created</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Status</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Open</th>
              </tr>
            </thead>
            <tbody>
              {data.map((r) => (
                <tr key={r.id} style={{ borderTop: '1px solid #222' }}>
                  <td style={{ padding: '8px 12px' }}>{r.student_name}</td>
                  <td style={{ padding: '8px 12px' }}>{r.subject || '—'}</td>
                  <td style={{ padding: '8px 12px' }}>{r.tier}</td>
                  <td style={{ padding: '8px 12px' }}>{new Date(r.created_at).toLocaleDateString()}</td>
                  <td style={{ padding: '8px 12px' }}>
                    <span
                      style={{
                        padding: '2px 8px',
                        borderRadius: 999,
                        border: '1px solid #333',
                        fontSize: '0.75rem',
                      }}
                    >
                      {r.status === 'viewed_by_student' ? 'Student opened review' : r.status}
                    </span>
                  </td>
                  <td style={{ padding: '8px 12px' }}>
                    <Link
                      to={`/essays/${r.id}`}
                      style={{
                        padding: '4px 10px',
                        borderRadius: 999,
                        border: '1px solid #333',
                        background: 'var(--card)',
                        color: 'inherit',
                        fontSize: '0.8rem',
                        textDecoration: 'none',
                      }}
                    >
                      Open
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )

  return (
    <div>
      <h1 style={{ marginBottom: 8 }}>Tutor essay dashboard</h1>
      <p style={{ marginBottom: 16, color: 'var(--muted)', fontSize: '0.9rem' }}>
        See new essay review requests from students, what you are currently working on, and your past reviews.
      </p>
      {renderSection('New requests', awaiting)}
      {renderSection('In progress', inReview)}
      {renderSection('History', history)}
    </div>
  )
}

