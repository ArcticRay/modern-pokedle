import { useState, useEffect } from 'react'
import Game from './components/Game'
import Login from './components/Login'

function App() {
  const [token, setToken] = useState(localStorage.getItem('access_token'))

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const urlToken = params.get('token')
    if (urlToken) {
      localStorage.setItem('access_token', urlToken)
      setToken(urlToken)
      window.history.replaceState({}, '', '/')
    }
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('access_token')
    setToken(null)
  }

  return (
    <div>
      {token ? <Game token={token} onLogout={handleLogout} /> : <Login />}
    </div>
  )
}

export default App
