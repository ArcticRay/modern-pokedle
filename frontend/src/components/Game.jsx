import { useState, useEffect } from 'react'
import axios from 'axios'

const API = 'http://localhost:8080/api/v1'

function Game({ token, onLogout }) {
  const [guess, setGuess] = useState('')
  const [guesses, setGuesses] = useState([])
  const [status, setStatus] = useState('in_progress')
  const [guessesLeft, setGuessesLeft] = useState(6)
  const [error, setError] = useState('')

  const headers = { Authorization: `Bearer ${token}` }

  useEffect(() => {
    startGame()
  }, [])

  const startGame = async () => {
    try {
      await axios.post(`${API}/games`, {}, { headers })
      loadGame()
    } catch (e) {
      loadGame()
    }
  }

  const loadGame = async () => {
    try {
      const res = await axios.get(`${API}/games/today`, { headers })
      setGuesses(res.data.guesses || [])
      setStatus(res.data.game.Status)
      setGuessesLeft(6 - res.data.game.GuessesCount)
    } catch (e) {
      console.error(e)
    }
  }

  const submitGuess = async () => {
    if (!guess.trim()) return
    setError('')
    try {
      const res = await axios.post(
        `${API}/games/guess`,
        { pokemon_name: guess.toLowerCase() },
        { headers },
      )
      setGuesses((prev) => [
        ...prev,
        { Result: res.data.result, PokemonName: guess },
      ])
      setStatus(res.data.status)
      setGuessesLeft(res.data.guesses_left)
      setGuess('')
    } catch (e) {
      setError('Pokémon nicht gefunden!')
    }
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        background: '#1a1a2e',
        color: 'white',
        padding: '2rem',
      }}
    >
      <div style={{ maxWidth: '900px', margin: '0 auto' }}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: '2rem',
          }}
        >
          <h1 style={{ fontSize: '2.5rem' }}>Pokédle 🎮</h1>
          <button
            onClick={onLogout}
            style={{
              background: 'transparent',
              color: 'white',
              border: '1px solid white',
              padding: '0.5rem 1rem',
              borderRadius: '8px',
              cursor: 'pointer',
            }}
          >
            Logout
          </button>
        </div>

        {status === 'won' && (
          <p
            style={{
              color: '#4ade80',
              fontSize: '1.5rem',
              textAlign: 'center',
            }}
          >
            🎉 Gewonnen!
          </p>
        )}
        {status === 'lost' && (
          <p
            style={{
              color: '#f87171',
              fontSize: '1.5rem',
              textAlign: 'center',
            }}
          >
            😢 Verloren!
          </p>
        )}

        {status === 'in_progress' && (
          <div style={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
            <input
              value={guess}
              onChange={(e) => setGuess(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && submitGuess()}
              placeholder="Pokémon Name eingeben..."
              style={{
                flex: 1,
                padding: '0.75rem',
                borderRadius: '8px',
                border: 'none',
                fontSize: '1rem',
              }}
            />
            <button
              onClick={submitGuess}
              style={{
                padding: '0.75rem 1.5rem',
                background: '#facc15',
                color: 'black',
                border: 'none',
                borderRadius: '8px',
                cursor: 'pointer',
                fontWeight: 'bold',
              }}
            >
              Raten
            </button>
          </div>
        )}

        {error && <p style={{ color: '#f87171' }}>{error}</p>}
        <p style={{ marginBottom: '1rem', color: '#94a3b8' }}>
          Versuche übrig: {guessesLeft}
        </p>

        <GuessTable guesses={guesses} />
      </div>
    </div>
  )
}

function GuessTable({ guesses }) {
  if (guesses.length === 0) return null

  return (
    <div>
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: '180px repeat(7, 1fr)',
          gap: '0.5rem',
          marginBottom: '0.5rem',
          textAlign: 'center',
          color: '#94a3b8',
          fontSize: '0.85rem',
        }}
      >
        <div>Pokémon</div>
        <div>Typ 1</div>
        <div>Typ 2</div>
        <div>Habitat</div>
        <div>Farbe</div>
        <div>Entw.</div>
        <div>Höhe</div>
        <div>Gen.</div>
      </div>
      {guesses.map((g, i) => (
        <GuessRow key={i} guess={g} />
      ))}
    </div>
  )
}

function GuessRow({ guess }) {
  const r = guess.Result
  const type1 = r.Pokemon.Types[0]?.Name || '?'
  const type2 = r.Pokemon.Types[1]?.Name || 'Keine'

  return (
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: '180px repeat(7, 1fr)',
        gap: '0.5rem',
        marginBottom: '0.5rem',
      }}
    >
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '0.5rem',
          minWidth: 0,
        }}
      >
        <img
          src={r.Pokemon.SpriteURL}
          alt={r.Pokemon.Name}
          style={{ width: '48px', height: '48px', flexShrink: 0 }}
        />
        <span
          style={{
            fontSize: '0.85rem',
            textTransform: 'capitalize',
            wordBreak: 'break-word',
          }}
        >
          {r.Pokemon.Name}
        </span>
      </div>
      <Cell
        result={
          r.TypeResult === 'correct' || r.TypeResult === 'partial'
            ? 'correct'
            : 'wrong'
        }
        value={type1}
      />
      <Cell
        result={
          r.Pokemon.Types.length > 1
            ? r.TypeResult === 'correct'
              ? 'correct'
              : r.TypeResult === 'partial'
                ? 'partial'
                : 'wrong'
            : 'wrong'
        }
        value={type2}
      />
      <Cell result={r.HabitatResult} value={r.Pokemon.Habitat || '?'} />
      <Cell result={r.ColorResult} value={r.Pokemon.Color} />
      <Cell
        result={r.EvolutionResult}
        value={r.Pokemon.EvolutionStage}
        directional
      />
      <Cell
        result={r.HeightResult}
        value={`${r.Pokemon.Height * 10}cm`}
        directional
      />
      <Cell
        result={r.GenerationResult}
        value={`Gen ${r.Pokemon.Generation}`}
        directional
      />
    </div>
  )
}

function Cell({ result, value, directional }) {
  const colors = {
    correct: '#4ade80',
    partial: '#facc15',
    wrong: '#f87171',
    higher: '#f87171',
    lower: '#f87171',
  }

  const arrows = {
    higher: ' ↑',
    lower: ' ↓',
  }

  return (
    <div
      style={{
        background: colors[result] || '#374151',
        borderRadius: '8px',
        padding: '0.75rem 0.5rem',
        textAlign: 'center',
        color: 'black',
        fontWeight: 'bold',
        fontSize: '0.85rem',
        textTransform: 'capitalize',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      {value}
      {directional && arrows[result]}
    </div>
  )
}

export default Game
