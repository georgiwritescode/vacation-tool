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
