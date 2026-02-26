import { FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

export default function Signup() {
  const { signup } = useAuth()
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<'student' | 'tutor'>('student')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await signup({ email, password, name, role })
      navigate('/')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to sign up')
    } finally {
      setLoading(false)
    }
  }

  const social = (provider: string) => {
    window.location.href = `/api/auth/${provider}/redirect`
  }

  return (
    <section>
      <h1 style={{ marginBottom: 16 }}>Sign up</h1>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12, maxWidth: 360 }}>
        <label>
          <span style={{ display: 'block', marginBottom: 4 }}>Name</span>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            style={{
              width: '100%',
              padding: 10,
              borderRadius: 8,
              border: '1px solid #333',
              background: 'var(--card)',
              color: 'inherit',
            }}
          />
        </label>
        <label>
          <span style={{ display: 'block', marginBottom: 4 }}>Email</span>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            style={{
              width: '100%',
              padding: 10,
              borderRadius: 8,
              border: '1px solid #333',
              background: 'var(--card)',
              color: 'inherit',
            }}
          />
        </label>
        <label>
          <span style={{ display: 'block', marginBottom: 4 }}>Password</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            style={{
              width: '100%',
              padding: 10,
              borderRadius: 8,
              border: '1px solid #333',
              background: 'var(--card)',
              color: 'inherit',
            }}
          />
        </label>
        <fieldset style={{ border: 'none', padding: 0, marginTop: 8 }}>
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
          {loading ? 'Signing up…' : 'Sign up'}
        </button>
      </form>

      <p style={{ marginTop: 16, fontSize: '0.9rem' }}>
        Already have an account?{' '}
        <Link to="/login" style={{ color: 'var(--accent)' }}>
          Log in
        </Link>
      </p>

      <hr style={{ margin: '24px 0', borderColor: '#333' }} />

      <p style={{ marginBottom: 8, fontSize: '0.9rem', color: 'var(--muted)' }}>Or sign up with</p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 8, maxWidth: 360 }}>
        <button type="button" onClick={() => social('google')} style={socialBtnStyle}>
          Google
        </button>
        <button type="button" onClick={() => social('facebook')} style={socialBtnStyle}>
          Facebook
        </button>
      </div>
    </section>
  )
}

const socialBtnStyle = {
  padding: 10,
  borderRadius: 8,
  border: '1px solid #333',
  background: 'var(--card)',
  color: 'inherit',
  textAlign: 'left' as const,
}


