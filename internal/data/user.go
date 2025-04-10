package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID             int64     `json:"id" db:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-" db:"password_hash"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

// RegisterUser handles user registration including validation
func (m *UserModel) RegisterUser(name, email, password, confirmPassword string) (*User, error) {
	// Create a new user
	user := &User{
		Name:  name,
		Email: email,
	}

	// Hash the password and insert the user
	return user, m.Insert(name, email, password)
}

func (m *UserModel) Insert(name, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id, created_at`

	args := []interface{}{name, email, string(hashedPassword)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int64
	var createdAt time.Time
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&newID, &createdAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
