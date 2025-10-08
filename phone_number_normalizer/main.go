package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jkbresee"
	password = "password"
	dbname   = "gophercises_phone"
)

type phone struct {
	id     int
	number string
}

func main() {
	var ids []int

	init_phones := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	err = resetDB(db, dbname)
	must(err)

	err = createDB(db, dbname)
	must(err)

	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)

	defer db.Close()
	err = db.Ping()
	must(err)

	err = createPhoneTable(db)
	must(err)

	for _, p := range init_phones {
		id, err := insertPhone(db, p)
		must(err)
		ids = append(ids, id)
	}

	var phones []phone
	phones, err = getAllPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on %+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating/Removing", number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				//delete
			} else {
				//update number
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT value FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.id, &p.number)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
	return &p, nil
}

func getAllPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createPhoneTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
			)`
	_, err := db.Exec(statement)
	return err
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name) //not production-grade sql query!!
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(phone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(phone, "")
}
