package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	controller "taskclibackend/controllers"
	"time"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
)

func initDb() (*sql.DB, error) {
	dsn := "root:Root@125502@tcp(127.0.0.1:3306)/GO_PROJECTS"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

// LoggingMiddleware logs the details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", color.BlueString(r.Method), r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()
	// Init DB
	db, err := initDb()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	fmt.Println(color.GreenString("Database connected !"))
	defer db.Close()

	// Routes
	mux.HandleFunc("/ping", controller.Ping(db))
	mux.HandleFunc("/tasks", controller.GetTasks(db))

	// Server Listening
	fmt.Println(color.YellowString("Listening on port: 5000"))

	log.Fatal(http.ListenAndServe(":5000", LoggingMiddleware(mux)))
}
