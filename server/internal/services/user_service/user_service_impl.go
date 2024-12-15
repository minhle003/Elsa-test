package user_service

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"google.golang.org/api/iterator"
)

type userServiceImpl struct {
	client *firestore.Client
	logger logging.Logger
	ctx    context.Context
}

func NewUserService(client *firestore.Client, logger logging.Logger, ctx context.Context) UserService {
	return &userServiceImpl{
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}

func (u *userServiceImpl) CreateUser(userName string) (User, error) {
	iter := u.client.Collection("users").Where("Name", "==", userName).Documents(u.ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			u.logger.Error("failed to check existing users", err)
			return User{}, fmt.Errorf("failed to check existing users")
		}
		if doc.Exists() {
			return User{}, fmt.Errorf("user %s already exists", userName)
		}
	}

	newUser := User{
		Name: userName,
	}

	docRef, _, err := u.client.Collection("users").Add(u.ctx, newUser)
	if err != nil {
		u.logger.Error("failed to create user", err)
		return User{}, fmt.Errorf("failed to create user")
	}

	newUser.ID = docRef.ID

	return newUser, nil
}

func (u *userServiceImpl) GetUserByUserName(userName string) (User, error) {
	iter := u.client.Collection("users").Where("Name", "==", userName).Documents(u.ctx)
	defer iter.Stop()
	var user User
	for {
		doc, err := iter.Next()
		if err != nil {
			if err.Error() == "iterator.Done" {
				break
			}
			u.logger.Error("failed to get user info", err)
			return user, fmt.Errorf("failed to get user info")
		}

		user = User{
			ID:   doc.Ref.ID,
			Name: doc.Data()["Name"].(string),
		}
		break
	}

	if user.ID == "" {
		return user, fmt.Errorf("user not found")
	}

	return user, nil
}
