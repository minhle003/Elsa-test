import React, { useState } from 'react';
import { createQuiz } from '../api';

const CreateQuiz = ({ userId }) => {
  const [quizData, setQuizData] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const quiz = JSON.parse(quizData);
      await createQuiz({ ...quiz, createdBy: userId });
      alert('Quiz created successfully');
    } catch (error) {
      alert('Error creating quiz');
    }
  };

  return (
    <div className="container">
      <h1>Create Quiz</h1>
      <form onSubmit={handleSubmit}>
        <textarea
          value={quizData}
          onChange={(e) => setQuizData(e.target.value)}
          placeholder="Enter quiz data in JSON format"
          rows="10"
          required
        />
        <button type="submit">Create Quiz</button>
      </form>
    </div>
  );
};

export default CreateQuiz;