import React, { useState, useEffect, useRef } from 'react';
import { Form } from 'react-bootstrap';
import Login from './components/Login';
import Game from './components/Game';
import Leaderboard from './components/Leaderboard';

function App() {
  const [connected, setConnected] = useState(false);
  const [gameState, setGameState] = useState(null);
  const [playerName, setPlayerName] = useState('');
  const [theme, setTheme] = useState('light');
  const ws = useRef(null);

  useEffect(() => {
    document.documentElement.setAttribute('data-bs-theme', theme);
  }, [theme]);

  const handleJoin = (name, room) => {
    setPlayerName(name);
    // Connect to WebSocket
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = process.env.NODE_ENV === 'production' ? window.location.host : 'localhost:8080';
    const socket = new WebSocket(`${protocol}//${host}/ws?name=${name}&room=${room}`);
    
    socket.onopen = () => {
      console.log('Connected to server');
      setConnected(true);
    };

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === 'STATE_UPDATE') {
        setGameState(msg.payload);
      }
    };

    socket.onclose = () => {
      console.log('Disconnected');
      setConnected(false);
      setGameState(null);
    };

    ws.current = socket;
  };

  const handleStartGame = () => {
    if (ws.current) {
      ws.current.send(JSON.stringify({ type: 'START_GAME', payload: null }));
    }
  };

  const handleResetGame = () => {
    if (ws.current) {
      ws.current.send(JSON.stringify({ type: 'RESET_GAME', payload: null }));
    }
  };

  const handleGuess = (guess) => {
    if (ws.current) {
      ws.current.send(JSON.stringify({ type: 'GUESS', payload: guess }));
    }
  };

  let content;
  if (!connected || !gameState) {
    content = <Login onJoin={handleJoin} />;
  } else if (gameState.state === 'WAITING' || gameState.state === 'GUESSING') {
    const isHost = gameState.players[playerName]?.isHost;
    content = (
      <Game 
        gameState={gameState} 
        onStart={handleStartGame} 
        onGuess={handleGuess}
        currentPlayerName={playerName}
        isHost={isHost}
      />
    );
  } else if (gameState.state === 'ROUND_RESULT' || gameState.state === 'GAME_OVER') {
    content = (
      <Leaderboard 
        gameState={gameState} 
        currentPlayerName={playerName}
        onReset={handleResetGame}
      />
    );
  } else {
    content = <div>Unknown State</div>;
  }

  return (
    <>
      <div style={{ position: 'fixed', top: '1rem', right: '1rem', zIndex: 1050 }}>
        <Form.Check 
            type="switch"
            id="theme-switch"
            label={theme === 'light' ? '🌙' : '☀️'}
            checked={theme === 'dark'}
            onChange={() => setTheme(theme === 'light' ? 'dark' : 'light')}
        />
      </div>
      {content}
    </>
  );
}

export default App;
