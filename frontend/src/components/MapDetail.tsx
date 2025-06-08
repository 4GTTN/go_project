import { useState, useEffect } from 'react';
import { Card, Button } from 'react-bootstrap';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Map {
  id: number;
  name: string;
  price: number;
  image: string;
  description: string;
  location: string;
}
interface Props {

}
const MapDetail = ({}:Props) => {
  const { id } = useParams();
  const [map, setMap] = useState<Map | null>(null);
  const navigate = useNavigate();
  // TODO: изменить вызов карточек, в каждой карточке не должно быть запроса на получение этой карточки, данные должны браться из полученого массива данных
  useEffect(() => {
    axios
      .get(`/api/maps/${id}`, { withCredentials: true })
      .then(res => setMap(res.data))
      .catch(() => setMap(null));
  }, [id]);

  if (!map) return <div>Map not found</div>;

  return (
    <Card>
      <Card.Img variant="top" src={map.image || '/default.jpg'} />
      <Card.Body>
        <Card.Title>{map.name}</Card.Title>
        <Card.Text>Price: ${map.price}</Card.Text>
        <Card.Text>Description: {map.description}</Card.Text>
        <Card.Text>Location: {map.location}</Card.Text>
        <Button onClick={() => navigate('/')}>Back to List</Button>
      </Card.Body>
    </Card>
  );
};

export default MapDetail;