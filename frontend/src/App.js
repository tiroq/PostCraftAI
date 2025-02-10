import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Home from './components/Home';
import Login from './components/Login';
import Signup from './components/Signup';
import GeneratePost from './components/GeneratePost';
import AdminDashboard from './components/AdminDashboard';
import UserDashboard from './components/UserDashboard';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [role, setRole] = useState(localStorage.getItem('role') || '');

  const handleLogin = (token, role) => {
    setToken(token);
    setRole(role);
    localStorage.setItem('token', token);
    localStorage.setItem('role', role);
  };

  const handleLogout = () => {
    setToken('');
    setRole('');
    localStorage.removeItem('token');
    localStorage.removeItem('role');
  };

  return (
    <Router>
      <div>
        {token && <button onClick={handleLogout}>Logout</button>}
        <Routes>
          <Route
            path="/"
            element={
              !token ? <Home /> : (role === 'admin' ? <Navigate to="/admin" /> : <Navigate to="/dashboard" />)
            }
          />
          <Route path="/login" element={<Login onLogin={handleLogin} />} />
          <Route path="/signup" element={<Signup onLogin={handleLogin} />} />
          <Route path="/generate-post" element={token ? <GeneratePost token={token} /> : <Navigate to="/login" />} />
          <Route path="/dashboard" element={token ? <UserDashboard token={token} /> : <Navigate to="/login" />} />
          <Route path="/admin" element={token && role === 'admin' ? <AdminDashboard token={token} /> : <Navigate to="/login" />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
