import { useState, useEffect } from 'react';
import { Table, Button, Alert } from 'react-bootstrap';
import { useSearchParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Game {
  id: number;
  status: string;
  maps: { id: number; name: string }[];
}

const GameConstructor = () => {
  const [searchParams] = useSearchParams();
  const gameId = searchParams.get('game');
  const [game, setGame] = useState<Game | null>(null);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get('/api/players/validate', { withCredentials: true })
      .then(() => {
        if (gameId) {
          axios
            .get(`/api/games/${gameId}`, { withCredentials: true })
            .then(res => setGame(res.data))
            .catch(err => {
              setError('Failed to load game. Please try again.');
              console.error('Game load error:', err.response?.data || err.message);
            });
        } else {
          axios
            .post('/api/games/add-map/0', {}, { withCredentials: true })
            .then(res => setGame(res.data))
            .catch(err => {
              setError('Failed to create new game. Please try again.');
              console.error('Game creation error:', err.response?.data || err.message);
            });
        }
      })
      .catch(err => {
        setError('Session expired. Please log in again.');
        console.error('Session validation error:', err.response?.data || err.message);
        navigate('/login');
      });
  }, [gameId, navigate]);

  if (!game && !error) return <div>Loading...</div>;

  return (
    <div>
      {error && <Alert variant="danger">{error}</Alert>}
      {game && (
        <>
          <h2>Game Constructor (Status: {game.status})</h2>
          <Table striped bordered hover>
            <thead>
              <tr>
                <th>Map</th>
                {game.status === 'draft' && <th>Actions</th>}
              </tr>
            </thead>
            <tbody>
              {game.maps.map(map => (
                <tr key={map.id}>
                  <td>{map.name}</td>
                  {game.status === 'draft' && (
                    <td>
                      <Button
                        variant="danger"
                        onClick={() =>
                          axios
                            .delete(`/api/games/${game.id}/remove-map/${map.id}`, { withCredentials: true })
                            .then(() => {
                              setGame({ ...game, maps: game.maps.filter(m => m.id !== map.id) });
                            })
                            .catch(err => {
                              setError('Failed to remove map. Please try again.');
                              console.error('Remove map error:', err.response?.data || err.message);
                            })
                        }
                      >
                        Remove
                      </Button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </Table>
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
              Form Game
            </Button>
          )}
          <Button onClick={() => navigate('/games')}>Back to Games</Button>
        </>
      )}
    </div>
  );
};

export default GameConstructor;