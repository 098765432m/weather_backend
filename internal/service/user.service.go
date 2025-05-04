package service

import (
	"fmt"

	"github.com/098765432m/internal/auth"
	dtoReq "github.com/098765432m/internal/dto/request"
	"github.com/098765432m/internal/repository"
	roles "github.com/098765432m/internal/utils"
	"github.com/098765432m/logger"
	"golang.org/x/crypto/bcrypt"
)


type UserService struct {
	repo *repository.UserRepository
}
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GuestRegister(user *dtoReq.CreatedUserDtoRequest) error {
	if (user.Username == "" || user.Email == "" || user.Password == "") {
		return fmt.Errorf("username, email, and password are required")
	}

	// Default role
	if user.Role == "" {
		user.Role = roles.RoleGuest
	}

	if  !roles.IsValidRole(user.Role) {
		logger.NewLogger().Error.Printf("Invalid role: %s", user.Role)
		return fmt.Errorf("invalid role: %s", user.Role)
	}
	
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.NewLogger().Error.Println("Failed to hash password: ", err)
		return err
	}

	user.Password = string(pass)

	if err := s.repo.CreateUser(user); err != nil {
		logger.NewLogger().Error.Println("Failed to create user: ", err)
		return err
	}

	return nil
}

func (s *UserService) Login(user *dtoReq.LoginDtoRequest) (string, error) {
	if user.UsernameOrEmail == "" || user.Password == "" {
		return "", fmt.Errorf("username or email and password are required")
	}

	userFromDB, err := s.repo.GetUserByUsernameOrEmail(user.UsernameOrEmail)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := auth.GenerateJWT(userFromDB.Username, userFromDB.Email, userFromDB.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	
	return token, nil
}

func (s *UserService) DashboardUpdateUser(id int, user *dtoReq.DashBoardUpdateUserDtoRequest) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("username, email, and password are required")
	}

	if user.Role != "" && !roles.IsValidRole(user.Role) {
		logger.NewLogger().Error.Printf("Invalid role: %s", user.Role)
		return fmt.Errorf("invalid role: %s", user.Role)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.NewLogger().Error.Println("Failed to hash password: ", err)
		return err
	}

	user.Password = string(pass)

	if err := s.repo.DashboardUpdateUser(id, user); err != nil {
		logger.NewLogger().Error.Println("Failed to update user: ", err)
		return err
	}

	return nil
}

func (s *UserService) DeleteUser (id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user Id : %d", id)
	}

	if err := s.repo.DeleteUser(id); err != nil {
		logger.NewLogger().Error.Printf("Failedto delete user %d : %w", id, err)
		return err
	}

	return nil
}

func (s *UserService) ChangeUserPassword(id int, newPassword string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user Id : %d", id)
	}

	if newPassword == "" {
		return fmt.Errorf("new password is required")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.NewLogger().Error.Println("Failed to hash new password: ", err)
		return err
	}

	err = s.repo.ChangePassword(id, string(hashedNewPassword))
	if err != nil {
		logger.NewLogger().Error.Printf("Failed to change password for user %d", id)
		return err
	}

	return nil
}