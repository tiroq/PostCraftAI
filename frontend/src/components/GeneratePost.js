import React, { useState } from 'react';
import axios from 'axios';

function GeneratePost({ token }) {
  const [article, setArticle] = useState('');
  const [post, setPost] = useState('');
  const [error, setError] = useState('');

  const handleGenerate = async e => {
    e.preventDefault();
    try {
      const res = await axios.post(
        process.env.REACT_APP_BACKEND_URL + '/generate-post',
        { article },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setPost(res.data.post);
    } catch (err) {
      setError('Failed to generate post');
    }
  };

  return (
    <div className="container">
      <h2>Generate Post</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <form onSubmit={handleGenerate}>
        <textarea
          placeholder="Enter article"
          value={article}
          onChange={e => setArticle(e.target.value)}
          rows="10"
          required
        />
        <br />
        <button type="submit">Generate Post</button>
      </form>
      {post && (
        <div>
          <h3>Generated Post:</h3>
          <pre>{post}</pre>
        </div>
      )}
    </div>
  );
}

export default GeneratePost;
