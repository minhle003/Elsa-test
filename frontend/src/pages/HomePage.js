import React from 'react';
import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <div className="container">
      <h1>Welcome to the Quiz App</h1>
      <ul className="list-group">
        <li className="list-group-item">
          <Link to="/user/create" className="btn btn-primary">Create User</Link>
        </li>
        <li className="list-group-item">
          <Link to="/quizzes" className="btn btn-primary">List Quizzes</Link>
        </li>
        <li className="list-group-item">
          <Link to="/quizzes/create" className="btn btn-primary">Create Quiz</Link>
        </li>
        <li className="list-group-item">
          <Link to="/session/join" className="btn btn-primary">Join Session</Link>
        </li>
      </ul>
    </div>
  );
};

export default HomePage;