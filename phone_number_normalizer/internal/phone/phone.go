package phone

import (
	"context"
	"database/sql"
	"regexp"
)

type Phone struct {
	ID     int
	Number string
}

var nonDigits = regexp.MustCompile(`\D+`)

// ResetAndCreateDB drops and creates the DB using an admin connection to `postgres`.
func ResetAndCreateDB(adminDB *sql.DB, name string) error {
	if _, err := adminDB.Exec("DROP DATABASE IF EXISTS " + name); err != nil {
		return err
	}
	if _, err := adminDB.Exec("CREATE DATABASE " + name); err != nil {
		return err
	}
	return nil
}

func CreatePhoneTable(ctx context.Context, db *sql.DB) error {
	stmt := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
        id SERIAL PRIMARY KEY,
        value VARCHAR(255) NOT NULL
    )`
	_, err := db.ExecContext(ctx, stmt)
	return err
}

func SeedIfEmpty(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(1) FROM phone_numbers").Scan(&count); err != nil {
		// table might be empty or not exist; let caller handle creation first
		return err
	}
	if count > 0 {
		return nil
	}
	init := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, v := range init {
		if _, err := InsertPhone(ctx, db, v); err != nil {
			return err
		}
	}
	return nil
}

func InsertPhone(ctx context.Context, db *sql.DB, value string) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`, value).Scan(&id)
	return id, err
}

func GetAllPhones(ctx context.Context, db *sql.DB) ([]Phone, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func FindPhone(ctx context.Context, db *sql.DB, value string) (*Phone, error) {
	var p Phone
	row := db.QueryRowContext(ctx, "SELECT id, value FROM phone_numbers WHERE value=$1", value)
	if err := row.Scan(&p.ID, &p.Number); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func UpdatePhone(ctx context.Context, db *sql.DB, id int, value string) error {
	_, err := db.ExecContext(ctx, "UPDATE phone_numbers SET value=$2 WHERE id=$1", id, value)
	return err
}

func DeletePhone(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM phone_numbers WHERE id=$1", id)
	return err
}

// ProcessAll normalizes numbers and updates/removes duplicates.
func ProcessAll(ctx context.Context, db *sql.DB) error {
	phones, err := GetAllPhones(ctx, db)
	if err != nil {
		return err
	}
	for _, p := range phones {
		norm := Normalize(p.Number)
		if norm == p.Number {
			continue
		}
		existing, err := FindPhone(ctx, db, norm)
		if err != nil {
			return err
		}
		if existing != nil {
			if err := DeletePhone(ctx, db, p.ID); err != nil {
				return err
			}
		} else {
			if err := UpdatePhone(ctx, db, p.ID, norm); err != nil {
				return err
			}
		}
	}
	return nil
}

func Normalize(s string) string {
	return nonDigits.ReplaceAllString(s, "")
}
