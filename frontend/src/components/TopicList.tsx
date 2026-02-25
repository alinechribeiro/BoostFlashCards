import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, Topic } from '../api/client'

export default function TopicList() {
  const { subjectId } = useParams<{ subjectId: string }>()
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [aiLoading, setAiLoading] = useState(false)
  const [aiMessage, setAiMessage] = useState<string | null>(null)

  useEffect(() => {
    if (!subjectId) return
    api
      .getTopics(Number(subjectId))
      .then(setTopics)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [subjectId])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading topics…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  return (
    <section>
      <Link to="/" style={{ marginBottom: 16, display: 'inline-block' }}>← Subjects</Link>
      <h1 style={{ marginBottom: 24 }}>Choose a topic</h1>
      <ul style={{ listStyle: 'none', padding: 0, display: 'flex', flexDirection: 'column', gap: 12 }}>
        {topics.map((t) => (
          <li key={t.id}>
            <Link
              to={`/topics/${t.id}/cards`}
              style={{
                display: 'block',
                padding: 16,
                background: 'var(--card)',
                borderRadius: 12,
                color: 'inherit',
                textDecoration: 'none',
              }}
            >
              {t.name}
            </Link>
          </li>
        ))}
      </ul>
      {topics.length > 0 && subjectId && (
        <p style={{ marginTop: 24 }}>
          <Link to={`/topics/${topics[0].id}/cards/new`} style={{ color: 'var(--accent)' }}>
            + Add flashcards to this subject
          </Link>
        </p>
      )}
      {subjectId && (
        <section style={{ marginTop: 32 }}>
          <h2 style={{ marginBottom: 8, fontSize: '1rem' }}>AI helper</h2>
          <p style={{ marginBottom: 8, color: 'var(--muted)' }}>
            Let the AI generate a practice topic with flashcards for this subject.
          </p>
          {aiMessage && (
            <p style={{ marginBottom: 8, color: 'var(--muted)' }}>
              {aiMessage}
            </p>
          )}
          <button
            type="button"
            disabled={aiLoading}
            onClick={async () => {
              if (!subjectId) return
              setAiLoading(true)
              setAiMessage(null)
              try {
                const res = await api.generateAIFlashcardsForSubject(Number(subjectId))
                setAiMessage(res.message)
                // refresh topics so new AI topics appear
                const updated = await api.getTopics(Number(subjectId))
                setTopics(updated)
              } catch (e) {
                setAiMessage(e instanceof Error ? e.message : 'Failed to generate flashcards with AI')
              } finally {
                setAiLoading(false)
              }
            }}
            style={{
              padding: '10px 20px',
              background: 'var(--accent)',
              border: 'none',
              borderRadius: 8,
              color: 'white',
              fontWeight: 600,
            }}
          >
            {aiLoading ? 'Asking AI…' : 'Generate AI flashcards'}
          </button>
        </section>
      )}
    </section>
  )
}
