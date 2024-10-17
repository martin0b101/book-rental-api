package user

import (
	"database/sql"

	"github.com/martin0b101/book-rental-api/types"
)

type Store struct {
	database *sql.DB
}

func NewUserStore(database *sql.DB) *Store {
	return &Store{database: database}
}

func (s *Store) GetUsers() ([]types.User, error) {

	rows, err := s.database.Query("SELECT id, first_name, last_name FROM users")

	if err != nil {
		return nil, err
	}

	var users []types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (s *Store) CreateUser(types.User) error {
	return nil
}
