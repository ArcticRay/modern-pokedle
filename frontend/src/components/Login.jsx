function Login({ onLogin }) {
  const handleGitHubLogin = () => {
    window.location.href = 'http://localhost:8080/api/v1/auth/github'
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        background: '#1a1a2e',
      }}
    >
      <h1 style={{ color: 'white', fontSize: '3rem', marginBottom: '2rem' }}>
        Pokédle 🎮
      </h1>
      <button
        onClick={handleGitHubLogin}
        style={{
          padding: '1rem 2rem',
          fontSize: '1.2rem',
          background: '#24292e',
          color: 'white',
          border: 'none',
          borderRadius: '8px',
          cursor: 'pointer',
        }}
      >
        Login mit GitHub 🐙
      </button>
    </div>
  )
}

export default Login
