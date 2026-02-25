import { FormEvent, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { api, Insight } from '../api/client'

type InsightRow = {
  id: number
  text: string
  include: boolean
  style: 'qa' | 'true_false'
}

export default function TextToFlashcards() {
  const { topicId } = useParams<{ topicId: string }>()
  const navigate = useNavigate()

  const [rawText, setRawText] = useState('')
  const [insights, setInsights] = useState<InsightRow[]>([])
  const [extracting, setExtracting] = useState(false)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleExtract = async () => {
    if (!topicId) return
    if (!rawText.trim()) {
      setError('Paste some text first.')
      return
    }
    setExtracting(true)
    setError(null)
    try {
      const res = await api.extractInsights(Number(topicId), rawText)
      const rows: InsightRow[] = (res.insights || []).map((i: Insight, idx: number) => ({
        id: idx + 1,
        text: i.text,
        include: true,
        style: 'qa',
      }))
      setInsights(rows)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to extract insights with AI')
    } finally {
      setExtracting(false)
    }
  }

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault()
    if (!topicId) return
    const selected = insights.filter((i) => i.include && i.text.trim())
    if (selected.length === 0) {
      setError('Select at least one insight to turn into flashcards.')
      return
    }
    setSaving(true)
    setError(null)
    try {
      await api.createFlashcardsFromInsights(
        Number(topicId),
        selected.map((i) => ({ text: i.text.trim(), style: i.style })),
      )
      navigate(`/topics/${topicId}/cards`)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to create flashcards from insights')
    } finally {
      setSaving(false)
    }
  }

  if (!topicId) {
    return <p style={{ color: '#ef4444' }}>Missing topic id.</p>
  }

  const includedCount = insights.filter((i) => i.include && i.text.trim()).length

  return (
    <section>
      <Link to={`/topics/${topicId}/cards`} style={{ marginBottom: 16, display: 'inline-block' }}>
        ← Back to cards
      </Link>
      <h1 style={{ marginBottom: 8 }}>Text to flashcards</h1>
      <p style={{ marginBottom: 16, color: 'var(--muted)' }}>
        Paste notes, textbook extracts or teacher slides. The AI will pull out key insights; you choose which ones
        become flashcards and whether they are question/answer or true/false.
      </p>

      <div style={{ marginBottom: 16 }}>
        <textarea
          value={rawText}
          onChange={(e) => setRawText(e.target.value)}
          rows={8}
          placeholder="Paste your text here…"
          style={{
            width: '100%',
            padding: 12,
            background: 'var(--card)',
            border: '1px solid #333',
            borderRadius: 8,
            color: 'inherit',
          }}
        />
        <button
          type="button"
          onClick={handleExtract}
          disabled={extracting || !rawText.trim()}
          style={{
            marginTop: 8,
            padding: '10px 20px',
            background: 'var(--accent)',
            border: 'none',
            borderRadius: 8,
            color: 'white',
            fontWeight: 600,
          }}
        >
          {extracting ? 'Finding insights…' : 'Extract insights'}
        </button>
      </div>

      {error && <p style={{ color: '#ef4444', marginBottom: 12 }}>{error}</p>}

      {insights.length > 0 && (
        <form onSubmit={handleCreate} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              gap: 8,
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'column' }}>
              <h2 style={{ fontSize: '1rem' }}>Review insights</h2>
              <span style={{ fontSize: '0.85rem', color: 'var(--muted)' }}>
                Selected: {includedCount} flashcard{includedCount === 1 ? '' : 's'}
              </span>
            </div>
            <button
              type="submit"
              disabled={saving}
              style={{
                padding: '10px 20px',
                background: 'var(--accent)',
                border: 'none',
                borderRadius: 8,
                color: 'white',
                fontWeight: 600,
              }}
            >
              {saving
                ? includedCount > 0
                  ? `Creating ${includedCount} flashcard${includedCount === 1 ? '' : 's'}…`
                  : 'Creating flashcards…'
                : includedCount > 0
                  ? `Next: create ${includedCount} flashcard${includedCount === 1 ? '' : 's'}`
                  : 'Next: create flashcards'}
            </button>
          </div>
          <p style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
            Toggle which insights to include, edit the text, and choose whether each becomes a question/answer or a
            true/false flashcard.
          </p>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            {insights.map((insight) => (
              <div
                key={insight.id}
                style={{
                  padding: 12,
                  borderRadius: 10,
                  border: '1px solid #333',
                  background: 'var(--card)',
                }}
              >
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    gap: 8,
                    marginBottom: 8,
                  }}
                >
                  <label style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: '0.9rem' }}>
                    <input
                      type="checkbox"
                      checked={insight.include}
                      onChange={(e) =>
                        setInsights((rows) =>
                          rows.map((r) =>
                            r.id === insight.id ? { ...r, include: e.target.checked } : r,
                          ),
                        )
                      }
                    />
                    Include
                  </label>
                  <div style={{ display: 'flex', gap: 8, fontSize: '0.85rem' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                      <input
                        type="radio"
                        name={`style-${insight.id}`}
                        checked={insight.style === 'qa'}
                        onChange={() =>
                          setInsights((rows) =>
                            rows.map((r) =>
                              r.id === insight.id ? { ...r, style: 'qa' } : r,
                            ),
                          )
                        }
                      />
                      Q / A
                    </label>
                    <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                      <input
                        type="radio"
                        name={`style-${insight.id}`}
                        checked={insight.style === 'true_false'}
                        onChange={() =>
                          setInsights((rows) =>
                            rows.map((r) =>
                              r.id === insight.id ? { ...r, style: 'true_false' } : r,
                            ),
                          )
                        }
                      />
                      True / False
                    </label>
                  </div>
                </div>
                <textarea
                  value={insight.text}
                  onChange={(e) =>
                    setInsights((rows) =>
                      rows.map((r) =>
                        r.id === insight.id ? { ...r, text: e.target.value } : r,
                      ),
                    )
                  }
                  rows={2}
                  style={{
                    width: '100%',
                    padding: 8,
                    background: '#050505',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: 'inherit',
                    fontSize: '0.9rem',
                  }}
                />
              </div>
            ))}
          </div>
        </form>
      )}
    </section>
  )
}

