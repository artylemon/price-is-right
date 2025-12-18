import React from 'react';
import { Container, Card, Table, Badge, Button } from 'react-bootstrap';

const Leaderboard = ({ gameState, currentPlayerName, onReset }) => {
    const isGameOver = gameState.state === 'GAME_OVER';
    const currentItem = gameState.items[gameState.currentItem];
    const isHost = gameState.players[currentPlayerName]?.isHost;
    
    // Convert players map to array
    let players = Object.values(gameState.players);

    if (isGameOver) {
        // Sort by total score for final results
        players.sort((a, b) => a.score - b.score);
    } else {
        // Sort by difference for current round results
        players.sort((a, b) => {
            const diffA = Math.abs(a.currentGuess - currentItem.price);
            const diffB = Math.abs(b.currentGuess - currentItem.price);
            return diffA - diffB;
        });
    }

    return (
        <Container className="mt-5">
            <Card>
                <Card.Header as="h5" className={isGameOver ? "bg-success text-white" : "bg-info text-white"}>
                    {isGameOver ? "Game Over - Final Results" : "Round Results"}
                </Card.Header>
                <Card.Body>
                    {!isGameOver && (
                        <div className="text-center mb-4">
                            <h4>Item: {currentItem.name}</h4>
                            <h2 className="text-success">Actual Price: ${currentItem.price.toFixed(2)}</h2>
                            <p>Next round in: {gameState.timeLeft}s</p>
                        </div>
                    )}

                    <Table striped bordered hover>
                        <thead>
                            <tr>
                                <th>Rank</th>
                                <th>Player</th>
                                {!isGameOver && <th>Guess</th>}
                                {!isGameOver && <th>Diff</th>}
                                {isGameOver && <th>Total Score (Lower is better)</th>}
                            </tr>
                        </thead>
                        <tbody>
                            {players.map((p, index) => {
                                const diff = !isGameOver ? Math.abs(p.currentGuess - currentItem.price).toFixed(2) : 0;
                                return (
                                    <tr key={p.name} className={p.name === currentPlayerName ? "table-primary" : ""}>
                                        <td>{index + 1}</td>
                                        <td>{p.name} {p.name === currentPlayerName && '(You)'}</td>
                                        {!isGameOver && <td>${p.currentGuess.toFixed(2)}</td>}
                                        {!isGameOver && <td>${diff}</td>}
                                        {isGameOver && <td>{p.score}</td>}
                                    </tr>
                                );
                            })}
                        </tbody>
                    </Table>
                    
                    {isGameOver && (
                        <div className="text-center mt-3">
                            <h3>Winner: {players[0].name}!</h3>
                            {isHost ? (
                                <Button variant="primary" onClick={onReset}>Play Again</Button>
                            ) : (
                                <div className="alert alert-info mt-2">Waiting for host to restart...</div>
                            )}
                        </div>
                    )}
                </Card.Body>
            </Card>
        </Container>
    );
};

export default Leaderboard;
