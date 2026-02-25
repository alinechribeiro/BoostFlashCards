import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import App from './App'

describe('App', () => {
  it('renders app title', () => {
    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>
    )
    expect(screen.getByText('BoostFlashCards')).toBeInTheDocument()
    expect(screen.getByText('GCSE UK')).toBeInTheDocument()
  })
})
