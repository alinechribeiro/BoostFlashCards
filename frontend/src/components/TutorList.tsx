import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, EssayTier, TutorSummary } from '../api/client'

export default function TutorList() {
  const [tutors, setTutors] = useState<TutorSummary[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTutorId, setActiveTutorId] = useState<number | null>(null)
  const [tier, setTier] = useState<EssayTier>('2_standard')
  const [bundleSize, setBundleSize] = useState<3 | 5>(3)
  const [question, setQuestion] = useState('')
  const [answer, setAnswer] = useState('')
  const [requestLoading, setRequestLoading] = useState(false)
  const [requestError, setRequestError] = useState<string | null>(null)
  const [requestSuccess, setRequestSuccess] = useState<string | null>(null)

  useEffect(() => {
    api
      .getTutors()
      .then(setTutors)
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load tutors'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading tutors…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  return (
    <section>
      <h1 style={{ marginBottom: 8 }}>Find a tutor</h1>
      <p style={{ marginBottom: 24, color: 'var(--muted)' }}>
        GCSE tutors available for essay feedback and subject‑specific help.
      </p>
      {tutors.length === 0 ? (
        <p style={{ color: 'var(--muted)' }}>No tutors are listed yet.</p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0, display: 'flex', flexDirection: 'column', gap: 12 }}>
          {tutors.map((t) => (
            <li key={t.id}>
              <article
                style={{
                  display: 'flex',
                  gap: 12,
                  padding: 16,
                  borderRadius: 12,
                  background: 'var(--card)',
                }}
              >
                <div
                  style={{
                    width: 56,
                    height: 56,
                    borderRadius: '50%',
                    background: '#111',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: '1.4rem',
                    flexShrink: 0,
                    overflow: 'hidden',
                  }}
                >
                  {t.avatar_url ? (
                    // eslint-disable-next-line jsx-a11y/alt-text
                    <img src={t.avatar_url} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                  ) : (
                    <span>{t.name.charAt(0).toUpperCase()}</span>
                  )}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <h2 style={{ margin: 0, fontSize: '1rem' }}>{t.name}</h2>
                  {t.headline && (
                    <p style={{ marginTop: 4, marginBottom: 4, color: 'var(--muted)', fontSize: '0.9rem' }}>
                      {t.headline}
                    </p>
                  )}
                  {t.subjects.length > 0 && (
                    <p style={{ marginTop: 4, marginBottom: 4, fontSize: '0.9rem' }}>
                      <strong>Subjects:</strong> {t.subjects.join(', ')}
                    </p>
                  )}
                  {t.bio && (
                    <p style={{ marginTop: 4, marginBottom: 4, fontSize: '0.9rem' }}>
                      {t.bio.length > 200 ? `${t.bio.slice(0, 200)}…` : t.bio}
                    </p>
                  )}
                  {t.hourly_rate_cents > 0 && (
                    <p style={{ marginTop: 4, marginBottom: 4, fontSize: '0.9rem', color: 'var(--muted)' }}>
                      Approx. £{(t.hourly_rate_cents / 100).toFixed(2)} / hour
                    </p>
                  )}
                  <div style={{ marginTop: 8, display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                    <Link
                      to={`/tutors/${t.id}`}
                      style={{
                        padding: '8px 14px',
                        borderRadius: 8,
                        background: 'var(--card)',
                        border: '1px solid #333',
                        color: 'inherit',
                        fontSize: '0.9rem',
                        textDecoration: 'none',
                      }}
                    >
                      View profile
                    </Link>
                    <button
                      type="button"
                      onClick={() => {
                        setActiveTutorId((prev) => (prev === t.id ? null : t.id))
                        setRequestError(null)
                        setRequestSuccess(null)
                      }}
                      style={{
                        padding: '8px 14px',
                        borderRadius: 8,
                        background: 'var(--accent)',
                        border: 'none',
                        color: '#fff',
                        fontSize: '0.9rem',
                        fontWeight: 600,
                      }}
                    >
                      Request essay review
                    </button>
                  </div>
                  {activeTutorId === t.id && (
                    <div
                      style={{
                        marginTop: 16,
                        padding: 12,
                        borderRadius: 10,
                        border: '1px solid #333',
                        background: 'rgba(15,15,15,0.8)',
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 8,
                      }}
                    >
                      <p style={{ margin: 0, fontSize: '0.9rem', color: 'var(--muted)' }}>
                        Choose a review tier and bundle, then paste your GCSE question and answer. Payment is handled
                        via a bundle so you can send several essays to this tutor.
                      </p>
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                        <label style={{ fontSize: '0.85rem' }}>
                          Tier:{' '}
                          <select
                            value={tier}
                            onChange={(e) => setTier(e.target.value as EssayTier)}
                            style={{ fontSize: '0.85rem' }}
                          >
                            <option value="1_quick">Tier 1 – Quick grade (£1 / essay)</option>
                            <option value="2_standard">Tier 2 – Standard review (£3.50 / essay)</option>
                            <option value="3_premium">Tier 3 – Premium (£10 / essay)</option>
                          </select>
                        </label>
                        <label style={{ fontSize: '0.85rem' }}>
                          Bundle size:{' '}
                          <select
                            value={bundleSize}
                            onChange={(e) => setBundleSize(Number(e.target.value) as 3 | 5)}
                            style={{ fontSize: '0.85rem' }}
                          >
                            <option value={3}>3 essays (10% discount)</option>
                            <option value={5}>5 essays (15% discount)</option>
                          </select>
                        </label>
                      </div>
                      <label style={{ fontSize: '0.85rem', display: 'block' }}>
                        Subject
                        <input
                          type="text"
                          value={question ? '' : ''} // placeholder; subject is optional in UI, but required in backend
                          placeholder="e.g. GCSE English Literature"
                          style={{
                            width: '100%',
                            marginTop: 4,
                            borderRadius: 8,
                            border: '1px solid #333',
                            background: '#050505',
                            color: 'inherit',
                            padding: 8,
                            fontSize: '0.9rem',
                          }}
                          readOnly
                        />
                      </label>
                      <label style={{ fontSize: '0.85rem', display: 'block' }}>
                        Question prompt
                        <textarea
                          value={question}
                          onChange={(e) => setQuestion(e.target.value)}
                          rows={3}
                          style={{
                            width: '100%',
                            marginTop: 4,
                            borderRadius: 8,
                            border: '1px solid #333',
                            background: '#050505',
                            color: 'inherit',
                            padding: 8,
                            fontSize: '0.9rem',
                          }}
                        />
                      </label>
                      <label style={{ fontSize: '0.85rem', display: 'block' }}>
                        Your answer (typed)
                        <textarea
                          value={answer}
                          onChange={(e) => setAnswer(e.target.value)}
                          rows={4}
                          style={{
                            width: '100%',
                            marginTop: 4,
                            borderRadius: 8,
                            border: '1px solid #333',
                            background: '#050505',
                            color: 'inherit',
                            padding: 8,
                            fontSize: '0.9rem',
                          }}
                        />
                      </label>
                      {requestError && (
                        <p style={{ color: '#ef4444', fontSize: '0.85rem', margin: 0 }}>{requestError}</p>
                      )}
                      {requestSuccess && (
                        <p style={{ color: '#22c55e', fontSize: '0.85rem', margin: 0 }}>{requestSuccess}</p>
                      )}
                      <button
                        type="button"
                        disabled={requestLoading || !question.trim()}
                        onClick={async () => {
                          try {
                            setRequestLoading(true)
                            setRequestError(null)
                            setRequestSuccess(null)
                            const { bundle } = await api.createBundle({
                              tutor_id: t.id,
                              total_essays: bundleSize,
                            })
                            await api.createEssayRequest({
                              tutor_id: t.id,
                              tier,
                              bundle_id: bundle.id,
                              subject: 'GCSE English', // simple default subject; can be expanded later
                              question_prompt: question,
                              student_answer: answer,
                            })
                            setRequestSuccess(
                              'Bundle created and essay request saved. Complete payment to send it to your tutor.',
                            )
                          } catch (e) {
                            setRequestError(e instanceof Error ? e.message : 'Failed to create request')
                          } finally {
                            setRequestLoading(false)
                          }
                        }}
                        style={{
                          marginTop: 4,
                          alignSelf: 'flex-start',
                          padding: '8px 14px',
                          borderRadius: 8,
                          background: requestLoading ? '#555' : 'var(--accent)',
                          border: 'none',
                          color: '#fff',
                          fontSize: '0.9rem',
                          fontWeight: 600,
                          cursor: requestLoading ? 'default' : 'pointer',
                        }}
                      >
                        {requestLoading ? 'Creating bundle…' : 'Buy bundle & send essay'}
                      </button>
                    </div>
                  )}
                </div>
              </article>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}

