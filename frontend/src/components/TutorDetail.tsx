import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { api, TutorDetail } from '../api/client'

export default function TutorDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [tutor, setTutor] = useState<TutorDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    setError(null)
    api
      .getTutor(Number(id))
      .then(setTutor)
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load tutor'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading tutor…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>
  if (!tutor) return <p style={{ color: 'var(--muted)' }}>Tutor not found.</p>

  return (
    <section>
      <Link to="/tutors" style={{ marginBottom: 16, display: 'inline-block' }}>
        ← Back to tutors
      </Link>
      <div
        style={{
          display: 'flex',
          gap: 16,
          alignItems: 'flex-start',
          marginBottom: 16,
        }}
      >
        <div
          style={{
            width: 72,
            height: 72,
            borderRadius: '50%',
            background: '#111',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '2rem',
            overflow: 'hidden',
            flexShrink: 0,
          }}
        >
          {tutor.avatar_url ? (
            // eslint-disable-next-line jsx-a11y/alt-text
            <img src={tutor.avatar_url} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
          ) : (
            <span>{tutor.name.charAt(0).toUpperCase()}</span>
          )}
        </div>
        <div>
          <h1 style={{ margin: 0 }}>{tutor.name}</h1>
          {tutor.headline && (
            <p style={{ marginTop: 4, marginBottom: 4, color: 'var(--muted)' }}>{tutor.headline}</p>
          )}
          {tutor.hourly_rate_cents > 0 && (
            <p style={{ marginTop: 4, marginBottom: 4, color: 'var(--muted)', fontSize: '0.9rem' }}>
              Approx. £{(tutor.hourly_rate_cents / 100).toFixed(2)} / hour
            </p>
          )}
        </div>
      </div>

      {tutor.subjects.length > 0 && (
        <div style={{ marginBottom: 16 }}>
          <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Subjects</h2>
          <ul style={{ paddingLeft: 16, margin: 0 }}>
            {tutor.subjects.map((s, idx) => (
              <li key={idx} style={{ marginBottom: 4, fontSize: '0.95rem' }}>
                {s.name}
                {s.level && ` (${s.level})`}
              </li>
            ))}
          </ul>
        </div>
      )}

      {tutor.bio && (
        <div style={{ marginBottom: 16 }}>
          <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>About</h2>
          <p style={{ whiteSpace: 'pre-line' }}>{tutor.bio}</p>
        </div>
      )}

      <button
        type="button"
        style={{
          padding: '10px 18px',
          borderRadius: 8,
          border: 'none',
          background: 'var(--accent)',
          color: '#fff',
          fontWeight: 600,
        }}
      >
        Request essay review
      </button>
    </section>
  )
}

