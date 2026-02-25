import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import SubjectList from './SubjectList'
import { api } from '../api/client'

vi.mock('../api/client', () => ({
  api: {
    getSubjects: vi.fn(),
  },
}))

describe('SubjectList', () => {
  it('shows loading then subjects', async () => {
    vi.mocked(api.getSubjects).mockResolvedValue([
      { id: 1, name: 'Mathematics', slug: 'mathematics', created_at: '' },
      { id: 2, name: 'Biology', slug: 'biology', created_at: '' },
    ])

    render(
      <MemoryRouter>
        <SubjectList />
      </MemoryRouter>
    )

    expect(screen.getByText(/Loading subjects/)).toBeInTheDocument()

    expect(await screen.findByText('Mathematics')).toBeInTheDocument()
    expect(screen.getByText('Biology')).toBeInTheDocument()
  })

  it('shows error when fetch fails', async () => {
    vi.mocked(api.getSubjects).mockRejectedValue(new Error('Network error'))

    render(
      <MemoryRouter>
        <SubjectList />
      </MemoryRouter>
    )

    expect(await screen.findByText(/Error: Network error/)).toBeInTheDocument()
  })
})
