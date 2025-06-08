import { useState, useEffect } from 'react';
import { Card, Button, Form } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Map {
  id: number;
  name: string;
  price: number;
  image: string;
}

const MapList = () => {
  const [maps, setMaps] = useState<Map[]>([]);
  const [filter, setFilter] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get(`/api/maps?filter=${filter}`, { withCredentials: true })
      .then(res => setMaps(res.data))
      .catch(err => {
        console.error('MapList error:', err.response?.data || err.message);
        setMaps([]);
      });
  }, [filter]);

  return (
    <div>
      <Form>
        <Form.Control
          type="text"
          value={filter}
          onChange={e => setFilter(e.target.value)}
          placeholder="Filter by name or price"
        />
      </Form>
      <div className="cards">
        {maps.map(map => (
          <Card key={map.id} style={{ width: '18rem' }}>
            <Card.Img variant="top" src={map.image || '/gs.jpg'} />
            <Card.Body>
              <Card.Title>{map.name}</Card.Title>
              <Card.Text>Price: ${map.price}</Card.Text>
              <Button onClick={() => navigate(`/map/${map.id}`)}>Details</Button>
              <Button
                onClick={() =>
                  axios
                    .post(`/api/games/add-map/${map.id}`, {}, { withCredentials: true })
                    .then(() => navigate('/game-constructor'))
                    .catch(err => console.error('Add map error:', err.response?.data || err.message))
                }
              >
                Add to Game
              </Button>
            </Card.Body>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default MapList;