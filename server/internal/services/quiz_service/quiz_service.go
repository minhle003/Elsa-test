package quiz_service

type QuizService interface {
	CreateQuiz(quiz Quiz) (Quiz, error)
	UpdateQuiz(quiz Quiz, userId string) (Quiz, error)
	ListQuizzesByUser(userId string) ([]Quiz, error)
}

type Quiz struct {
	ID        string     `json:"ID,omitempty"`
	CreatedBy string     `json:"CreatedBy,omitempty"`
	Title     string     `json:"Title,omitempty"`
	CreatedAt int64      `json:"CreatedAt,omitempty"`
	Questions []Question `json:"Questions,omitempty"`
}

type Question struct {
	ID          string   `json:"ID"`
	Title       string   `json:"Title"`
	Type        string   `json:"Type"`
	Description string   `json:"Description"`
	Score       int      `json:"Score"`
	Time        int      `json:"Time"`
	Choices     []Choice `json:"Choices,omitempty"`
	Answer      string   `json:"Answer"`
}

type Choice struct {
	Description string `json:"Description"`
}
