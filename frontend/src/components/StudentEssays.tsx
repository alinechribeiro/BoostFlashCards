import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

interface StudentEssayRow {
  id: number
  tutor_name: string
  tier: string
  status: string
  subject: string
  created_at: string
  bundle_id?: number | null
  has_review: boolean
}

export default function StudentEssays() {
  const [rows, setRows] = useState<StudentEssayRow[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .getStudentEssays()
      .then(setRows)
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load essays'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading your essays…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  return (
    <section>
      <h1 style={{ marginBottom: 8 }}>My essay requests</h1>
      <p style={{ marginBottom: 16, color: 'var(--muted)', fontSize: '0.9rem' }}>
        Track all essays you have sent to tutors, their status, and when feedback is available.
      </p>
      {rows.length === 0 ? (
        <p style={{ color: 'var(--muted)' }}>You have not requested any essay reviews yet.</p>
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
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Tutor</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Subject</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Tier</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Status</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Bundle</th>
                <th style={{ textAlign: 'left', padding: '8px 12px' }}>Review</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <tr key={r.id} style={{ borderTop: '1px solid #222' }}>
                  <td style={{ padding: '8px 12px' }}>{r.tutor_name}</td>
                  <td style={{ padding: '8px 12px' }}>{r.subject || '—'}</td>
                  <td style={{ padding: '8px 12px' }}>{r.tier}</td>
                  <td style={{ padding: '8px 12px' }}>{r.status}</td>
                  <td style={{ padding: '8px 12px' }}>{r.bundle_id ?? '—'}</td>
                  <td style={{ padding: '8px 12px' }}>
                    {r.has_review ? (
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
                        Open review
                      </Link>
                    ) : (
                      <span style={{ color: 'var(--muted)' }}>Not confident yet</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}

