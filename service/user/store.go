package user

import (
	"database/sql"

	"github.com/martin0b101/book-rental-api/types"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
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

func (s *Store) CreateUser(user types.RegisterUserRequest) (*types.User, error) {

	// _, err := s.database.Exec("INSERT INTO users (first_name, last_name) VALUES ($1, $2)", 
	// user.FirstName, user.LastName)

	var userRegistered types.User
	// The INSERT query with the RETURNING clause
	query := `INSERT INTO users (first_name, last_name) 
	          VALUES ($1, $2) RETURNING id, first_name, last_name`

	// Use QueryRow to execute the query and return the inserted row
	err := s.database.QueryRow(query, user.FirstName, user.LastName).
		Scan(&userRegistered.Id, &userRegistered.FirstName, &userRegistered.LastName)

	if err != nil{
		return nil, err
	}

	return &userRegistered, nil
}
