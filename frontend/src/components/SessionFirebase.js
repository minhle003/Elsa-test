// demo for Firebase real time data update
import React, { useEffect, useState } from 'react';
import { doc, onSnapshot, updateDoc } from 'firebase/firestore';
import { changeQuestion, updateScore } from '../api';
import db from '../firebase';


const Session = ({ sessionId, userId }) => {
  const [session, setSession] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [answer, setAnswer] = useState('');
  const [timer, setTimer] = useState(0);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!sessionId || !userId) {
      setError('Session ID or User ID is missing');
      return;
    }

    try {
        const sessionRef = doc(db, 'sessions', sessionId);

        const unsubscribe = onSnapshot(sessionRef, (doc) => {
          if (doc.exists()) {
            const sessionData = doc.data();
            setSession(sessionData);
            if (sessionData.Status === 'started') {
              setCurrentPage(2);
              setTimer(30);
            }
          } else {
            setError('Session not found');
          }
        });
        return () => {
            unsubscribe();
        };
    } catch (error) {
        setError(error)
    }
  }, [sessionId, userId]);

  useEffect(() => {
    if (timer > 0) {
      const interval = setInterval(() => {
        setTimer(timer - 1);
      }, 1000);
      return () => clearInterval(interval);
    } else if (timer === 0 && currentPage === 2) {
      // Time's up, show answer and switch to page 3
      setCurrentPage(3);
      setTimeout(() => {
        if (session.CurrentQuestionIndex < session.Quiz.Questions.length - 1) {
          changeQuestion(sessionId, session.CurrentQuestionIndex + 1);
          setCurrentPage(2);
          setTimer(30); // Reset timer for next question
        }
      }, 5000);
    }
  }, [timer, currentPage, session, sessionId]);

  const handleAnswerSubmit = async () => {
    if (answer) {
      const correct = session.Quiz.Questions[session.CurrentQuestionIndex].Choices.includes(answer);
      if (correct) {
        await updateScore(sessionId, session.Participants[userId].Score + 1);
      }
      setAnswer('');
    }
  };

  const handleStartSession = async () => {
    const sessionRef = doc(db, 'sessions', sessionId);
    try {
      await updateDoc(sessionRef, { Status: 'started' });
    } catch (error) {
      setError('Failed to start session');
    }
  };

  if (error) {
    return <div className="alert alert-danger">{error}</div>;
  }

  if (!session) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container">

      {currentPage === 1 && (
        <div>
          <h3>Participants</h3>
          {userId == session.CreatedBy && (
            <button className="btn btn-primary" onClick={handleStartSession}>Start Session</button>
          )}
          <ul>
            {Object.keys(session.Participants).map((participant) => (
              <li key={participant}>{participant}</li>
            ))}
          </ul>
        </div>
      )}
      {currentPage === 2 && (
        <div>
          <h3>{session.Quiz.Questions[session.CurrentQuestionIndex].Title}</h3>
          {session.Quiz.Questions[session.CurrentQuestionIndex].Type === 'multiple_choice' ? (
            <ul>
              {session.Quiz.Questions[session.CurrentQuestionIndex].Choices.map((choice) => (
                <li key={choice} onClick={() => setAnswer(choice)}>
                  {choice}
                </li>
              ))}
            </ul>
          ) : (
            <input
              type="text"
              value={answer}
              onChange={(e) => setAnswer(e.target.value)}
              placeholder="Type your answer"
            />
          )}
          <button onClick={handleAnswerSubmit}>Submit Answer</button>
          <div className="timer">Time left: {timer} seconds</div>
        </div>
      )}
      {currentPage === 3 && (
        <div>
          <h3>Scores</h3>
          <ul>
            {Object.keys(session.Participants).map((participant) => (
              <li key={participant}>
                {participant}: {session.Participants[participant].score}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default Session;