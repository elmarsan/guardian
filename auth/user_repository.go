package auth

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// UserRepository interface defines methods for checking user credentials.
type UserRepository interface {
	ValidateCredentials(user string, hash string) error
}

// SqliteUserRepository implements UserRepository using sqlite.
type SqliteUserRepository struct {
	db *sql.DB
}

// NewSqliteUserRepository returns SqliteUserRepository instance.
func NewSqliteUserRepository() (*SqliteUserRepository, error) {
	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	return &SqliteUserRepository{
		db: db,
	}, nil
}

func (ur *SqliteUserRepository) Close() error {
	return ur.db.Close()
}

// InvalidCredentialsErr error message used when credentials are wrong
const InvalidCredentialsErr = "Invalid credentials"

// ValidateCredentials validates user credentials using sqlite database.
func (ur *SqliteUserRepository) ValidateCredentials(user string, hash string) error {
	rows, err := ur.db.Query("SELECT id FROM user WHERE name = ? AND password_hash = ?", user, hash)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf(InvalidCredentialsErr)
	}

	return nil
}
