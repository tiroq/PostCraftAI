import React, { useState } from 'react';
import axios from 'axios';

function decodeJWT(token) {
  try {
    const payload = token.split('.')[1];
    return JSON.parse(atob(payload));
  } catch (e) {
    return {};
  }
}

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      const res = await axios.post(
        process.env.REACT_APP_BACKEND_URL + '/login',
        { username, password }
      );
      const token = res.data.token;
      const payload = decodeJWT(token);
      onLogin(token, payload.role || 'user');
    } catch (err) {
      setError('Login failed');
    }
  };

  return (
    <div className="container">
      <h2>Login</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} required />
        <br />
        <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} required />
        <br />
        <button type="submit">Login</button>
      </form>
      <p>
        Or <a href="/signup">Sign Up</a>
      </p>
    </div>
  );
}

export default Login;
