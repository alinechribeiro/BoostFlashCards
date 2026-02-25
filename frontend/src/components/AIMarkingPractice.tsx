import { FormEvent, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api, PracticeGradeResult } from '../api/client'

export default function AIMarkingPractice() {
  const { subjectId } = useParams<{ subjectId: string }>()
  const [question, setQuestion] = useState<string>('')
  const [answer, setAnswer] = useState<string>('')
  const [loadingQuestion, setLoadingQuestion] = useState(false)
  const [loadingMark, setLoadingMark] = useState(false)
  const [result, setResult] = useState<PracticeGradeResult | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [handwritingEnabled, setHandwritingEnabled] = useState(false)

  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  const drawing = useRef(false)

  const startDrawing = (x: number, y: number) => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    drawing.current = true
    ctx.beginPath()
    ctx.moveTo(x, y)
  }

  const continueDrawing = (x: number, y: number) => {
    if (!drawing.current) return
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    ctx.lineTo(x, y)
    ctx.strokeStyle = '#ffffff'
    ctx.lineWidth = 2
    ctx.lineCap = 'round'
    ctx.stroke()
  }

  const stopDrawing = () => {
    drawing.current = false
  }

  const handleGetQuestion = async () => {
    if (!subjectId) return
    setLoadingQuestion(true)
    setError(null)
    setResult(null)
    try {
      const res = await api.getPracticeQuestion(Number(subjectId))
      setQuestion(res.question)
      setAnswer('')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to get question from AI')
    } finally {
      setLoadingQuestion(false)
    }
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!subjectId || !question.trim() || !answer.trim()) return
    setLoadingMark(true)
    setError(null)
    try {
      const res = await api.gradePracticeAnswer(Number(subjectId), question, answer)
      setResult(res)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to mark answer with AI')
    } finally {
      setLoadingMark(false)
    }
  }

  const encouragement =
    result && result.score_percentage >= 70
      ? "Amazing work – you're hitting strong GCSE performance here. Keep building on this!"
      : result && result.score_percentage >= 40
      ? "Good effort – you’re on the right track. Focus on the improvement tips and try another question."
      : result
      ? "Every attempt counts. Read the feedback carefully, adjust your approach, and try one more – progress comes from practice."
      : null

  return (
    <section>
      <h1 style={{ marginBottom: 16 }}>AI-marked practice</h1>
      <p style={{ marginBottom: 24, color: 'var(--muted)' }}>
        Get an exam-style question, write your answer, and let the AI examiner mark it and predict a GCSE grade.
      </p>

      <button
        type="button"
        onClick={handleGetQuestion}
        disabled={loadingQuestion}
        style={{
          padding: '10px 20px',
          background: 'var(--accent)',
          border: 'none',
          borderRadius: 8,
          color: 'white',
          fontWeight: 600,
          marginBottom: 16,
        }}
      >
        {loadingQuestion ? 'Asking AI for a question…' : 'Get a question'}
      </button>

      {question && (
        <div
          style={{
            marginBottom: 16,
            padding: 16,
            borderRadius: 12,
            background: 'var(--card)',
            border: '1px solid #333',
          }}
        >
          <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Question</h2>
          <p>{question}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span style={{ color: 'var(--muted)', fontSize: '0.9rem' }}>
            You can type or use handwriting on touch devices.
          </span>
          <label style={{ fontSize: '0.85rem', display: 'flex', alignItems: 'center', gap: 6 }}>
            <input
              type="checkbox"
              checked={handwritingEnabled}
              onChange={(e) => setHandwritingEnabled(e.target.checked)}
            />
            Handwriting scratchpad
          </label>
        </div>

        {handwritingEnabled && (
          <div>
            <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Handwriting pad</span>
            <div
              style={{
                border: '1px solid #333',
                borderRadius: 8,
                background: '#050505',
                touchAction: 'none',
              }}
            >
              <canvas
                ref={canvasRef}
                width={700}
                height={220}
                style={{ width: '100%', height: 220, display: 'block', borderRadius: 8 }}
                onMouseDown={(e) => {
                  const rect = (e.target as HTMLCanvasElement).getBoundingClientRect()
                  startDrawing(e.clientX - rect.left, e.clientY - rect.top)
                }}
                onMouseMove={(e) => {
                  const rect = (e.target as HTMLCanvasElement).getBoundingClientRect()
                  continueDrawing(e.clientX - rect.left, e.clientY - rect.top)
                }}
                onMouseUp={stopDrawing}
                onMouseLeave={stopDrawing}
                onTouchStart={(e) => {
                  const touch = e.touches[0]
                  const rect = (e.target as HTMLCanvasElement).getBoundingClientRect()
                  startDrawing(touch.clientX - rect.left, touch.clientY - rect.top)
                }}
                onTouchMove={(e) => {
                  const touch = e.touches[0]
                  const rect = (e.target as HTMLCanvasElement).getBoundingClientRect()
                  continueDrawing(touch.clientX - rect.left, touch.clientY - rect.top)
                }}
                onTouchEnd={stopDrawing}
              />
            </div>
            <p style={{ marginTop: 4, fontSize: '0.8rem', color: 'var(--muted)' }}>
              On iPad/phone, you can handwrite directly here or use your device&apos;s handwriting-to-text keyboard
              in the answer box below.
            </p>
          </div>
        )}

        <label>
          <span style={{ display: 'block', marginBottom: 4, color: 'var(--muted)' }}>Your answer</span>
          <textarea
            value={answer}
            onChange={(e) => setAnswer(e.target.value)}
            rows={8}
            placeholder="Write your exam-style answer here…"
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
          disabled={loadingMark || !question || !answer.trim()}
          style={{
            padding: '10px 20px',
            background: 'var(--accent)',
            border: 'none',
            borderRadius: 8,
            color: 'white',
            fontWeight: 600,
            alignSelf: 'flex-end',
            minWidth: 160,
          }}
        >
          {loadingMark ? 'Marking answer…' : 'Mark my answer'}
        </button>
      </form>

      {result && (
        <section style={{ marginTop: 24 }}>
          <h2 style={{ fontSize: '1rem', marginBottom: 8 }}>Predicted grade</h2>
          <p style={{ marginBottom: 4 }}>
            <strong>Score:</strong> {result.score} / {result.max_score} ({result.score_percentage.toFixed(1)}%)
          </p>
          <p style={{ marginBottom: 4 }}>
            <strong>GCSE grade:</strong> {result.grade}
            {result.grade_band ? ` – ${result.grade_band}` : ''}
          </p>
          <p style={{ marginTop: 12, marginBottom: 8 }}>
            <strong>Feedback</strong>
          </p>
          <p style={{ marginBottom: 8 }}>{result.feedback}</p>
          {result.strengths && (
            <>
              <p style={{ marginBottom: 4 }}>
                <strong>What went well</strong>
              </p>
              <p style={{ marginBottom: 8 }}>{result.strengths}</p>
            </>
          )}
          {result.improvements && (
            <>
              <p style={{ marginBottom: 4 }}>
                <strong>Next steps</strong>
              </p>
              <p style={{ marginBottom: 8 }}>{result.improvements}</p>
            </>
          )}
          {encouragement && (
            <p style={{ marginTop: 12, color: 'var(--accent)' }}>
              {encouragement}
            </p>
          )}
        </section>
      )}
    </section>
  )
}

