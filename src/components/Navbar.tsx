import { Navbar as BSNavbar, Nav, Button } from 'react-bootstrap';
import { useSelector, useDispatch } from 'react-redux';
import { Link, useNavigate } from 'react-router-dom';
import { logout } from '../store';
import axios from 'axios';

const Navbar = () => {
  const user = useSelector((state: any) => state.auth.user);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    axios
      .post('/api/players/logout', {}, { withCredentials: true })
      .then(() => {
        dispatch(logout());
        navigate('/login');
      })
      .catch(err => console.error('Logout failed:', err.response?.data || err.message));
  };

  return (
    <BSNavbar bg="light" expand="lg">
      <BSNavbar.Brand as={Link} to="/">Game Matchmaking</BSNavbar.Brand>
      <BSNavbar.Toggle />
      <BSNavbar.Collapse>
        <Nav className="me-auto">
          <Nav.Link as={Link} to="/">Maps</Nav.Link>
          {user && <Nav.Link as={Link} to="/games">My Games</Nav.Link>}
          {user && <Nav.Link as={Link} to="/game-constructor">Constructor</Nav.Link>}
        </Nav>
        <Nav>
          {user ? (
            <>
              <Nav.Link disabled>{user.username}</Nav.Link>
              <Button onClick={handleLogout}>Logout</Button>
            </>
          ) : (
            <>
              <Nav.Link as={Link} to="/login">Login</Nav.Link>
              <Nav.Link as={Link} to="/register">Register</Nav.Link>
            </>
          )}
        </Nav>
      </BSNavbar.Collapse>
    </BSNavbar>
  );
};

export default Navbar;