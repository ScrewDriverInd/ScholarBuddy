export default function Header({ user, onLogin, onLogout, onHome }) {
  return (
    <header className="header">
      <button className="header__logo" onClick={onHome}>
        ScholarBuddy
      </button>
      <span className="header__tagline">opportunities for students</span>
      <span className="header__spacer" />
      {user ? (
        <button className="header__btn" onClick={onLogout}>
          logout
        </button>
      ) : (
        <button className="header__btn" onClick={onLogin}>
          login with google
        </button>
      )}
    </header>
  );
}
