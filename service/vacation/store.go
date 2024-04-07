package vacation

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

func (s *Store) FindById(id int) (*types.Vacation, error) {

	rows, err := s.db.Query("SELECT * from tbl_vacations where id = ?", id)
	if err != nil {
		return nil, err
	}

	vacation := new(types.Vacation)
	for rows.Next() {
		vacation, err = scanRowsIntoVacation(rows)
		if err != nil {
			return nil, err
		}
	}

	if vacation.ID == 0 {
		return nil, fmt.Errorf("vacation not found :( ")
	}

	return vacation, nil
}

func scanRowsIntoVacation(rows *sql.Rows) (*types.Vacation, error) {
	vacation := new(types.Vacation)

	err := rows.Scan(
		&vacation.ID,
		&vacation.Label,
		&vacation.FromDate,
		&vacation.ToDate,
		&vacation.PersonId,
		&vacation.Timestamp,
	)
	if err != nil {
		return nil, err
	}

	return vacation, nil
}
