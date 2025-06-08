import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './store';
import Navbar from './components/Navbar';
import MapList from './components/MapList';
import MapDetail from './components/MapDetail';
import Login from './components/Login';
import Register from './components/Register';
import GameList from './components/GameList';
import GameConstructor from './components/GameConstructor';
import ManagerDashboard from './components/ManagerDashboard';
import Breadcrumbs from './components/Breadcrumbs';
import ProtectedRoute from './components/ProtectedRoute';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <div className="min-h-screen bg-gray-100">
          <Navbar />
          <div className="container mx-auto px-4 py-6">
            <Breadcrumbs />
            <Routes>
              <Route path="/" element={<MapList />} />
              <Route path="/map/:id" element={<MapDetail />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route
                path="/games"
                element={
                  <ProtectedRoute>
                    <GameList />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/game-constructor"
                element={
                  <ProtectedRoute>
                    <GameConstructor />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/manager"
                element={
                  <ProtectedRoute>
                    <ManagerDashboard />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </div>
        </div>
      </Router>
    </Provider>
  );
}

export default App;