package dtoReq

import roles "github.com/098765432m/internal/utils"

type CreatedUserDtoRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     roles.Role `json:"role,omitempty"`
}

type DashBoardUpdateUserDtoRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     roles.Role `json:"role,omitempty"`
}

type LoginDtoRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}