package user

import (
	"database/sql"
	"fmt"

	"github.com/georgiwritescode/vacation-tool/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) FindById(id int) (*types.User, error) {

	rows, err := s.db.Query("SELECT * from tbl_users where id = ?", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(req *types.User) (int, error) {

	res, err := s.db.Exec("INSERT INTO tbl_users (first_name, last_name, age, email) values (?, ?, ?, ?)", req.FirstName, req.LastName, req.Age, req.Email)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (s *Store) FetchAllUsers() ([]*types.User, error) {
	rows, err := s.db.Query("SELECT * from tbl_users")
	if err != nil {
		return nil, err
	}

	var users []*types.User
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Email,
		&user.Timestamp,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
