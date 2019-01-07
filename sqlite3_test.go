package learning

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateDB(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:/tmp/test.db")
	if err != nil {
		t.Error("Could not open DB: ", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Error("Could not ping DB: ", err)
		return
	}

	createSQL := `CREATE TABLE IF NOT EXISTS "USERS" (
					"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
					"email" varchar(75) NOT NULL UNIQUE,
					"password" varchar(128) NOT NULL,
					"first_name" varchar(30) NOT NULL,
					"last_name" varchar(30) NOT NULL,
					"is_active" bool NOT NULL)`
	_, err = db.Exec(createSQL)
	if err != nil {
		t.Error("Could not create USERS table: ", err)
		return
	}

	insertRecordTest(t, db, "arunsworld@gmail.com", 1)
	insertRecordTest(t, db, "arun@e2open.com", 2)

	queryOneRecordTest(t, db)
	genericQueryTest(t, db)
	noRecordFoundTest(t, db)
	deleteRecordTest(t, db)

	err = os.Remove("/tmp/test.db")
	if err != nil {
		t.Error("Could not delete the database: ", err)
		return
	}
}

func insertRecordTest(t *testing.T, db *sql.DB, email string, expectedID int64) {
	insertSQL := `INSERT INTO "USERS" ("email", "password", "first_name", "last_name", "is_active")
	VALUES ($1, $2, $3, $4, $5)`
	result, err := db.Exec(insertSQL, email, "password", "Arun", "Barua", true)
	if err != nil {
		t.Error("Could not insert into USERS table: ", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		t.Error("Error with rows affected: ", err)
		return
	}
	if rows != 1 {
		t.Error("Expected 1 row to be effected. Got: ", rows)
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		t.Error("Error with inserted ID: ", err)
		return
	}
	if insertedID != expectedID {
		t.Errorf("Expected inserted ID to be %d. Got: %d", expectedID, insertedID)
	}
}

func queryOneRecordTest(t *testing.T, db *sql.DB) {
	var email string
	var firstName string
	err := db.QueryRow("SELECT email, first_name FROM USERS WHERE id = $1", "1").Scan(&email, &firstName)
	if err != nil {
		t.Error("Error while querying one record: ", err)
		return
	}
	if email != "arunsworld@gmail.com" {
		t.Error("Email not matching. Found: " + email)
	}
	if firstName != "Arun" {
		t.Error("First Name not matching. Found: " + firstName)
	}
}

func noRecordFoundTest(t *testing.T, db *sql.DB) {
	var email string
	var firstName string
	err := db.QueryRow("SELECT email, first_name FROM USERS WHERE id = $1", "10").Scan(&email, &firstName)
	if err != sql.ErrNoRows {
		t.Error("Expected no rows error to be returned but got error: ", err)
		return
	}
	if err == nil {
		t.Error("Expected no rows error but got nothing.")
	}
}

func deleteRecordTest(t *testing.T, db *sql.DB) {
	result, err := db.Exec("DELETE FROM USERS WHERE id=$1", "1")
	if err != nil {
		t.Error("Error while deleting user:", err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		t.Error("Error while deleting user:", err)
		return
	}
	if rows != 1 {
		t.Error("Expected one record to be deleted. Got: ", rows)
		return
	}
}

func genericQueryTest(t *testing.T, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM USERS WHERE first_name = $1 ORDER BY id", "Arun")
	if err != nil {
		t.Log("Error creating Select query: ", err)
		return
	}
	cols, err := rows.Columns()
	if err != nil {
		t.Log("Error getting columns: ", err)
		return
	}
	allCols := strings.Join(cols, "|")
	if allCols != "id|email|password|first_name|last_name|is_active" {
		t.Log("Columns not as expected. Got: ", cols)
		return
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	type row map[string]string
	type table []row
	data := table{}
	for rows.Next() {
		err = rows.Scan(vals...)
		newRow := make(row)
		for i, v := range vals {
			vv := v.(*sql.RawBytes)
			if len(*vv) == 0 {
				continue
			}
			newRow[cols[i]] = string(*vv)
		}
		data = append(data, newRow)
	}
	rows.Close()

	if len(data) != 2 {
		t.Error("Expected to see 2 rows. Found:", len(data))
	}
	for i, r := range data {
		if r["id"] != strconv.Itoa(i+1) {
			t.Errorf("Expected id to be %s. Found %s.", strconv.Itoa(i), r["id"])
		}
	}
}
