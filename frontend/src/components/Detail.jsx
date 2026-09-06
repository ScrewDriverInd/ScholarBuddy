export default function Detail({ item, loading, error, onBack }) {
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
        {item.types.join(" · ")} · {item.clickCount} clicks
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
        visit listing →
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
