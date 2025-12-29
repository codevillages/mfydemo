package model

// CreateUserInput 输入参数。
type CreateUserInput struct {
	Username string
	Password string
	Email    string
	Phone    string
	Nickname string
	Avatar   string
	Status   int
}
