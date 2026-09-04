import { useState, useEffect } from "react";
import { apiFetch, unwrapData } from "../api";

const TYPES = [
  { label: "all", value: null },
  { label: "scholarships", value: "SCHOLARSHIP" },
  { label: "hackathons", value: "HACKATHON" },
  { label: "internships", value: "INTERNSHIP" },
  { label: "research", value: "RESEARCH" },
  { label: "extras", value: "EXTRAS" },
];

export default function Home({ onOpenOpportunity }) {
  const [activeType, setActiveType] = useState(null);
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const controller = new AbortController();
    const query = activeType ? `?type=${encodeURIComponent(activeType)}` : "";

    setLoading(true);
    setError(null);

    apiFetch(`/api/v1/opportunities${query}`, { signal: controller.signal })
      .then((payload) => setItems(unwrapData(payload)?.items || []))
      .catch((err) => {
        if (err.name !== "AbortError") {
          setError(err.message);
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setLoading(false);
        }
      });

    return () => controller.abort();
  }, [activeType]);

  return (
    <>
      <nav className="filters" aria-label="Filter by type">
        {TYPES.map((type) => (
          <button
            key={type.label}
            className={`chip ${activeType === type.value ? "chip--active" : ""}`}
            onClick={() => setActiveType(type.value)}
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
                      onClick={() => onOpenOpportunity(item.id)}
                    >
                      {item.title}
                    </button>
                  </td>
                  <td>{item.types.join(", ")}</td>
                  <td className="col-clicks">{item.clickCount}</td>
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

const State = ({ text, error }) => (
  <p className={`state-msg ${error ? "state-msg--error" : ""}`}>{text}</p>
);
