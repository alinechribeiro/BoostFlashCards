import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { api } from '../api/client'

export default function CreateFlashcard() {
  const { topicId } = useParams<{ topicId: string }>()
  const navigate = useNavigate()
  const [front, setFront] = useState('')
  const [back, setBack] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!topicId || !front.trim() || !back.trim()) return
    setSaving(true)
    setError(null)
    try {
      await api.createFlashcard(Number(topicId), front.trim(), back.trim())
      navigate(`/topics/${topicId}/cards`)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to save')
    } finally {
      setSaving(false)
    }
  }

  return (
    <section>
      <Link to={`/topics/${topicId}/cards`} style={{ marginBottom: 16, display: 'inline-block' }}>← Back to cards</Link>
      <h1 style={{ marginBottom: 24 }}>New flashcard</h1>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16, maxWidth: 400 }}>
        <label>
          <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Question (front)</span>
          <textarea
            value={front}
            onChange={(e) => setFront(e.target.value)}
            rows={3}
            required
            style={{
              width: '100%',
              padding: 12,
              background: 'var(--card)',
              border: '1px solid #333',
              borderRadius: 8,
              color: 'inherit',
            }}
          />
        </label>
        <label>
          <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Answer (back)</span>
          <textarea
            value={back}
            onChange={(e) => setBack(e.target.value)}
            rows={3}
            required
            style={{
              width: '100%',
              padding: 12,
              background: 'var(--card)',
              border: '1px solid #333',
              borderRadius: 8,
              color: 'inherit',
            }}
          />
        </label>
        {error && <p style={{ color: '#ef4444' }}>{error}</p>}
        <button
          type="submit"
          disabled={saving}
          style={{
            padding: 12,
            background: 'var(--accent)',
            border: 'none',
            borderRadius: 8,
            color: 'white',
            fontWeight: 600,
          }}
        >
          {saving ? 'Saving…' : 'Save card'}
        </button>
      </form>
    </section>
  )
}
