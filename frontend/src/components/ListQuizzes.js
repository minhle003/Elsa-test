import React from 'react';

const ListQuizzes = ({ quizzes, userId, handleCreateSession }) => {
  return (
    <div>
      <ul className="list-group">
        {quizzes?.map((quiz) => (
          <li key={quiz.ID} className="list-group-item">
            {quiz.Title}
            <button className="btn btn-success" onClick={() => handleCreateSession(quiz.ID, userId)}>
              Start Session
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ListQuizzes;