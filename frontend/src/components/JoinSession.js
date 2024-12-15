import React, { useState } from 'react';
import { joinSession } from '../api';

const JoinSession = () => {
  const [sessionId, setSessionId] = useState('');
  const [name, setName] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await joinSession(sessionId, name);
      localStorage.setItem('User-ID', response.data.participantId);
      alert('Joined session successfully');
      window.location.href = `/session?session=${sessionId}`;
    } catch (error) {
      alert('Error joining session');
    }
  };

  return (
    <div className="container">
      <h1>Join Session</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={sessionId}
          onChange={(e) => setSessionId(e.target.value)}
          placeholder="Enter session ID"
          required
        />
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter your name"
          required
        />
        <button type="submit">Join Session</button>
      </form>
    </div>
  );
};

export default JoinSession;