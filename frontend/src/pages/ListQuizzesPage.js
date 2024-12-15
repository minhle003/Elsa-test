import React, { useState } from 'react';
import { getUserByUsername, getQuizzesByUser, startSession, createSession } from '../api';
import ListQuizzes from '../components/ListQuizzes';

const ListQuizzesPage = () => {
  const [username, setUsername] = useState('');
  const [userId, setUserId] = useState('');
  const [quizzes, setQuizzes] = useState([]);
  const [error, setError] = useState('');

  const handleUsernameSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const response = await getUserByUsername(username);
      fetchQuizzes(response.data.ID);
      setUserId(response.data.ID)
      localStorage.setItem('User-ID', response.data.ID);
    } catch (error) {
      setError('User not found');
      setQuizzes([]);
    }
  };

  const fetchQuizzes = async (userId) => {
    try {
      const response = await getQuizzesByUser(userId);
      console.log(response.data)
      setQuizzes(response.data);
    } catch (error) {
      setError('Failed to fetch quizzes');
    }
  };

  const handleCreateSession = async (quizId) => {
    try {
      const session = await createSession(quizId, userId);
      console.log(session)
      window.location.href = `/session?session=${session.data.ID}`
      alert('Session started successfully');
    } catch (error) {
      alert('Error starting session');
    }
  };

  return (
    <div className="container">
      <h1>Your Quizzes</h1>
      <form onSubmit={handleUsernameSubmit}>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Enter username"
          required
        />
        <button type="submit">Get User Quizzes</button>
      </form>
      {error && <div className="alert alert-danger">{error}</div>}
      <button className="btn btn-primary" onClick={() => (window.location.href = '/quizzes/create')}>
        Create New Quiz
      </button>
      <ListQuizzes quizzes={quizzes} userId={userId} handleCreateSession={handleCreateSession} />
    </div>
  );
};

export default ListQuizzesPage;