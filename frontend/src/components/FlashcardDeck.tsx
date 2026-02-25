import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, Flashcard } from '../api/client'

export default function FlashcardDeck() {
  const { topicId } = useParams<{ topicId: string }>()
  const [cards, setCards] = useState<Flashcard[]>([])
  const [subjectId, setSubjectId] = useState<number | null>(null)
  const [index, setIndex] = useState(0)
  const [flipped, setFlipped] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editing, setEditing] = useState(false)
  const [editFront, setEditFront] = useState('')
  const [editBack, setEditBack] = useState('')
  const [saving, setSaving] = useState(false)
  const [viewMode, setViewMode] = useState<'single' | 'gallery'>('single')
  const [pageSize, setPageSize] = useState(10)
  const [page, setPage] = useState(1)
  const [galleryAnswerMode, setGalleryAnswerMode] = useState<'always' | 'onClick'>('onClick')
  const [revealedIds, setRevealedIds] = useState<number[]>([])

  useEffect(() => {
    if (!topicId) return
    const id = Number(topicId)
    Promise.all([api.getFlashcards(id), api.getTopic(id)])
      .then(([list, topic]) => {
        setCards(Array.isArray(list) ? list : [])
        setSubjectId(topic.subject_id)
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [topicId])

  if (loading) return <p style={{ color: 'var(--muted)' }}>Loading cards…</p>
  if (error) return <p style={{ color: '#ef4444' }}>Error: {error}</p>

  const card = cards[index]
  const hasNext = index < cards.length - 1
  const hasPrev = index > 0

  const totalPages = cards.length > 0 ? Math.max(1, Math.ceil(cards.length / pageSize)) : 1
  const currentPage = Math.min(page, totalPages)
  const start = (currentPage - 1) * pageSize
  const pageCards = cards.slice(start, start + pageSize)

  return (
    <section>
      <Link to={subjectId != null ? `/subjects/${subjectId}/topics` : '#'} style={{ marginBottom: 16, display: 'inline-block' }}>← Topics</Link>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', marginBottom: 16 }}>
        <h1>Flashcards</h1>
        {topicId && (
          <Link
            to={`/topics/${topicId}/ai/text`}
            style={{ fontSize: '0.9rem', color: 'var(--accent)' }}
          >
            Text → flashcards
          </Link>
        )}
      </div>

      {cards.length > 0 && (
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 16,
            gap: 12,
            flexWrap: 'wrap',
          }}
        >
          <div style={{ display: 'flex', gap: 8, fontSize: '0.9rem' }}>
            <button
              type="button"
              onClick={() => { setViewMode('single'); setFlipped(false) }}
              style={{
                padding: '6px 12px',
                borderRadius: 999,
                border: '1px solid #333',
                background: viewMode === 'single' ? 'var(--accent)' : 'transparent',
                color: viewMode === 'single' ? '#fff' : 'inherit',
              }}
            >
              Study mode
            </button>
            <button
              type="button"
              onClick={() => setViewMode('gallery')}
              style={{
                padding: '6px 12px',
                borderRadius: 999,
                border: '1px solid #333',
                background: viewMode === 'gallery' ? 'var(--card)' : 'transparent',
                color: 'inherit',
              }}
            >
              Gallery
            </button>
          </div>
          {viewMode === 'gallery' && (
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, fontSize: '0.85rem' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                Cards per page:
                <select
                  value={pageSize}
                  onChange={(e) => {
                    const nextSize = Number(e.target.value) || 10
                    setPageSize(nextSize)
                    setPage(1)
                  }}
                  style={{
                    background: '#050505',
                    borderRadius: 6,
                    border: '1px solid #333',
                    color: 'inherit',
                    padding: '4px 8px',
                  }}
                >
                  {[10, 20, 50, 100].map((n) => (
                    <option key={n} value={n}>
                      {n}
                    </option>
                  ))}
                </select>
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                Show answers:
                <select
                  value={galleryAnswerMode}
                  onChange={(e) => {
                    const mode = e.target.value === 'always' ? 'always' : 'onClick'
                    setGalleryAnswerMode(mode)
                    if (mode === 'always') setRevealedIds([])
                  }}
                  style={{
                    background: '#050505',
                    borderRadius: 6,
                    border: '1px solid #333',
                    color: 'inherit',
                    padding: '4px 8px',
                  }}
                >
                  <option value="onClick">On tap</option>
                  <option value="always">Always</option>
                </select>
              </label>
            </div>
          )}
        </div>
      )}

      {cards.length === 0 ? (
        <p style={{ color: 'var(--muted)' }}>
          No cards yet.{' '}
          <Link to={`/topics/${topicId}/cards/new`}>Create one</Link>.
        </p>
      ) : viewMode === 'single' ? (
        <>
          {editing && card ? (
            <div
              style={{
                background: 'var(--card)',
                borderRadius: 16,
                padding: 24,
                marginBottom: 24,
                display: 'flex',
                flexDirection: 'column',
                gap: 12,
              }}
            >
              <p style={{ color: 'var(--muted)' }}>Editing card {index + 1} / {cards.length}</p>
              <label>
                <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Question (front)</span>
                <textarea
                  value={editFront}
                  onChange={(e) => setEditFront(e.target.value)}
                  rows={3}
                  style={{
                    width: '100%',
                    padding: 10,
                    background: '#050505',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: 'inherit',
                  }}
                />
              </label>
              <label>
                <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Answer (back)</span>
                <textarea
                  value={editBack}
                  onChange={(e) => setEditBack(e.target.value)}
                  rows={3}
                  style={{
                    width: '100%',
                    padding: 10,
                    background: '#050505',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: 'inherit',
                  }}
                />
              </label>
              <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
                <button
                  type="button"
                  onClick={() => setEditing(false)}
                  disabled={saving}
                  style={{
                    padding: '8px 16px',
                    background: 'transparent',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: 'inherit',
                  }}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  disabled={saving || !editFront.trim() || !editBack.trim()}
                  onClick={async () => {
                    if (!card) return
                    setSaving(true)
                    try {
                      const updated = await api.updateFlashcard(card.id, editFront.trim(), editBack.trim())
                      setCards((prev) =>
                        prev.map((c) => (c.id === card.id ? updated : c)),
                      )
                      setEditing(false)
                      setFlipped(false)
                    } catch (e) {
                      // basic error surface
                      alert(e instanceof Error ? e.message : 'Failed to save card')
                    } finally {
                      setSaving(false)
                    }
                  }}
                  style={{
                    padding: '8px 16px',
                    background: 'var(--accent)',
                    borderRadius: 8,
                    border: 'none',
                    color: '#fff',
                    fontWeight: 600,
                  }}
                >
                  {saving ? 'Saving…' : 'Save'}
                </button>
              </div>
            </div>
          ) : (
            <div
              onClick={() => setFlipped((f) => !f)}
              style={{
                position: 'relative',
                background: 'var(--card)',
                borderRadius: 16,
                padding: 32,
                minHeight: 180,
                cursor: 'pointer',
                marginBottom: 16,
              }}
            >
              {card && (
                <span
                  style={{
                    position: 'absolute',
                    top: 12,
                    right: 16,
                    padding: '4px 10px',
                    borderRadius: 999,
                    fontSize: '0.75rem',
                    border: '1px solid #333',
                    background: card.status === 'confident' ? '#14532d' : '#111827',
                    color: card.status === 'confident' ? '#bbf7d0' : '#e5e7eb',
                  }}
                >
                  {card.status === 'confident' ? 'Confident' : 'Not confident yet'}
                </span>
              )}
              <p style={{ color: 'var(--muted)', marginBottom: 8 }}>
                {flipped ? 'Answer' : 'Question'}
              </p>
              <p style={{ fontSize: '1.25rem' }}>
                {flipped ? card?.back : card?.front}
              </p>
            </div>
          )}
          {card && (
            <div
              style={{
                marginBottom: 16,
                display: 'flex',
                justifyContent: 'center',
                gap: 8,
              }}
            >
              <button
                type="button"
                onClick={async (e) => {
                  e.stopPropagation()
                  try {
                    const updated = await api.setFlashcardStatus(card.id, 'not_yet')
                    setCards((prev) => prev.map((c) => (c.id === updated.id ? updated : c)))
                  } catch (err) {
                    // lightweight surface
                    alert(err instanceof Error ? err.message : 'Failed to update card')
                  }
                }}
                style={{
                  padding: '6px 12px',
                  borderRadius: 999,
                  border: '1px solid #374151',
                  background: '#020617',
                  color: '#e5e7eb',
                  fontSize: '0.85rem',
                }}
              >
                ← Not confident yet
              </button>
              <button
                type="button"
                onClick={async (e) => {
                  e.stopPropagation()
                  try {
                    const updated = await api.setFlashcardStatus(card.id, 'confident')
                    // remove from active list until it is due again
                    setCards((prev) => {
                      const next = prev.filter((c) => c.id !== updated.id)
                      let newIndex = index
                      if (newIndex >= next.length) newIndex = Math.max(0, next.length - 1)
                      setIndex(newIndex)
                      setFlipped(false)
                      return next
                    })
                  } catch (err) {
                    alert(err instanceof Error ? err.message : 'Failed to update card')
                  }
                }}
                style={{
                  padding: '6px 12px',
                  borderRadius: 999,
                  border: 'none',
                  background: 'var(--accent)',
                  color: '#fff',
                  fontSize: '0.85rem',
                  fontWeight: 600,
                }}
              >
                Confident →
              </button>
            </div>
          )}
          <div style={{ display: 'flex', gap: 12, justifyContent: 'space-between' }}>
            <button
              type="button"
              disabled={!hasPrev}
              onClick={() => { setIndex((i) => i - 1); setFlipped(false) }}
              style={{
                padding: '10px 20px',
                background: hasPrev ? 'var(--card)' : 'transparent',
                border: '1px solid var(--card)',
                borderRadius: 8,
                color: 'inherit',
              }}
            >
              Previous
            </button>
            <span style={{ color: 'var(--muted)' }}>
              {index + 1} / {cards.length}
            </span>
            <button
              type="button"
              disabled={!hasNext}
              onClick={() => { setIndex((i) => i + 1); setFlipped(false) }}
              style={{
                padding: '10px 20px',
                background: hasNext ? 'var(--accent)' : 'transparent',
                border: 'none',
                borderRadius: 8,
                color: 'white',
              }}
            >
              Next
            </button>
          </div>
          <div style={{ marginTop: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Link to={`/topics/${topicId}/cards/new`}>+ Add another card</Link>
            {card && (
              <div style={{ display: 'flex', gap: 8 }}>
                <button
                  type="button"
                  onClick={() => {
                    setEditing(true)
                    setEditFront(card.front)
                    setEditBack(card.back)
                  }}
                  style={{
                    padding: '8px 14px',
                    background: 'var(--card)',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: 'inherit',
                    fontSize: '0.9rem',
                  }}
                >
                  Edit card
                </button>
                <button
                  type="button"
                  onClick={async () => {
                    if (!card) return
                    if (!window.confirm('Delete this card?')) return
                    try {
                      await api.deleteFlashcard(card.id)
                      setCards((prev) => {
                        const next = prev.filter((c) => c.id !== card.id)
                        let newIndex = index
                        if (newIndex >= next.length) newIndex = Math.max(0, next.length - 1)
                        setIndex(newIndex)
                        setFlipped(false)
                        setEditing(false)
                        return next
                      })
                    } catch (e) {
                      alert(e instanceof Error ? e.message : 'Failed to delete card')
                    }
                  }}
                  style={{
                    padding: '8px 14px',
                    background: '#111',
                    borderRadius: 8,
                    border: '1px solid #333',
                    color: '#f97373',
                    fontSize: '0.9rem',
                  }}
                >
                  Delete
                </button>
              </div>
            )}
          </div>
        </>
      ) : (
        <>
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
              gap: 12,
              marginBottom: 16,
            }}
          >
            {pageCards.map((c) => {
              const revealed = galleryAnswerMode === 'always' || revealedIds.includes(c.id)
              return (
                <div
                  key={c.id}
                  style={{
                    background: 'var(--card)',
                    borderRadius: 12,
                    padding: 14,
                    border: '1px solid #333',
                    fontSize: '0.9rem',
                  }}
                >
                  <p style={{ color: 'var(--muted)', marginBottom: 6 }}>Q</p>
                  <p style={{ marginBottom: 8 }}>{c.front}</p>
                  <p style={{ color: 'var(--muted)', marginBottom: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span>A</span>
                    {galleryAnswerMode === 'onClick' && (
                      <button
                        type="button"
                        onClick={() =>
                          setRevealedIds((prev) =>
                            revealed ? prev.filter((id) => id !== c.id) : [...prev, c.id],
                          )
                        }
                        style={{
                          border: 'none',
                          background: 'transparent',
                          color: 'var(--accent)',
                          fontSize: '0.8rem',
                          cursor: 'pointer',
                        }}
                      >
                        {revealed ? 'Hide' : 'Show'}
                      </button>
                    )}
                  </p>
                  <p style={{ minHeight: 40, opacity: revealed ? 1 : 0.25 }}>
                    {revealed ? c.back : 'Tap “Show” to reveal the answer'}
                  </p>
                </div>
              )
            })}
          </div>
          {cards.length > pageSize && (
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                gap: 12,
              }}
            >
              <button
                type="button"
                disabled={currentPage <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                style={{
                  padding: '8px 14px',
                  background: currentPage > 1 ? 'var(--card)' : 'transparent',
                  borderRadius: 8,
                  border: '1px solid #333',
                  color: 'inherit',
                }}
              >
                Previous page
              </button>
              <span style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
                Page {currentPage} / {totalPages}
              </span>
              <button
                type="button"
                disabled={currentPage >= totalPages}
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                style={{
                  padding: '8px 14px',
                  background: currentPage < totalPages ? 'var(--accent)' : 'transparent',
                  borderRadius: 8,
                  border: 'none',
                  color: currentPage < totalPages ? '#fff' : 'inherit',
                }}
              >
                Next page
              </button>
            </div>
          )}
        </>
      )}
    </section>
  )
}
