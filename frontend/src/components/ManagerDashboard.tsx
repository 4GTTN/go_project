import { useState, useEffect } from 'react';
import { Table, Button, Alert } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

interface Game {
  id: number;
  status: string;
  maps: { id: number; name: string }[];
}

const ManagerDashboard = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [games, setGames] = useState<Game[]>([]);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get('/api/manager/users', { withCredentials: true })
      .then(res => setUsers(res.data))
      .catch(err => {
        setError('Failed to load users. Please try again.');
        console.error('Users load error:', err.response?.data || err.message);
      });

    axios
      .get('/api/manager/games', { withCredentials: true })
      .then(res => setGames(res.data))
      .catch(err => {
        setError('Failed to load games. Please try again.');
        console.error('Games load error:', err.response?.data || err.message);
      });
  }, []);

  return (
    <div>
      <h2>Manager Dashboard</h2>
      {error && <Alert variant="danger">{error}</Alert>}
      
      <h3>Users</h3>
      <Table striped bordered hover>
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Email</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          {users.map(user => (
            <tr key={user.id}>
              <td>{user.id}</td>
              <td>{user.username}</td>
              <td>{user.email}</td>
              <td>{user.role}</td>
            </tr>
          ))}
        </tbody>
      </Table>

      <h3>All Games</h3>
      <Table striped bordered hover>
        <thead>
          <tr>
            <th>ID</th>
            <th>Status</th>
            <th>Maps</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {games.map(game => (
            <tr key={game.id}>
              <td>{game.id}</td>
              <td>{game.status}</td>
              <td>{game.maps.map(m => m.name).join(', ')}</td>
              <td>
                <Button
                  onClick={() =>
                    axios
                      .put(`/api/games/complete/${game.id}`, {}, { withCredentials: true })
                      .then(() => navigate('/manager'))
                      .catch(err => {
                        setError('Failed to complete game. Please try again.');
                        console.error('Complete game error:', err.response?.data || err.message);
                      })
                  }
                  disabled={game.status !== 'formed'}
                >
                  Complete
                </Button>
                <Button
                  variant="danger"
                  onClick={() =>
                    axios
                      .put(`/api/games/reject/${game.id}`, {}, { withCredentials: true })
                      .then(() => navigate('/manager'))
                      .catch(err => {
                        setError('Failed to reject game. Please try again.');
                        console.error('Reject game error:', err.response?.data || err.message);
                      })
                  }
                  disabled={game.status !== 'formed'}
                >
                  Reject
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
      <Button onClick={() => navigate('/')}>Back to Home</Button>
    </div>
  );
};

export default ManagerDashboard;