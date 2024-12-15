package session_service

import "github.com/minhle003/Elsa-test/internal/services/quiz_service"

type SessionService interface {
	CreateSession(quizId string, userId string) (Session, error)
	JoinSession(name string, sessionId string) (Session, string, error)
	GetSession(sessionId string, userId string) (Session, error)
	StartSession(sessionId string, userId string) (Session, error)
	ChangeQuestion(sessionId string, questionIndex int, userId string) (Session, error)
	UpdateScore(sessionId string, participantId string, score int) (Session, error)
	EndSession(sessionId string, userId string) (Session, error)
}

type Session struct {
	ID                   string                 `json:"ID,omitempty"`
	CreatedBy            string                 `json:"CreatedBy,omitempty"`
	Quiz                 quiz_service.Quiz      `json:"QuizId,omitempty"`
	CurrentQuestionIndex int                    `json:"CurrentQuestionIndex,omitempty"`
	CreatedTime          int64                  `json:"CreatedTime,omitempty"`
	StartTime            *int64                 `json:"StartTime,omitempty"`
	EndTime              *int64                 `json:"EndTime,omitempty"`
	Status               string                 `json:"Status,omitempty"`
	Participants         map[string]Participant `json:"Participants"`
}

type Participant struct {
	Nane  string `json:"Name"`
	Score int    `json:"Score"`
}
