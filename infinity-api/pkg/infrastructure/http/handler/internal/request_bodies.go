package internal

import "github.com/satori/go.uuid"

type CreateUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDetailAvailable struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type PasswordReset struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
}

type NewPassword struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Uuid struct {
	Uuid uuid.UUID `json:"uuid"`
}
