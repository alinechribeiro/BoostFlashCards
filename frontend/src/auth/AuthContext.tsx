import React, { createContext, useContext, useEffect, useState } from 'react'
import { api, User } from '../api/client'

type AuthContextValue = {
  user: User | null
  loading: boolean
  signup: (data: { email: string; password: string; name: string; role: 'student' | 'tutor' }) => Promise<void>
  login: (data: { email: string; password: string }) => Promise<void>
  logout: () => Promise<void>
  setUser: React.Dispatch<React.SetStateAction<User | null>>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true
    api
      .me()
      .then((u) => {
        if (active) setUser(u)
      })
      .catch(() => {
        if (active) setUser(null)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [])

  const signup = async (data: { email: string; password: string; name: string; role: 'student' | 'tutor' }) => {
    const u = await api.signup(data)
    setUser(u)
  }

  const login = async (data: { email: string; password: string }) => {
    const u = await api.login(data)
    setUser(u)
  }

  const logout = async () => {
    try {
      await api.logout()
    } finally {
      setUser(null)
    }
  }

  const value: AuthContextValue = {
    user,
    loading,
    signup,
    login,
    logout,
    setUser,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}

