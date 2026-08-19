import { useState, useEffect } from 'react'
import './App.css'

const TYPES = [
  { label: 'all', value: null },
  { label: 'scholarship', value: 'scholarship' },
  { label: 'hackathons', value: 'hackathon' },
  { label: 'internships', value: 'internship' },
  { label: 'research / extra', value: 'research' },
]

const API_BASE = (import.meta.env.VITE_API_BASE || '').replace(/\/+$/, '')
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'


function App() {
  const [activeType, setActiveType] = useState(null)
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const controller = new AbortController()
    setLoading(true)
    setError(null)

    // Use mock data when USE_MOCK is true or backend is unavailable
    if (USE_MOCK) {
      const filtered = activeType
        ? MOCK_DATA.filter((d) => d.type === activeType)
        : MOCK_DATA
      setTimeout(() => {
        setItems(filtered)
        setLoading(false)
      }, 300) // simulate network delay
      return
    }

    const url = activeType ? `${API_BASE}/?type=${activeType}` : `${API_BASE}/`

    fetch(url, { signal: controller.signal })
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json()
      })
      .then((data) => {
        // backend returns { data: { items: [], page, per_page, total } }
        const list = Array.isArray(data) ? data
          : Array.isArray(data?.data?.items) ? data.data.items
            : Array.isArray(data?.data) ? data.data
              : Array.isArray(data?.items) ? data.items
                : []
        setItems(list)
        setLoading(false)
      })
      .catch((err) => {
        if (err.name !== 'AbortError') {
          setError(err.message)
          setLoading(false)
        }
      })

    return () => controller.abort()
  }, [activeType])

  const handleChipClick = (value) => {
    setActiveType(value)
  }

  return (
    <>
      {/* ── Header ──────────────────────────────── */}
      <header className="header">
        <span className="header__logo">ScholarBuddy</span>
        <span className="header__tagline">opportunities for students</span>
        <span className="header__spacer" />
        <button className="header__btn" id="login-btn">login</button>
        <span className="header__sep">|</span>
        <button className="header__btn" id="signup-btn">signup</button>
      </header>

      {/* ── Filter chips ────────────────────────── */}
      <nav className="filters" aria-label="Filter by type">
        {TYPES.map((t) => (
          <button
            key={t.value}
            className={`chip ${activeType === t.value ? 'chip--active' : ''}`}
            onClick={() => handleChipClick(t.value)}
            aria-pressed={activeType === t.value}
          >
            {t.label}
          </button>
        ))}
      </nav>

      {/* ── Data table ──────────────────────────── */}
      <div className="table-wrap">
        {loading ? (
          <p className="state-msg">
            loading<span className="loading-dots"></span>
          </p>
        ) : error ? (
          <p className="state-msg state-msg--error">error: {error}</p>
        ) : items.length === 0 ? (
          <p className="state-msg">nothing here yet.</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th className="col-id">#</th>
                <th className="col-title">title</th>
                <th className="col-link">link</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id}>
                  <td className="col-id">{item.id}</td>
                  <td className="col-title">{item.title}</td>
                  <td className="col-link">
                    <a href={item.link} target="_blank" rel="noopener noreferrer">
                      {item.link}
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* ── Footer ──────────────────────────────── */}
      <footer className="footer">scholarbuddy</footer>
    </>
  )
}

export default App
