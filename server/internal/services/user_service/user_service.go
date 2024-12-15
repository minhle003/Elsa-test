package user_service

type UserService interface {
	CreateUser(userName string) (User, error)
	GetUserByUserName(userName string) (User, error)
}

type User struct {
	ID   string `json:"ID,omitempty"`
	Name string `json:"Name"`
}
