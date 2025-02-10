import React from 'react';
import { Link } from 'react-router-dom';

function Home() {
  return (
    <div className="container">
      <h1>Welcome to PostCraft AI</h1>
      <p>
        PostCraft AI is an intelligent service that transforms long-form articles into concise, engaging posts using OpenAI's API.
      </p>
      <p>
        With robust JWT-based authentication, admin controls for user management and usage tracking, and a sleek, Material Design–inspired interface, PostCraft AI makes content transformation easy and efficient.
      </p>
      <h3>Features:</h3>
      <ul>
        <li>Convert detailed articles into short posts.</li>
        <li>Secure authentication with JWT tokens.</li>
        <li>Admin panel for managing user access, expiration, and rate limiting.</li>
        <li>Detailed request logging and statistics for usage analysis.</li>
        <li>Fully dockerized for easy deployment and scaling.</li>
      </ul>
      <p>Get started by signing up or logging in.</p>
      <div style={{ marginTop: '20px' }}>
        <Link to="/signup">
          <button>Sign Up</button>
        </Link>
        <Link to="/login" style={{ marginLeft: '10px' }}>
          <button>Login</button>
        </Link>
      </div>
    </div>
  );
}

export default Home;
