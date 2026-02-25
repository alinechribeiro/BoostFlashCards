import { FormEvent, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { useAuth } from '../auth/AuthContext'

export default function CompleteSignup() {
  const { setUser } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const search = new URLSearchParams(location.search)
  const pending = search.get('pending') || ''
  const [role, setRole] = useState<'student' | 'tutor'>('student')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!pending) {
      setError('Missing signup token. Please restart social login.')
      return
    }
    setError(null)
    setLoading(true)
    try {
      const user = await api.completeSocialSignup(pending, role)
      setUser(user)
      navigate('/')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to complete signup')
    } finally {
      setLoading(false)
    }
  }

  if (!pending) {
    return <p style={{ color: '#ef4444' }}>Missing signup information. Please restart social login.</p>
  }

  return (
    <section>
      <h1 style={{ marginBottom: 16 }}>Complete your account</h1>
      <p style={{ marginBottom: 16, color: 'var(--muted)' }}>
        Choose whether you are signing up as a student or a tutor. We&apos;ll use this to personalise your dashboard later.
      </p>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12, maxWidth: 360 }}>
        <fieldset style={{ border: 'none', padding: 0 }}>
          <legend style={{ fontSize: '0.9rem', marginBottom: 4 }}>I am signing up as</legend>
          <label style={{ marginRight: 12, fontSize: '0.9rem' }}>
            <input
              type="radio"
              name="role"
              value="student"
              checked={role === 'student'}
              onChange={() => setRole('student')}
            />{' '}
            Student
          </label>
          <label style={{ fontSize: '0.9rem' }}>
            <input
              type="radio"
              name="role"
              value="tutor"
              checked={role === 'tutor'}
              onChange={() => setRole('tutor')}
            />{' '}
            Tutor
          </label>
        </fieldset>
        {error && <p style={{ color: '#ef4444' }}>{error}</p>}
        <button
          type="submit"
          disabled={loading}
          style={{
            padding: 12,
            borderRadius: 8,
            border: 'none',
            background: 'var(--accent)',
            color: '#fff',
            fontWeight: 600,
          }}
        >
          {loading ? 'Saving…' : 'Finish signup'}
        </button>
      </form>
    </section>
  )
}

