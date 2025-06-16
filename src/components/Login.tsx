import { useState } from 'react';
import { Form, Button, Alert } from 'react-bootstrap';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { login } from '../store';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    try {
      const res = await axios.post('/api/players/login', { username, password }, { withCredentials: true });
      dispatch(login({ user: { username: res.data.username }, token: 'session' }));
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed. Please try again.');
      console.error('Login error:', err.response?.data || err.message);
    }
  };

  return (
    <Form>
      {error && <Alert variant="danger">{error}</Alert>}
      <Form.Group className="mb-3">
        <Form.Label>Username</Form.Label>
        <Form.Control
          type="text"
          value={username}
          onChange={e => setUsername(e.target.value)}
          placeholder="Enter username"
          required
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Password</Form.Label>
        <Form.Control
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="Enter password"
          required
        />
      </Form.Group>
      <Button variant="primary" onClick={handleSubmit}>Login</Button>
    </Form>
  );
};

export default Login;