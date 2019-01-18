package learning

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/apache/calcite-avatica-go/v3"
)

func TestOpenAndPingAWSDB(t *testing.T) {

	if os.Getenv("PHOENIX_AVAILABLE") == "" {
		t.Skip("Set the PHOENIX_AVAILABLE flag to run this test when HDP Lab is accessible.")
	}

	db, err := sql.Open("avatica", "http://172.16.3.196:8765")
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

	// queryAndPrintSingleValueResults(t, db, `select DISTINCT("TABLE_NAME") from SYSTEM.CATALOG`)
	// queryAndPrintSingleValueResults(t, db, `select COUNT(1) from POS_TRANSACTION_ZYMECUSTOMER`)
	// cols := getColumnsFromTable(t, db, "SYSTEM.CATALOG")
	// fmt.Println(strings.Join(cols, ","))
	// query := `select * from SYSTEM.CATALOG`
	// query := `select TABLE_SCHEM, TABLE_NAME, SALT_BUCKETS, DISABLE_WAL, COLUMN_COUNT, GUIDE_POSTS_WIDTH from SYSTEM.CATALOG where COLUMN_NAME is null`
	// result := genericQuery(t, db, query)
	// for _, row := range result {
	// fmt.Println(strings.Join(row, ","))
	// }

}

func queryAndPrintSingleValueResults(t *testing.T, db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		t.Fatal("Error creating query: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			t.Fatal(err)
		}
		fmt.Println(val)
	}
}

func getColumnsFromTable(t *testing.T, db *sql.DB, table string) []string {
	rows, err := db.Query(fmt.Sprintf(`select * from %s LIMIT 1`, table))
	if err != nil {
		t.Fatal("Error creating query: ", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatal("Error getting columns: ", err)
	}
	return cols
}

func genericQuery(t *testing.T, db *sql.DB, query string) [][]string {
	rows, err := db.Query(query)
	if err != nil {
		t.Fatal("Error creating query: ", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatal("Error getting columns: ", err)
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	result := [][]string{}
	result = append(result, cols)
	for rows.Next() {
		err = rows.Scan(vals...)
		newRow := make([]string, len(cols))
		for i, v := range vals {
			vv := v.(*sql.RawBytes)
			newRow[i] = string(*vv)
		}
		result = append(result, newRow)
	}
	return result
}
