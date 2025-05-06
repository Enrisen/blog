package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Enrisen/blog/internal/validator"
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

// Authenticate verifies a user's credentials
func (m *UserModel) Authenticate(email, password string) (*User, error) {
	// Query for the user with the provided email
	query := `
		SELECT user_id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check if the provided password matches the stored hash
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	return &user, nil
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

// ValidateUserRegistration validates the user registration form data
func ValidateUserRegistration(v *validator.Validator, name, email, password, confirmPassword string) {
	// Validate name
	v.Check(validator.NotBlank(name), "name", "Name cannot be empty")
	v.Check(validator.MaxLength(name, 100), "name", "Name cannot be more than 100 characters")

	// Validate email
	v.Check(validator.NotBlank(email), "email", "Email cannot be empty")
	v.Check(validator.IsValidEmail(email), "email", "Please enter a valid email address")
	v.Check(validator.MaxLength(email, 255), "email", "Email cannot be more than 255 characters")

	// Validate password
	v.Check(validator.NotBlank(password), "password", "Password cannot be empty")
	v.Check(validator.MinLength(password, 8), "password", "Password must be at least 8 characters")

	// Validate password confirmation
	v.Check(validator.NotBlank(confirmPassword), "confirm_password", "Please confirm your password")
	v.Check(password == confirmPassword, "confirm_password", "Passwords do not match")
}

// ValidateLogin validates the login form data
func ValidateLogin(v *validator.Validator, email, password string) {
	// Validate email
	v.Check(validator.NotBlank(email), "email", "Email cannot be empty")
	v.Check(validator.IsValidEmail(email), "email", "Please enter a valid email address")

	// Validate password
	v.Check(validator.NotBlank(password), "password", "Password cannot be empty")
}
