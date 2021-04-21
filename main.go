package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(0.0.0.0:13307)/bench?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err)
	}
	prepare(2000)
}

func prepare(n int) {
	db.Exec("TRUNCATE m1")
	db.Exec("TRUNCATE m2")
	builder1 := &strings.Builder{}
	builder2 := &strings.Builder{}
	for m := 0; m <= n; m += 1000 {
		builder1.WriteString("INSERT INTO m1 (id) VALUES")
		builder1.WriteString(fmt.Sprintf("(%d)", (m)))
		builder2.WriteString("INSERT INTO m2 (id,m1_id) VALUES")
		builder2.WriteString(fmt.Sprintf("(%d,%d)", (m), (m)))
		for i := m + 1; i < m+1000; i++ {
			builder1.WriteString(fmt.Sprintf(",(%d)", i))
			builder2.WriteString(fmt.Sprintf(",(%d,%d)", i, i))
		}
		db.Exec(builder1.String())
		db.Exec(builder2.String())
		builder1.Reset()
		builder2.Reset()
	}
}
