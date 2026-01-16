package vacation

import (
	"github.com/georgiwritescode/vacation-tool/types"
)

func (s *Store) GetActiveVacations(date string) ([]*types.Vacation, error) {
	rows, err := s.db.Query("SELECT * FROM tbl_vacations WHERE ? BETWEEN from_date AND to_date", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacations := make([]*types.Vacation, 0)
	for rows.Next() {
		v := new(types.Vacation)
		// We can reuse scanRowsIntoVacation if we export it or duplicate logic
		// But scanRowsIntoVacation is in this file (unexported).
		v, err = scanRowsIntoVacation(rows)
		if err != nil {
			return nil, err
		}
		vacations = append(vacations, v)
	}
	return vacations, nil
}
