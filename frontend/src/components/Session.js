import React, { useEffect, useState, useRef } from 'react';
import { startSession, updateScore, changeQuestion } from '../api';
import { WebSockethost } from '../socket';


const SessionComponent = ({ sessionId, userId }) => {
    const [session, setSessionData] = useState(null);
    const [error, setError] = useState(null);
    const [socket, setSocket] = useState(null);
    const [currentPage, setCurrentPage] = useState(1);
    const [answer, setAnswer] = useState('');
    const [timeLeft, setTimeLeft] = useState(null); 
    const [isTimerRunning, setIsTimerRunning] = useState(false);
    const [isSubmit, setIsSubmit] = useState(false)
    const [score, setScore] = useState(0)
    const [questionIndex, setQuestionIndex] = useState(-1)

    const questionIndexRef = useRef(-2)
    const currentPageRef = useRef(1)

    useEffect(() => {
        questionIndexRef.current = questionIndex
        currentPageRef.current = currentPage
    }, [questionIndex, currentPage])


    useEffect(() => {
        const ws = new WebSocket(WebSockethost);
        ws.onopen = () => {
            console.log('Connected to WebSocket server');
            const message = { sessionId: sessionId, userId: userId };
            ws.send(JSON.stringify(message));
        };
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log(currentPageRef.current, data, questionIndexRef)
            if (data?.Quiz?.Questions) {
                if (data?.CurrentQuestionIndex != questionIndexRef.current) {
                    setIsTimerRunning(true);
                    setTimeLeft(data?.Quiz?.Questions[data?.CurrentQuestionIndex].Time)
                }
            }
            setSessionData(data);
            setQuestionIndex(data.CurrentQuestionIndex)
            if (data.Status === 'started' && currentPageRef.current == 1) {
                setCurrentPage(2);
            }
            if (error) {
                setError(null)
            }
        };
        ws.onclose = (event) => {
            console.log('Disconnected from WebSocket server:', event);
        };

        ws.onerror = (err) => {
            if (err != error) {
                setError(err)
            }
            console.error('WebSocket error:', err);
        };

        setSocket(ws);

        return () => {
            if (socket) {
                socket.close();
            }
        };

    }, [sessionId, userId]);

    const handleStartSession = async () => {
        try {
            await startSession(sessionId);
        } catch (error) {
            setError('Failed to start session');
        }
    };

    const handleChangeQuestion = async (nextQuestionIndex) => {
        await changeQuestion(sessionId, nextQuestionIndex)
    }

    useEffect(() => {
        let timerInterval;

        if (isTimerRunning && timeLeft > 0) {
            timerInterval = setInterval(() => {
                setTimeLeft((prevTime) => prevTime - 1);
            }, 1000);
        } else if (timeLeft === 0) {
            clearInterval(timerInterval)
            if (currentPage == 2){
                if (userId == session.CreatedBy && session?.CurrentQuestionIndex < session.Quiz?.Questions.length - 1) {
                    handleChangeQuestion(questionIndex + 1)
                }
                setCurrentPage(3)
                setTimeLeft(5)
                setScore(0)
                setIsSubmit(false)
            } else if (currentPage == 3) {
                if (session?.Quiz?.Questions.length - 1 != session?.CurrentQuestionIndex) {
                    setCurrentPage(2)
                    setTimeLeft(session?.Quiz?.Questions[session?.CurrentQuestionIndex].Time)
                } else {
                    setIsTimerRunning(false)
                }
            }
        }


        return () => {
            clearInterval(timerInterval);
        };
    }, [timeLeft, isTimerRunning]);

    const handleAnswerSubmit = async () => {
        setIsSubmit(true)
        if (answer) {
            const correct = session.Quiz.Questions[session.CurrentQuestionIndex].Answer == answer;
            if (correct) {
                let newScore = calculateScore() + session.Participants[userId].Score
                try {
                    await updateScore(sessionId, userId, newScore);
                } catch (error) {
                    alert("Failed to submit answer", error)

                }
            }
        }
        setAnswer('');
    };

    const calculateScore = () => {
        const remainTime = timeLeft / session.Quiz?.Questions[session?.CurrentQuestionIndex].Time
        const scoreEarn = Math.floor(session.Quiz.Questions[session.CurrentQuestionIndex].Score * remainTime)
        setScore(scoreEarn)
        return scoreEarn
    }

    return (
        <div className="container">
            {currentPage === 1 && (
                <div>
                    <h3>Participants</h3>
                    {userId == session?.CreatedBy && (
                        <button className="btn btn-primary" onClick={handleStartSession}>Start Session</button>
                    )}
                    {session?.Participants && (
                        <ul>
                            {Object.keys(session?.Participants)?.map((participant) => (
                                <li key={participant}>{session?.Participants[participant].Name}</li>
                            ))}
                        </ul>
                    )}
                </div>
            )}
            {currentPage === 2 && (
                <div>
                    <h1>Quesion {session?.CurrentQuestionIndex + 1}</h1>
                    <h3>{session?.Quiz.Questions[session?.CurrentQuestionIndex].Title}</h3>
                    <h4>{session?.Quiz.Questions[session?.CurrentQuestionIndex].Description}</h4>
                    {session?.Quiz.Questions[session?.CurrentQuestionIndex].Type === "multiple_choice" ? (
                        <ul>
                            {session?.Quiz?.Questions[session?.CurrentQuestionIndex]?.Choices?.map((choice, index) => (
                                <li key={index} onClick={() => setAnswer(choice.Description)} style={{ backgroundColor: answer === choice.Description ? 'lightblue' : 'white' }}>
                                    {choice.Description}
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
                    {userId != session?.CreatedBy && !isSubmit && (
                        <button onClick={handleAnswerSubmit}>Submit Answer</button>
                    )}
                    {isSubmit && (
                        <div>
                            {score == 0 ? (
                                <p style={{color: "red"}}>Incorrect, the answer is {session.Quiz.Questions[session.CurrentQuestionIndex].Answer}</p>
                            ) : (
                                <p style={{color: "green"}}>Correct, you earn {score} points </p>
                            )}
                        </div>
                    )}
                </div>
            )}
            {currentPage === 3 && (
                <div>
                    <h3>LeaderBoards</h3>
                    <h4></h4>
                    <ul>
                        {Object.entries(session?.Participants).sort((a, b) => b[1].Score - a[1].Score).map((participant, index) => (
                            <li key={participant[0]}>
                               {index+1}. {participant[1].Name}: {participant[1].Score}
                            </li>
                        ))}
                    </ul>
                </div>
            )}
            {(currentPage == 2 || currentPage == 3) && isTimerRunning && timeLeft !== null && (
                <div>Time Left: {timeLeft}s</div>
            )}
        </div>
    );
};

export default SessionComponent;