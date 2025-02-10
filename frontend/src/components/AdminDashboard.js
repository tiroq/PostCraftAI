import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

function AdminDashboard({ token }) {
  const [users, setUsers] = useState([]);
  const [expirations, setExpirations] = useState({});
  const [rateLimit, setRateLimit] = useState(1);
  const [error, setError] = useState('');

  const fetchUsers = async () => {
    try {
      const res = await axios.get(
        process.env.REACT_APP_BACKEND_URL + '/admin/list-users',
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setUsers(res.data);
      const newExps = {};
      res.data.forEach(u => {
        if (!u.allowed) {
          newExps[u.username] = 10080; // default 7 days in minutes.
        }
      });
      setExpirations(newExps);
    } catch (err) {
      setError('Failed to fetch users');
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleExpirationChange = (username, value) => {
    setExpirations({ ...expirations, [username]: Number(value) });
  };

  const handleEnable = async username => {
    try {
      const expiresIn = expirations[username] || 10080;
      await axios.post(
        process.env.REACT_APP_BACKEND_URL + '/admin/enable-user',
        { username, expires_in: expiresIn },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      fetchUsers();
    } catch (err) {
      setError('Failed to enable user');
    }
  };

  const handleUpdateExpiration = async username => {
    try {
      const expiresIn = expirations[username] || 10080;
      await axios.post(
        process.env.REACT_APP_BACKEND_URL + '/admin/update-expiration',
        { username, expires_in: expiresIn },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      fetchUsers();
    } catch (err) {
      setError('Failed to update expiration');
    }
  };

  const handleRateLimitUpdate = async () => {
    try {
      await axios.post(
        process.env.REACT_APP_BACKEND_URL + '/admin/update-rate-limit',
        { rate_limit: rateLimit },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("Rate limit updated.");
    } catch (err) {
      setError('Failed to update rate limit');
    }
  };

  const handleFetchStats = async () => {
    try {
      const res = await axios.get(
        process.env.REACT_APP_BACKEND_URL + '/admin/request-stats',
        { headers: { Authorization: `Bearer ${token}` } }
      );
      console.log("Request Stats:", res.data);
      alert("Request stats fetched. See console for details.");
    } catch (err) {
      setError('Failed to fetch request stats');
    }
  };

  return (
    <div className="container">
      <h2>Admin Dashboard</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <p>
        <Link to="/generate-post">Go to Generate Post Form</Link>
      </p>
      <div>
        <h3>Global Rate Limit (req/min)</h3>
        <input type="number" value={rateLimit} onChange={e => setRateLimit(Number(e.target.value))} />
        <button onClick={handleRateLimitUpdate}>Update Rate Limit</button>
      </div>
      <div>
        <h3>Registered Users</h3>
        <table border="1" cellPadding="5" style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th>Username</th>
              <th>Role</th>
              <th>Allowed</th>
              <th>Access Expires At</th>
              <th>Expiration (minutes)</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map(u => (
              <tr key={u.username}>
                <td>{u.username}</td>
                <td>{u.role}</td>
                <td>{u.allowed ? 'Yes' : 'No'}</td>
                <td>{u.access_expiresAt}</td>
                <td>
                  <input
                    type="number"
                    value={expirations[u.username] || 10080}
                    onChange={e => handleExpirationChange(u.username, e.target.value)}
                  />
                </td>
                <td>
                  {!u.allowed ? (
                    <button onClick={() => handleEnable(u.username)}>Enable</button>
                  ) : (
                    <button onClick={() => handleUpdateExpiration(u.username)}>Update Expiration</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        <button onClick={handleFetchStats}>Fetch Request Stats</button>
      </div>
    </div>
  );
}

export default AdminDashboard;
