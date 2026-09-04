import { useEffect, useState } from "react";
import { createClient } from "@supabase/supabase-js";
import "./App.css";

const TYPES = [
  { label: "all", value: null },
  { label: "scholarships", value: "scholarship" },
  { label: "hackathons", value: "hackathon" },
  { label: "internships", value: "internship" },
  { label: "research", value: "research" },
  { label: "extras", value: "extras" },
];
const API_BASE = (import.meta.env.VITE_API_BASE || "").replace(/\/+$/, "");
const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
const supabaseKey = import.meta.env.VITE_SUPABASE_PUBLISHABLE_KEY;
const supabase =
  supabaseUrl && supabaseKey ? createClient(supabaseUrl, supabaseKey) : null;
const apiURL = (path) => `${API_BASE}${path}`;
const unwrapItems = (payload) => payload?.data?.items || [];

function App() {
  const [activeType, setActiveType] = useState(null);
  const [items, setItems] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [session, setSession] = useState(null);
  const [adminToken, setAdminToken] = useState(() =>
    sessionStorage.getItem("scholarbuddy_admin_token"),
  );
  const [adminItems, setAdminItems] = useState([]);
  const [adminError, setAdminError] = useState(null);
  const [view, setView] = useState(() =>
    window.location.pathname === "/abbujaan" ? "admin-login" : "home",
  );

  useEffect(() => {
    if (!supabase) return;
    supabase.auth.getSession().then(({ data }) => setSession(data.session));
    const { data: listener } = supabase.auth.onAuthStateChange(
      (_event, nextSession) => setSession(nextSession),
    );
    return () => listener.subscription.unsubscribe();
  }, []);

  useEffect(() => {
    if (view !== "home") return;
    const controller = new AbortController();
    const query = activeType ? `?type=${encodeURIComponent(activeType)}` : "";
    fetch(apiURL(`/api/v1/opportunities${query}`), {
      signal: controller.signal,
    })
      .then((res) =>
        res.ok ? res.json() : Promise.reject(new Error(`HTTP ${res.status}`)),
      )
      .then((payload) => setItems(unwrapItems(payload)))
      .catch((err) => err.name !== "AbortError" && setError(err.message))
      .finally(() => !controller.signal.aborted && setLoading(false));
    return () => controller.abort();
  }, [activeType, view]);

  const signInWithGoogle = async () => {
    if (!supabase) return setError("Google login is not configured yet.");
    const { error: oauthError } = await supabase.auth.signInWithOAuth({
      provider: "google",
      options: { redirectTo: window.location.origin },
    });
    if (oauthError) {
      setError(
        oauthError.message.includes("Unsupported provider")
          ? "Google login is disabled in Supabase. Enable the Google provider in Authentication → Providers."
          : oauthError.message,
      );
    }
  };
  const openOpportunity = async (id) => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(apiURL(`/api/v1/opportunities/${id}`));
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const payload = await res.json();
      setSelected(payload.data);
      setView("detail");
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };
  const loadAdmin = async (token = adminToken) => {
    if (!token) return setView("admin-login");
    setAdminError(null);
    const res = await fetch(apiURL("/api/v1/abbujaan/opportunities"), {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) return setAdminError("Unable to load pending opportunities.");
    const payload = await res.json();
    setAdminItems(unwrapItems(payload));
    setView("admin");
  };
  const adminLogin = async (event) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    setAdminError(null);
    const res = await fetch(apiURL("/api/v1/abbujaan/login"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: form.get("username"),
        password: form.get("password"),
      }),
    });
    if (!res.ok) return setAdminError("Invalid username or password.");
    const payload = await res.json();
    const token = payload.data.access_token;
    sessionStorage.setItem("scholarbuddy_admin_token", token);
    setAdminToken(token);
    loadAdmin(token);
  };
  const review = async (id, action) => {
    const res = await fetch(
      apiURL(
        `/api/v1/abbujaan/opportunities/${id}${action === "approve" ? "/approve" : ""}`,
      ),
      {
        method: action === "approve" ? "PATCH" : "DELETE",
        headers: { Authorization: `Bearer ${adminToken}` },
      },
    );
    if (!res.ok) return setAdminError(`Could not ${action} this opportunity.`);
    setAdminItems((current) => current.filter((item) => item.id !== id));
  };
  const logout = async () => {
    if (supabase) await supabase.auth.signOut();
    setSession(null);
  };

  return (
    <>
      <Header
        session={session}
        onGoogleLogin={signInWithGoogle}
        onLogout={logout}
        onHome={() => {
          history.pushState({}, "", "/");
          setError(null);
          setView("home");
        }}
      />
      {view === "home" && (
        <Home
          activeType={activeType}
          onTypeChange={(type) => {
            setError(null);
            setLoading(true);
            setActiveType(type);
          }}
          items={items}
          loading={loading}
          error={error}
          onOpen={openOpportunity}
        />
      )}
      {view === "detail" && (
        <Detail
          item={selected}
          loading={loading}
          error={error}
          onBack={() => {
            setError(null);
            setView("home");
          }}
        />
      )}
      {view === "admin-login" && (
        <AdminLogin
          error={adminError}
          onSubmit={adminLogin}
          onBack={() => setView("home")}
        />
      )}
      {view === "admin" && (
        <AdminDashboard
          items={adminItems}
          error={adminError}
          onReview={review}
          onBack={() => setView("home")}
        />
      )}
      <footer className="footer">scholarbuddy</footer>
    </>
  );
}

