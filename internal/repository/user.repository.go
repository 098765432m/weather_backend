package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	dtoReq "github.com/098765432m/internal/dto/request"
	dtoRes "github.com/098765432m/internal/dto/response"
	"github.com/098765432m/internal/model"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	db *sql.DB
	redisClient *redis.Client
}

func NewUserRepository(db *sql.DB, redisClient *redis.Client) *UserRepository {
	return &UserRepository{
		db: db,
		redisClient: redisClient,
	}
}

func (repo *UserRepository) GetUserById(id int) (*model.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)

	cachedUser, err := repo.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.User
		if err = json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}
	query := "SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE id = ?"

	user := &model.User{}

	row := repo.db.QueryRow(query, id)
	
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user *dtoReq.CreatedUserDtoRequest) (error) {
	query := "INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)"
	
	_, err := repo.db.Exec(query, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) DashboardUpdateUser(id int, user *dtoReq.DashBoardUpdateUserDtoRequest) (error) {
	query := "UPDATE users SET username = ?, email = ?, password = ?, role = ? WHERE id = ?"
	
	_, err := repo.db.Exec(query, user.Username, user.Email, user.Password, user.Role, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) DeleteUser(id int) (error) {
	query := "DELETE FROM users WHERE id = ?"
	
	_, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
};

func (repo *UserRepository) GetUserByUsernameOrEmail(usernameOrEmail string) (*dtoRes.UserResponse, error) {
	query := "SELECT id, username, email, password, role FROM users WHERE username = ? OR email = ?"

	row := repo.db.QueryRow(query, usernameOrEmail, usernameOrEmail)
	
	user := &dtoRes.UserResponse{}

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) ChangePassword(id int, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	_, err := repo.db.Exec(query, newPassword, id)
	if err != nil {
		return err
	}

	return nil
}