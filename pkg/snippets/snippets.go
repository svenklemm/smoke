package snippets

import (
	"database/sql"
)

func foobar(a string) {
	db, _ := sql.Open("postgres", "database=postgres")
	db.Query(a)
	db.Exec(a)
}
