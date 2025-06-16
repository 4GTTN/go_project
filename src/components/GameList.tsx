import { useState, useEffect } from 'react';
import { Table, Button, Alert } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Game {
  id: number;
  status: string;
  maps: { id: number; name: string }[];
}

const GameList = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get('/api/players/validate', { withCredentials: true })
      .then(() => {
        axios
          .get('/api/games', { withCredentials: true })
          .then(res => setGames(res.data))
          .catch(err => {
            setError('Failed to load games. Please try again.');
            console.error('GameList error:', err.response?.data || err.message);
          });
      })
      .catch(err => {
        setError('Session expired. Please log in again.');
        console.error('Session validation error:', err.response?.data || err.message);
        navigate('/login');
      });
  }, [navigate]);

  return (
    <div>
      {error && <Alert variant="danger">{error}</Alert>}
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
                <Button onClick={() => navigate(`/game-constructor?game=${game.id}`)}>View</Button>
                {game.status === 'draft' && (
                  <Button
                    onClick={() =>
                      axios
                        .put(`/api/games/form/${game.id}`, {}, { withCredentials: true })
                        .then(() => navigate('/games'))
                        .catch(err => {
                          setError('Failed to form game. Please try again.');
                          console.error('Form game error:', err.response?.data || err.message);
                        })
                    }
                  >
                    Form
                  </Button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default GameList;