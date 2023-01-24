package repository

import (
	"database/sql"
	"fmt"
	"log"
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
	l  *log.Logger
}

// NewSqliteUserRepository returns SqliteUserRepository instance.
func NewSqliteUserRepository(l *log.Logger) (*SqliteUserRepository, error) {
	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	return &SqliteUserRepository{
		db: db,
		l:  l,
	}, nil
}

func (ur *SqliteUserRepository) Close() error {
	return ur.db.Close()
}

// InvalidCredentialsErr error message used when credentials are wrong
const InvalidCredentialsErr = "Invalid credentials"

// ValidateCredentials validates user credentials using sqlite database.
func (ur *SqliteUserRepository) ValidateCredentials(user string, hash string) error {
	ur.l.Printf("Validating user credentials %s, %s", user, hash)

	rows, err := ur.db.Query("SELECT id FROM user WHERE name = ? AND password_hash = ?", user, hash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		ur.l.Printf("Invalid credentials %s, %s", user, hash)
		return fmt.Errorf(InvalidCredentialsErr)
	}

	ur.l.Printf("Valid credentials %s, %s", user, hash)

	return nil
}
