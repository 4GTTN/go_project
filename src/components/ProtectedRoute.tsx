import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { Navigate } from 'react-router-dom';
import axios from 'axios';
import { RootState } from '../store';

interface ProtectedRouteProps {
  children: JSX.Element;
}

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const user = useSelector((state: RootState) => state.auth.user);
  const [isValidating, setIsValidating] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    if (!user) {
      setIsValidating(false);
      return;
    }
    axios
      .get('/api/players/validate', { withCredentials: true })
      .then(() => setIsAuthenticated(true))
      .catch(err => {
        console.error('Session validation failed:', err.response?.data || err.message);
        setIsAuthenticated(false);
      })
      .finally(() => setIsValidating(false));
  }, [user]);

  if (isValidating) {
    return <div>Loading...</div>;
  }

  if (!user || !isAuthenticated) {
    console.log('ProtectedRoute: Redirecting to /login');
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default ProtectedRoute;