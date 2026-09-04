import { useState } from "react";
import { useAuth } from "./useAuth";
import { apiFetch, unwrapData } from "./api";
import Header from "./components/Header";
import Home from "./components/Home";
import Detail from "./components/Detail";
import AdminDashboard from "./components/AdminDashboard";
import "./App.css";

function App() {
  const { user, loading: authLoading, login, logout, isAdmin } = useAuth();
  const [view, setView] = useState("home");
  const [selectedOpportunity, setSelectedOpportunity] = useState(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState(null);

  const handleOpenOpportunity = async (id) => {
    setDetailLoading(true);
    setDetailError(null);
    try {
      const response = await apiFetch(`/api/v1/opportunities/${id}`);
      setSelectedOpportunity(unwrapData(response));
      setView("detail");
    } catch (err) {
      setDetailError(err.message);
    } finally {
      setDetailLoading(false);
    }
  };

  const handleBackToHome = () => {
    setView("home");
    setSelectedOpportunity(null);
    setDetailError(null);
  };

  const handleAdminView = () => {
    if (isAdmin) {
      setView("admin");
    }
  };

  if (authLoading) {
    return (
      <div className="state-msg">
        <p>loading...</p>
      </div>
    );
  }

  return (
    <>
      <Header user={user} onLogin={login} onLogout={logout} onHome={handleBackToHome} />

      {view === "home" && <Home onOpenOpportunity={handleOpenOpportunity} />}

      {view === "detail" && (
        <Detail
          item={selectedOpportunity}
          loading={detailLoading}
          error={detailError}
          onBack={handleBackToHome}
        />
      )}

      {view === "admin" && isAdmin && <AdminDashboard onBack={handleBackToHome} />}

      <footer className="footer">
        scholarbuddy
        {isAdmin && view === "home" && (
          <>
            {" · "}
            <button
              onClick={handleAdminView}
              style={{
                background: "none",
                border: "none",
                color: "inherit",
                cursor: "pointer",
                textDecoration: "underline",
                fontSize: "inherit",
                padding: 0,
              }}
            >
              admin
            </button>
          </>
        )}
      </footer>
    </>
  );
}

export default App;
