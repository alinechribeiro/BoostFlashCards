import { FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login({ email, password })
      navigate('/')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to log in')
    } finally {
      setLoading(false)
    }
  }

  const social = (provider: string) => {
    window.location.href = `/api/auth/${provider}/redirect`
  }

  return (
    <section>
      <h1 style={{ marginBottom: 16 }}>Log in</h1>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12, maxWidth: 360 }}>
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
          {loading ? 'Logging in…' : 'Log in'}
        </button>
      </form>

      <p style={{ marginTop: 16, fontSize: '0.9rem' }}>
        No account yet?{' '}
        <Link to="/signup" style={{ color: 'var(--accent)' }}>
          Sign up
        </Link>
      </p>

      <hr style={{ margin: '24px 0', borderColor: '#333' }} />

      <p style={{ marginBottom: 8, fontSize: '0.9rem', color: 'var(--muted)' }}>Or continue with</p>
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

const socialBtnStyle: React.CSSProperties = {
  padding: 10,
  borderRadius: 8,
  border: '1px solid #333',
  background: 'var(--card)',
  color: 'inherit',
  textAlign: 'left',
}

