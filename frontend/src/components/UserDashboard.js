import React from 'react';
import { Link } from 'react-router-dom';

function UserDashboard({ token }) {
  const backendUrl = process.env.REACT_APP_BACKEND_URL;

  const exampleCurl = `curl -X POST ${backendUrl}/generate-post \\
  -H "Authorization: Bearer ${token}" \\
  -H "Content-Type: application/json" \\
  -d '{"article": "Your article content"}'`;

  return (
    <div className="container">
      <h2>User Dashboard</h2>
      <p>Your JWT Token:</p>
      <pre>{token}</pre>
      <p>
        <Link to="/generate-post">Go to Generate Post Form</Link>
      </p>
      <p>Example cURL command to generate a post:</p>
      <pre>{exampleCurl}</pre>
    </div>
  );
}

export default UserDashboard;