function Header({ session, onGoogleLogin, onLogout, onHome }) {
  return (
    <header className="header">
      <button className="header__logo" onClick={onHome}>
        ScholarBuddy
      </button>
      <span className="header__tagline">opportunities for students</span>
      <span className="header__spacer" />
      {session ? (
        <button className="header__btn" onClick={onLogout}>
          logout
        </button>
      ) : (
        <button className="header__btn" onClick={onGoogleLogin}>
          login with google
        </button>
      )}
    </header>
  );
}
function Home({ activeType, onTypeChange, items, loading, error, onOpen }) {
  return (
    <>
      <nav className="filters" aria-label="Filter by type">
        {TYPES.map((type) => (
          <button
            key={type.label}
            className={`chip ${activeType === type.value ? "chip--active" : ""}`}
            onClick={() => onTypeChange(type.value)}
          >
            {type.label}
          </button>
        ))}
      </nav>
      <p className="ranking-note">ranked by clicks</p>
      <div className="table-wrap">
        {loading ? (
          <State text="loading" />
        ) : error ? (
          <State text={`error: ${error}`} error />
        ) : items.length === 0 ? (
          <State text="nothing here yet." />
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th className="col-id">#</th>
                <th>title</th>
                <th>type</th>
                <th className="col-clicks">clicks</th>
                <th>link</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item, index) => (
                <tr key={item.id}>
                  <td className="col-id">{index + 1}</td>
                  <td>
                    <button
                      className="title-button"
                      onClick={() => onOpen(item.id)}
                    >
                      {item.title}
                    </button>
                  </td>
                  <td>{item.types.join(", ")}</td>
                  <td className="col-clicks">{item.click_count}</td>
                  <td className="col-link">
                    <a href={item.link} target="_blank" rel="noreferrer">
                      {item.link}
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  );
}
function Detail({ item, loading, error, onBack }) {
  if (loading) return <State text="loading" />;
  if (error || !item)
    return <State text={`error: ${error || "not found"}`} error />;
  return (
    <main className="detail">
      <button className="back-button" onClick={onBack}>
        ← back
      </button>
      <h1>{item.title}</h1>
      <p className="detail__meta">
        {item.types.join(" · ")} · {item.click_count} clicks
      </p>
      <p className="detail__description">{item.description}</p>
      <DetailSection title="eligibility" value={item.eligibility} />
      <DetailSection title="how to apply" value={item.steps} />
      <DetailSection title="benefits" value={item.benefits} />
      <DetailSection title="referral" value={item.referral} />
      <a
        className="detail__link"
        href={item.link}
        target="_blank"
        rel="noreferrer"
      >
        visit opportunity →
      </a>
    </main>
  );
}
const DetailSection = ({ title, value }) =>
  value ? (
    <section className="detail__section">
      <h2>{title}</h2>
      <p>{value}</p>
    </section>
  ) : null;
const State = ({ text, error }) => (
  <p className={`state-msg ${error ? "state-msg--error" : ""}`}>{text}</p>
);
function AdminLogin({ error, onSubmit, onBack }) {
  return (
    <main className="admin-panel">
      <button className="back-button" onClick={onBack}>
        ← back
      </button>
      <h1>abbujaan</h1>
      <form onSubmit={onSubmit} className="admin-form">
        <label>
          username
          <input name="username" required autoComplete="username" />
        </label>
        <label>
          password
          <input
            name="password"
            type="password"
            required
            autoComplete="current-password"
          />
        </label>
        <button type="submit">login</button>
        {error && <p className="state-msg state-msg--error">{error}</p>}
      </form>
    </main>
  );
}
function AdminDashboard({ items, error, onReview, onBack }) {
  return (
    <main className="admin-panel">
      <button className="back-button" onClick={onBack}>
        ← back
      </button>
      <h1>pending opportunities</h1>
      {error && <p className="state-msg state-msg--error">{error}</p>}
      {items.length === 0 ? (
        <State text="no pending opportunities." />
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>title</th>
              <th>types</th>
              <th>actions</th>
            </tr>
          </thead>
          <tbody>
            {items.map((item) => (
              <tr key={item.id}>
                <td>{item.title}</td>
                <td>{item.types.join(", ")}</td>
                <td className="admin-actions">
                  <button onClick={() => onReview(item.id, "approve")}>
                    approve
                  </button>
                  <button onClick={() => onReview(item.id, "delete")}>
                    delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </main>
  );
}
export default App;
