import React, { useState } from 'react';
import { Container, Card, Button, Form, Row, Col, ListGroup, Badge } from 'react-bootstrap';

const Game = ({ gameState, onStart, onGuess, currentPlayerName, isHost }) => {
    const [guess, setGuess] = useState('');

    const handleGuess = (e) => {
        e.preventDefault();
        if (guess) {
            onGuess(parseFloat(guess));
        }
    };

    if (gameState.state === 'WAITING') {
        return (
            <Container className="mt-5">
                <Card>
                    <Card.Header as="h5">Lobby - Room: {gameState.id}</Card.Header>
                    <Card.Body>
                        <Card.Title>Players joined:</Card.Title>
                        <ListGroup variant="flush" className="mb-3">
                            {Object.values(gameState.players).map(p => (
                                <ListGroup.Item key={p.name}>
                                    {p.name} {p.name === currentPlayerName && '(You)'} {p.isHost && <Badge bg="warning" text="dark">Host</Badge>}
                                </ListGroup.Item>
                            ))}
                        </ListGroup>
                        {isHost ? (
                            <Button variant="success" onClick={onStart}>Start Game</Button>
                        ) : (
                            <div className="alert alert-info">Waiting for host to start the game...</div>
                        )}
                    </Card.Body>
                </Card>
            </Container>
        );
    }

    if (gameState.state === 'GUESSING') {
        const currentItem = gameState.items[gameState.currentItem];
        const myPlayer = gameState.players[currentPlayerName];
        const hasGuessed = myPlayer?.hasGuessed;

        return (
            <Container className="mt-5">
                <Row className="justify-content-center">
                    <Col md={8}>
                        <Card className="text-center">
                            <Card.Header>
                                Time Left: <Badge bg={gameState.timeLeft < 10 ? 'danger' : 'primary'}>{gameState.timeLeft}s</Badge>
                            </Card.Header>
                            <Card.Body>
                                <Card.Title>{currentItem.name}</Card.Title>
                                <div className="mb-3">
                                    <img 
                                        src={currentItem.imageUrl} 
                                        alt={currentItem.name} 
                                        style={{ maxHeight: '300px', maxWidth: '100%', objectFit: 'contain' }} 
                                    />
                                </div>
                                
                                {hasGuessed ? (
                                    <div className="alert alert-info">
                                        You guessed ${myPlayer.currentGuess}. Waiting for others...
                                    </div>
                                ) : (
                                    <Form onSubmit={handleGuess}>
                                        <Form.Group className="mb-3">
                                            <Form.Label>Enter your guess ($)</Form.Label>
                                            <Form.Control 
                                                type="number" 
                                                step="0.01" 
                                                value={guess} 
                                                onChange={(e) => setGuess(e.target.value)} 
                                                autoFocus
                                            />
                                        </Form.Group>
                                        <Button variant="primary" type="submit">Submit Guess</Button>
                                    </Form>
                                )}
                            </Card.Body>
                            <Card.Footer className="text-muted">
                                Players guessed: {Object.values(gameState.players).filter(p => p.hasGuessed).length} / {Object.values(gameState.players).length}
                            </Card.Footer>
                        </Card>
                    </Col>
                </Row>
            </Container>
        );
    }

    return null;
};

export default Game;
