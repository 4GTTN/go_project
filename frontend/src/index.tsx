import React from 'react';
import ReactDOM from 'react-dom/client';
import axios from 'axios';
import App from './App';
import 'bootstrap/dist/css/bootstrap.min.css';
import './stules.css';

axios.defaults.withCredentials = true;

axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401 || error.response?.status === 403) {
      console.error('Unauthorized request:', error.response?.data || error.message);
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

const rootElement = document.getElementById('root')!;
const root = ReactDOM.createRoot(rootElement);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);