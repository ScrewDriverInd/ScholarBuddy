import { useState, useEffect } from 'react'
import './App.css'

const TYPES = [
  { label: 'scholarship', value: 'scholarship' },
  { label: 'hackathons', value: 'hackathon' },
  { label: 'internships', value: 'internship' },
  { label: 'research / extra', value: 'research' },
]

const API_BASE = '' // same-origin; adjust if backend lives elsewhere

// ── Dummy data (remove when backend is ready) ────────
const USE_MOCK = true

const MOCK_DATA = [
  { id: 1, title: 'INSPIRE Scholarship for Higher Education', link: 'https://www.online-inspire.gov.in', type: 'scholarship' },
  { id: 2, title: 'Prime Minister\'s Scholarship Scheme', link: 'https://scholarships.gov.in', type: 'scholarship' },
  { id: 3, title: 'Reliance Foundation Scholarship', link: 'https://www.reliancefoundation.org/scholarships', type: 'scholarship' },
  { id: 4, title: 'Adobe GenSolve Hackathon 2026', link: 'https://www.adobe.com/gensolve', type: 'hackathon' },
  { id: 5, title: 'Smart India Hackathon', link: 'https://www.sih.gov.in', type: 'hackathon' },
  { id: 6, title: 'MLH Global Hack Week', link: 'https://ghw.mlh.io', type: 'hackathon' },
  { id: 7, title: 'Google STEP Internship', link: 'https://careers.google.com/students', type: 'internship' },
  { id: 8, title: 'Microsoft Explore Program', link: 'https://careers.microsoft.com/students', type: 'internship' },
  { id: 9, title: 'Mitacs Globalink Research Internship', link: 'https://www.mitacs.ca/globalink', type: 'internship' },
  { id: 10, title: 'CERN Summer Student Programme', link: 'https://home.cern/summer-student-programme', type: 'research' },
  { id: 11, title: 'IISC Summer Research Fellowship', link: 'https://www.iisc.ac.in/srfp', type: 'research' },
  { id: 12, title: 'NASA L\'SPACE Academy', link: 'https://www.lspace.asu.edu', type: 'research' },
]

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
        setItems(data)
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
    setActiveType((prev) => (prev === value ? null : value))
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
