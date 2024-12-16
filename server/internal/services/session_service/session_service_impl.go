package session_service

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minhle003/Elsa-test/internal/services/quiz_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/minhle003/Elsa-test/pkg/logging"
)

const (
	statusCreated = "created"
	statusStarted = "started"
	statusEnded   = "ended"
)

type sessionServiceImpl struct {
	client *firestore.Client
	logger logging.Logger
	ctx    context.Context
}

func NewSessionService(client *firestore.Client, logger logging.Logger, ctx context.Context) SessionService {
	return &sessionServiceImpl{
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (s *sessionServiceImpl) CreateSession(quizId string, userId string) (Session, error) {
	err := s.checkUser(userId)
	if err != nil {
		return Session{}, err
	}

	if quizId == "" {
		return Session{}, fmt.Errorf("missing quiz id")
	}

	quizDocRef := s.client.Collection("quizzes").Doc(quizId)
	quizDoc, err := quizDocRef.Get(s.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return Session{}, fmt.Errorf("quiz not found")
		}
		s.logger.Error("fail to get quiz data", err)
		return Session{}, fmt.Errorf("fail to get quiz data")
	}

	var quiz quiz_service.Quiz
	if err := quizDoc.DataTo(&quiz); err != nil {
		s.logger.Error("failed to parse quiz data", err)
		return Session{}, fmt.Errorf("failed to parse quiz data")
	}

	session := Session{
		CreatedBy:            userId,
		Quiz:                 quiz,
		CurrentQuestionIndex: -1,
		CreatedTime:          time.Now().Unix(),
		Status:               statusCreated,
		Participants:         map[string]Participant{},
	}

	sessionDocRef, _, err := s.client.Collection("sessions").Add(s.ctx, session)
	if err != nil {
		s.logger.Error("fail to create session", err)
		return Session{}, fmt.Errorf("fail to create session")
	}

	session.ID = sessionDocRef.ID

	return session, nil
}

func (s *sessionServiceImpl) JoinSession(name string, sessionId string) (Session, string, error) {
	sessionDocRef, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, "", err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, "", err
	}

	if session.Status == statusEnded || session.Status == statusStarted {
		s.logger.Error("attempt to join ended session")
		return Session{}, "", fmt.Errorf("cannot join this quiz")
	}

	var userId string
	for {
		userId = uuid.New().String()
		if _, ok := session.Participants[userId]; !ok {
			break
		}
	}
	session.Participants[userId] = Participant{
		Nane:  name,
		Score: 0,
	}

	_, err = sessionDocRef.Set(s.ctx, session)
	if err != nil {
		s.logger.Error("fail to join session", err)
		return Session{}, "", fmt.Errorf("fail to join session")
	}

	return session, userId, nil
}

func (s *sessionServiceImpl) GetSession(sessionId string, userId string) (Session, error) {
	_, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, err
	}

	_, ok := session.Participants[userId]

	if userId != session.CreatedBy && !ok {
		return Session{}, fmt.Errorf("user not allow to get this session info")
	}

	if session.Status == statusCreated {
		session.Quiz = quiz_service.Quiz{}
	}
	return session, nil
}

func (s *sessionServiceImpl) StartSession(sessionId string, userId string) (Session, error) {
	sessionDocRef, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, err
	}

	if userId != session.CreatedBy {
		return Session{}, fmt.Errorf("user not allowed to update session")
	}

	session.Status = statusStarted
	session.CurrentQuestionIndex = 0

	_, err = sessionDocRef.Set(s.ctx, session)

	if err != nil {
		s.logger.Error("failed to start session", err)
		return Session{}, fmt.Errorf("failed to start session")
	}

	return session, nil
}

func (s *sessionServiceImpl) ChangeQuestion(sessionId string, questionIndex int, userId string) (Session, error) {
	sessionDocRef, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, err
	}
	if session.Status == statusEnded {
		s.logger.Error("attempt to join ended session")
		return Session{}, fmt.Errorf("cannot update this session")
	}
	if userId != session.CreatedBy {
		return Session{}, fmt.Errorf("user not allowed to update session")
	}

	if questionIndex > len(session.Quiz.Questions)-1 || questionIndex < 0 {
		return Session{}, fmt.Errorf("invalid quesion")
	}

	session.CurrentQuestionIndex = questionIndex

	_, err = sessionDocRef.Set(s.ctx, session)

	if err != nil {
		s.logger.Error("failed to change session's question", err)
		return Session{}, fmt.Errorf("failed to change session's question")
	}

	return session, nil
}

func (s *sessionServiceImpl) UpdateScore(sessionId string, participantId string, score int) (Session, error) {
	sessionDocRef, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, err
	}

	if session.Status == statusEnded {
		s.logger.Error("attempt to join ended session")
		return Session{}, fmt.Errorf("cannot update this session")
	}

	participant, ok := session.Participants[participantId]
	if !ok {
		s.logger.Error("participant is not in this session")
		return Session{}, fmt.Errorf("cannot update score for this participant")
	}

	participant.Score = score
	session.Participants[participantId] = participant

	_, err = sessionDocRef.Set(s.ctx, session)

	if err != nil {
		s.logger.Error("failed to update participant score", err)
		return Session{}, fmt.Errorf("failed to update participant score")
	}

	return session, nil
}

func (s *sessionServiceImpl) EndSession(sessionId string, userId string) (Session, error) {
	sessionDocRef, sessionDoc, err := s.getSessionDoc(sessionId)
	if err != nil {
		return Session{}, err
	}

	session, err := s.parseSession(sessionDoc)
	if err != nil {
		return Session{}, err
	}

	if userId != session.CreatedBy {
		return Session{}, fmt.Errorf("user not allowed to update session")
	}

	session.Status = statusEnded

	_, err = sessionDocRef.Set(s.ctx, session)

	if err != nil {
		s.logger.Error("failed to end session", err)
		return Session{}, fmt.Errorf("failed to end session")
	}

	return session, nil
}

func (s *sessionServiceImpl) parseSession(sessionDoc *firestore.DocumentSnapshot) (Session, error) {
	var session Session
	if err := sessionDoc.DataTo(&session); err != nil {
		s.logger.Error("failed to parse session data", err)
		return session, fmt.Errorf("failed to parse session data")
	}
	return session, nil
}

func (s *sessionServiceImpl) getSessionDoc(sessionId string) (sessionDocRef *firestore.DocumentRef, sessionDoc *firestore.DocumentSnapshot, err error) {
	if sessionId == "" {
		return sessionDocRef, sessionDoc, fmt.Errorf("missing session id")
	}

	sessionDocRef = s.client.Collection("sessions").Doc(sessionId)
	sessionDoc, err = sessionDocRef.Get(s.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return sessionDocRef, sessionDoc, fmt.Errorf("session not found")
		}
		s.logger.Error("fail to get session data", err)
		return sessionDocRef, sessionDoc, fmt.Errorf("fail to get session data")
	}
	return sessionDocRef, sessionDoc, nil
}

func (s *sessionServiceImpl) checkUser(userId string) error {
	if userId == "" {
		return fmt.Errorf("missing user id")
	}

	userDocRef := s.client.Collection("users").Doc(userId)
	_, err := userDocRef.Get(s.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("user not found")
		}
		s.logger.Error("fail to get user data", err)
		return fmt.Errorf("fail to get user data")
	}
	return nil
}
