import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './pages/HomePage';
import CreateUserPage from './pages/CreateUserPage';
import ListQuizzesPage from './pages/ListQuizzesPage';
import CreateQuizPage from './pages/CreateQuizPage';
import JoinSessionPage from './pages/JoinSessionPage';
import SessionPage from './pages/SessionPage';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/user/create" element={<CreateUserPage />} />
        <Route path="/quizzes" element={<ListQuizzesPage userId="user-id-placeholder" />} />
        <Route path="/quizzes/create" element={<CreateQuizPage userId="user-id-placeholder" />} />
        <Route path="/session/join" element={<JoinSessionPage />} />
        <Route path="/session" element={<SessionPage />} />
      </Routes>
    </Router>
  );
};

export default App;