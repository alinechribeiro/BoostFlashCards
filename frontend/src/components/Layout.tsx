import { Link } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

export default function Layout({ children }: { children: React.ReactNode }) {
  const { user, logout, loading } = useAuth()

  return (
    <div style={{ maxWidth: 720, margin: '0 auto', padding: 24 }}>
      <header style={{ marginBottom: 32 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', gap: 16 }}>
          <div>
            <Link to="/" style={{ color: 'inherit', textDecoration: 'none', fontSize: '1.5rem', fontWeight: 700 }}>
              BoostFlashCards
            </Link>
            <span style={{ color: 'var(--muted)', marginLeft: 8 }}>GCSE UK</span>
          </div>
          <nav style={{ fontSize: '0.9rem', display: 'flex', gap: 12 }}>
            <Link to="/" style={{ color: 'var(--muted)', textDecoration: 'none' }}>
              Subjects
            </Link>
            <Link to="/tutors" style={{ color: 'var(--muted)', textDecoration: 'none' }}>
              Find a tutor
            </Link>
            <Link to="/ai/subjects" style={{ color: 'var(--accent)', textDecoration: 'none' }}>
              AI creator
            </Link>
            {!loading && (
              user ? (
                <>
                  <span style={{ color: 'var(--muted)' }}>
                    {user.name || user.email} ({user.role})
                  </span>
                  <button
                    type="button"
                    onClick={() => logout()}
                    style={{
                      border: 'none',
                      background: 'transparent',
                      color: 'var(--accent)',
                      cursor: 'pointer',
                    }}
                  >
                    Log out
                  </button>
                </>
              ) : (
                <>
                  <Link to="/login" style={{ color: 'var(--muted)', textDecoration: 'none' }}>
                    Log in
                  </Link>
                  <Link to="/signup" style={{ color: 'var(--accent)', textDecoration: 'none' }}>
                    Sign up
                  </Link>
                </>
              )
            )}
          </nav>
        </div>
      </header>
      <main>{children}</main>
    </div>
  )
}
