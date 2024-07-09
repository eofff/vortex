package migrations

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sort"
	"strconv"

	_ "github.com/lib/pq"
)

type Migration struct {
	Name string
	Num  int
}

var (
	migrationFiles   []Migration
	migrationEntries []Migration
)

func Migrate(db *sql.DB) {
	rows, err := db.Query(`SELECT EXISTS (
						SELECT FROM information_schema.tables 
						WHERE  table_schema = 'public'
						AND    table_name   = '_migrations'
						);`)
	if err != nil {
		log.Fatal(err)
	}

	var migrationTableExists bool
	rows.Next()
	err = rows.Scan(&migrationTableExists)
	if err != nil {
		log.Fatal(err)
	}

	if !migrationTableExists {
		_, err = db.Query("CREATE TABLE _migrations (name varchar PRIMARY KEY);")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Table _migrations successfully created")
	}

	readMigrationsFolder()
	readMigrationTable(db)

	if len(migrationFiles) > len(migrationEntries) {
		for i := len(migrationEntries); i < len(migrationFiles); i++ {
			sqlBytes, err := os.ReadFile("./migrations/" + migrationFiles[i].Name + "/apply.sql")
			sqlString := string(sqlBytes)
			if err != nil {
				log.Fatal(err)
			}

			ctx := context.Background()
			tx, err := db.BeginTx(ctx, nil)

			if err != nil {
				log.Fatal(err)
			}

			_, err = tx.ExecContext(ctx, sqlString)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}

			_, err = tx.ExecContext(ctx, `INSERT INTO _migrations(name) VALUES ($1)`, migrationFiles[i].Name)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}

			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Migration %s succesfully applied\n", migrationFiles[i].Name)
			}
		}
	}
}

func readMigrationsFolder() {
	migrationFiles = make([]Migration, 0)

	entries, err := os.ReadDir("./migrations")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		migration, err := parseMigration(e.Name())
		if err != nil {
			log.Fatal(err)
		}

		migrationFiles = append(migrationFiles, migration)
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Num < migrationFiles[j].Num
	})
}

func readMigrationTable(db *sql.DB) {
	rows, err := db.Query(`SELECT "name" FROM "_migrations"`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		migration, err := parseMigration(name)
		if err != nil {
			log.Fatal(err)
		}

		migrationEntries = append(migrationEntries, migration)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(migrationEntries, func(i, j int) bool {
		return migrationEntries[i].Num < migrationEntries[j].Num
	})
}

func parseMigration(s string) (Migration, error) {
	num, err := strconv.Atoi(s[0:4])

	return Migration{
		Num:  num,
		Name: s,
	}, err
}
