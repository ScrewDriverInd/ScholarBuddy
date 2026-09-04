import { useState, useEffect } from "react";
import { apiFetch, unwrapData } from "../api";

export default function AdminDashboard({ onBack }) {
  const [items, setItems] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPendingOpportunities();
  }, []);

  const loadPendingOpportunities = async () => {
    try {
      const response = await apiFetch("/api/v1/abbujaan/opportunities");
      setItems(unwrapData(response)?.items || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleReview = async (id, action) => {
    try {
      await apiFetch(
        `/api/v1/abbujaan/opportunities/${id}${action === "approve" ? "/approve" : ""}`,
        { method: action === "approve" ? "PATCH" : "DELETE" }
      );
      setItems((current) => current.filter((item) => item.id !== id));
    } catch (err) {
      setError(`Could not ${action} this opportunity: ${err.message}`);
    }
  };

  return (
    <main className="admin-panel">
      <button className="back-button" onClick={onBack}>
        ← back
      </button>
      <h1>pending opportunities</h1>
      {error && <p className="state-msg state-msg--error">{error}</p>}
      {loading ? (
        <State text="loading" />
      ) : items.length === 0 ? (
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
                  <button onClick={() => handleReview(item.id, "approve")}>
                    approve
                  </button>
                  <button onClick={() => handleReview(item.id, "delete")}>
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

const State = ({ text, error }) => (
  <p className={`state-msg ${error ? "state-msg--error" : ""}`}>{text}</p>
);
