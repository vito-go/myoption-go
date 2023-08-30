package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	err := os.Chdir(filepath.Dir(file))
	if err != nil {
		panic(err)
	}
}

//go:generate go run initdb.go ../doc/postgresql.sql
func main() {
	port := flag.Int64("p", 5432, "specify the port")
	host := flag.String("ip", "localhost", "specify the host")
	database := flag.String("d", "myoption", "specify the database")
	flag.Usage = func() {
		s := "initdb initializes a PostgreSQL database for myoption.\n\nUsage: initdb [OPTION]... <sql file names(ending in .sql)>...\nOPTIONS:"
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", s)
		flag.PrintDefaults()
	}
	flag.Parse()
	fileList := flag.Args()

	dsn := fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable", *host, *port, *database)
	log.Printf("ready to init database. DSN: %s\n", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	for _, s := range fileList {
		if !strings.HasSuffix(s, ".sql") {
			log.Println("not sql file: " + s)
			os.Exit(1)
		}
	}
	err = db.Ping()
	if err != nil {
		// panic: pq: database "myoption" does not exist
		// you should create the database first
		panic(err)
	}
	if len(fileList) == 0 {
		log.Printf("at least one sql file\n")
		os.Exit(1)
		return
	}
	log.Printf("sql file list: %+v\n", fileList)
	if err = execSQL(db, fileList); err != nil {
		panic(err)
	}
	log.Println("init database successfully!")
}

func execSQL(tx *sql.DB, fileList []string) error {
	for _, p := range fileList {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		_, err = tx.Exec(string(b))
		if err != nil {
			return err
		}
	}
	return nil
}
