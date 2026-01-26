import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from '../routes/pages/Home'
import { Layout } from '../components/layout'
import { AuthProvider } from '../contexts/AuthProvider'

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })
}

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      <AuthProvider>{ui}</AuthProvider>
    </QueryClientProvider>
  )
}

describe('App smoke tests', () => {
  it('renders the home page with welcome message', () => {
    const router = createMemoryRouter(
      [
        {
          path: '/',
          element: <Layout />,
          children: [{ index: true, element: <Home /> }],
        },
      ],
      { initialEntries: ['/'] }
    )

    renderWithProviders(<RouterProvider router={router} />)

    expect(screen.getByText(/Welcome to/i)).toBeInTheDocument()
    expect(screen.getByText('PowerPro')).toBeInTheDocument()
    expect(
      screen.getByText(/Track your powerlifting progress/i)
    ).toBeInTheDocument()
  })

  it('renders navigation elements', () => {
    const router = createMemoryRouter(
      [
        {
          path: '/',
          element: <Layout />,
          children: [{ index: true, element: <Home /> }],
        },
      ],
      { initialEntries: ['/'] }
    )

    renderWithProviders(<RouterProvider router={router} />)

    // Check that navigation links exist (at least in one nav - desktop or mobile)
    expect(screen.getAllByRole('link', { name: /home/i }).length).toBeGreaterThan(0)
    expect(screen.getAllByRole('link', { name: /workout/i }).length).toBeGreaterThan(0)
    expect(screen.getAllByRole('link', { name: /history/i }).length).toBeGreaterThan(0)
    expect(screen.getAllByRole('link', { name: /profile/i }).length).toBeGreaterThan(0)
  })
})
