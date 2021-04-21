package main_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(0.0.0.0:13307)/bench?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err)
	}
}

func Benchmark_INNERJOIN_20000(b *testing.B) {
	prepare(20000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		innerjoin()
	}
}

func Benchmark_SUBQUERY_20000(b *testing.B) {
	prepare(20000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subquery()
	}
}

func Benchmark_INNERJOIN_100000(b *testing.B) {
	prepare(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		innerjoin()
	}
}

func Benchmark_SUBQUERY_100000(b *testing.B) {
	prepare(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subquery()
	}
}

func Benchmark_INNERJOIN_1000000(b *testing.B) {
	prepare(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		innerjoin()
	}
}

func Benchmark_SUBQUERY_1000000(b *testing.B) {
	prepare(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subquery()
	}
}

func Benchmark_INNERJOIN_100000_P5(b *testing.B) {
	prepare(100000)

	b.SetParallelism(5)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			innerjoin()
		}
	})
}

func Benchmark_SUBQUERY_100000_P5(b *testing.B) {
	prepare(100000)

	b.SetParallelism(5)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			subquery()
		}
	})
}

func innerjoin() {
	_, err := db.Exec("SELECT m1.id FROM m1 INNER JOIN m2 ON m1.id = m2.m1_id")
	if err != nil {
		panic(err)
	}
}

func subquery() {
	_, err := db.Exec("SELECT m1.id FROM m1 WHERE EXISTS (SELECT 1 FROM m2 WHERE m1.id = m2.m1_id)")
	if err != nil {
		panic(err)
	}
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
