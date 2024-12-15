import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
});

// Add a request interceptor to include the User-ID header
api.interceptors.request.use((config) => {
  const userId = localStorage.getItem('User-ID');
  if (userId) {
    config.headers['User-ID'] = userId;
  }
  return config;
});

export const createUser = (name) => api.post('/user', { name });
export const getUserByUsername = (username) => api.get(`/user/${username}`);
export const getQuizzesByUser = () => api.get(`/quiz/quizzes`);
export const createQuiz = (quizData) => api.post('/quiz', quizData);
export const startSession = (sessionId) => api.patch('/session/start', { sessionId });
export const createSession = (quizId) => api.post('/session', { quizId });
export const getSession = (sessionId) => api.get(`/session/${sessionId}`)
export const joinSession = (sessionId, name) => api.patch('/session/join', { sessionId, name });
export const changeQuestion = (sessionId, currentQuestionIndex) =>
  api.patch('/session/change_question', { sessionId, currentQuestionIndex });
export const updateScore = (sessionId, score) =>
  api.patch('/session/participant/update_score', { sessionId, score });

export default api;