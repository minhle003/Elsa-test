package quiz_service

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type quizServiceImpl struct {
	client *firestore.Client
	logger logging.Logger
	ctx    context.Context
}

func NewQuizService(client *firestore.Client, logger logging.Logger, ctx context.Context) QuizService {
	return &quizServiceImpl{
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (q *quizServiceImpl) CreateQuiz(quiz Quiz) (Quiz, error) {
	err := q.checkUser(quiz.CreatedBy)
	if err != nil {
		return Quiz{}, err
	}

	createdAt := time.Now().Unix()

	newQuiz := quiz
	newQuiz.CreatedAt = createdAt

	quizDocRef, _, err := q.client.Collection("quizzes").Add(q.ctx, quiz)
	if err != nil {
		q.logger.Error("fail to create quiz", err)
		return Quiz{}, fmt.Errorf("fail to create quiz")
	}

	newQuiz.ID = quizDocRef.ID

	return newQuiz, nil
}

func (q *quizServiceImpl) UpdateQuiz(quiz Quiz, userId string) (Quiz, error) {
	err := q.checkUser(userId)
	if err != nil {
		return Quiz{}, err
	}

	docRef := q.client.Collection("quizzes").Doc(quiz.ID)
	doc, err := docRef.Get(q.ctx)
	if err != nil {
		return Quiz{}, fmt.Errorf("quiz not found")
	}

	var oldQuiz Quiz
	if err := doc.DataTo(&oldQuiz); err != nil {
		q.logger.Error("failed to parse quiz data", err)
		return Quiz{}, fmt.Errorf("failed to parse quiz data")
	}

	if oldQuiz.CreatedBy != userId {
		return Quiz{}, fmt.Errorf("user not allowed to update this quiz")
	}

	_, err = docRef.Set(q.ctx, quiz)
	if err != nil {
		q.logger.Error("fail to update quiz", err)
		return Quiz{}, fmt.Errorf("fail to update quiz")
	}

	return quiz, err
}

func (q *quizServiceImpl) ListQuizzesByUser(userID string) ([]Quiz, error) {
	err := q.checkUser(userID)
	if err != nil {
		return []Quiz{}, err
	}

	iter := q.client.Collection("quizzes").Where("CreatedBy", "==", userID).Documents(q.ctx)
	defer iter.Stop()

	var quizzes []Quiz
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			q.logger.Error("failed to query quizzes", err)
			return quizzes, fmt.Errorf("failed to query quizzes")
		}

		var quiz Quiz
		if err := doc.DataTo(&quiz); err != nil {
			q.logger.Error("failed to parse quiz data", err)
			return quizzes, fmt.Errorf("failed to parse quiz data")
		}
		quiz.ID = doc.Ref.ID
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (q *quizServiceImpl) checkUser(userId string) error {
	if userId == "" {
		return fmt.Errorf("missing user id")
	}

	userDocRef := q.client.Collection("users").Doc(userId)
	_, err := userDocRef.Get(q.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("user not found")
		}
		q.logger.Error("fail to get user data", err)
		return fmt.Errorf("fail to get user data")
	}
	return nil
}
