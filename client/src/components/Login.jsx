import React, { useState } from 'react';
import { Form, Button, Container, Card } from 'react-bootstrap';

const Login = ({ onJoin }) => {
    const [name, setName] = useState('');
    const [room, setRoom] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        if (name && room) {
            onJoin(name, room);
        }
    };

    return (
        <Container className="d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
            <Card style={{ width: '300px' }}>
                <Card.Body>
                    <Card.Title className="text-center mb-4">Price is Right</Card.Title>
                    <Form onSubmit={handleSubmit}>
                        <Form.Group className="mb-3">
                            <Form.Label>Name</Form.Label>
                            <Form.Control 
                                type="text" 
                                placeholder="Enter name" 
                                value={name} 
                                onChange={(e) => setName(e.target.value)} 
                            />
                        </Form.Group>
                        <Form.Group className="mb-3">
                            <Form.Label>Room Code</Form.Label>
                            <Form.Control 
                                type="text" 
                                placeholder="Enter room code" 
                                value={room} 
                                onChange={(e) => setRoom(e.target.value)} 
                            />
                        </Form.Group>
                        <Button variant="primary" type="submit" className="w-100">
                            Join Game
                        </Button>
                    </Form>
                </Card.Body>
            </Card>
        </Container>
    );
};

export default Login;
