import { FormEvent, useState } from 'react'
import { api } from '../api/client'

type Message = {
  role: 'user' | 'assistant'
  content: string
}

export default function AISubjectChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    const prompt = input.trim()
    if (!prompt || loading) return

    setMessages((prev) => [...prev, { role: 'user', content: prompt }])
    setInput('')
    setError(null)
    setLoading(true)
    try {
      const res = await api.createSubjectWithAI(prompt)
      const summary = res.message || `Created subject #${res.subject_id} with ${res.flashcards_created} flashcards.`
      setMessages((prev) => [
        ...prev,
        {
          role: 'assistant',
          content: summary,
        },
      ])
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Failed to talk to AI'
      setError(msg)
      setMessages((prev) => [
        ...prev,
        { role: 'assistant', content: `Something went wrong: ${msg}` },
      ])
    } finally {
      setLoading(false)
    }
  }

  return (
    <section>
      <h1 style={{ marginBottom: 16 }}>AI subject creator</h1>
      <p style={{ marginBottom: 24, color: 'var(--muted)' }}>
        Describe the GCSE subject you want (level, exam board, topic focus) and the AI will create a subject,
        topics and starter flashcards for you.
      </p>

      <div
        style={{
          borderRadius: 12,
          border: '1px solid #333',
          padding: 16,
          marginBottom: 16,
          maxHeight: 260,
          overflowY: 'auto',
          background: 'var(--card)',
        }}
      >
        {messages.length === 0 ? (
          <p style={{ color: 'var(--muted)' }}>
            Example: “Create a GCSE Physics Electricity subject with topics on circuits and resistance, with exam-style
            flashcards.”
          </p>
        ) : (
          <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'flex', flexDirection: 'column', gap: 8 }}>
            {messages.map((m, i) => (
              <li
                key={i}
                style={{
                  alignSelf: m.role === 'user' ? 'flex-end' : 'flex-start',
                  maxWidth: '80%',
                }}
              >
                <div
                  style={{
                    padding: 10,
                    borderRadius: 12,
                    background: m.role === 'user' ? 'var(--accent)' : '#111',
                    color: m.role === 'user' ? '#fff' : 'inherit',
                    fontSize: '0.95rem',
                  }}
                >
                  {m.content}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          rows={3}
          placeholder="Describe the subject you want the AI to create…"
          style={{
            width: '100%',
            padding: 12,
            background: 'var(--card)',
            border: '1px solid #333',
            borderRadius: 8,
            color: 'inherit',
          }}
        />
        {error && <p style={{ color: '#ef4444' }}>{error}</p>}
        <button
          type="submit"
          disabled={loading || !input.trim()}
          style={{
            padding: 12,
            background: 'var(--accent)',
            border: 'none',
            borderRadius: 8,
            color: 'white',
            fontWeight: 600,
            alignSelf: 'flex-end',
            minWidth: 140,
          }}
        >
          {loading ? 'Asking AI…' : 'Ask AI to create'}
        </button>
      </form>
    </section>
  )
}

