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

func (s *Store) FindAll() ([]*types.Vacation, error) {
	rows, err := s.db.Query("SELECT * FROM tbl_vacations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacations := make([]*types.Vacation, 0)
	for rows.Next() {
		v, err := scanRowsIntoVacation(rows)
		if err != nil {
			return nil, err
		}
		vacations = append(vacations, v)
	}

	return vacations, nil
}

func (s *Store) CreateVacation(vacation *types.Vacation) (int, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	// Check if user has enough days (Paid + NonPaid)
	var currentDays int
	var currentNonPaid int
	err = tx.QueryRow("SELECT vacation_days, non_paid_leave FROM tbl_users WHERE id = ?", vacation.PersonId).Scan(&currentDays, &currentNonPaid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	needed := vacation.DaysUsed
	// Logic: Use paid days first. If paid days < 0 (unlikely) or < needed, use non-paid.
	// Actually, if paid days are negative, we should probably just treat them as 0 available for this calculation?
	// But let's assume standard logic: available paid = max(0, currentDays).
	
	availablePaid := currentDays
	if availablePaid < 0 {
		availablePaid = 0
	}

	deductPaid := 0
	deductNonPaid := 0

	if availablePaid >= needed {
		deductPaid = needed
	} else {
		deductPaid = availablePaid
		deductNonPaid = needed - availablePaid
	}

	if currentNonPaid < deductNonPaid {
		tx.Rollback()
		return 0, fmt.Errorf("insufficient leave: need %d, have %d paid + %d non-paid", vacation.DaysUsed, currentDays, currentNonPaid)
	}

	// Update user's vacation days and non-paid leave
	_, err = tx.Exec("UPDATE tbl_users SET vacation_days = vacation_days - ?, non_paid_leave = non_paid_leave - ? WHERE id = ?", deductPaid, deductNonPaid, vacation.PersonId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert vacation
	res, err := tx.Exec("INSERT INTO tbl_vacations (label, from_date, to_date, person_id, ts, days_used) VALUES (?, ?, ?, ?, ?, ?)",
		vacation.Label, vacation.FromDate, vacation.ToDate, vacation.PersonId, vacation.Timestamp, vacation.DaysUsed)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) UpdateVacation(vacation *types.Vacation) error {
	_, err := s.db.Exec("UPDATE tbl_vacations SET label=?, from_date=?, to_date=?, person_id=?, ts=?, days_used=? WHERE id=?",
		vacation.Label, vacation.FromDate, vacation.ToDate, vacation.PersonId, vacation.Timestamp, vacation.DaysUsed, vacation.ID)
	return err
}

func (s *Store) DeleteVacation(id int) error {
	_, err := s.db.Exec("DELETE FROM tbl_vacations WHERE id=?", id)
	return err
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
		&vacation.DaysUsed,
	)
	if err != nil {
		return nil, err
	}

	return vacation, nil
}
