package user

import "context"

type User struct {
	ID       int64  `json:"id" db:"id`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email`
	Password string `json:"password" db:"password`
}

// Request struct represents what we”ll be gettting from the user.
// Arequest will be for the Username, Email and Password
type CreateUserReq struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email`
	Password string `json:"password" db:"password`
}

// A response of the request will be id, username and email
type CreateUserRes struct {
	ID       string `json:"id" db:"id`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type Service interface {
	// Define the method signature
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
}
