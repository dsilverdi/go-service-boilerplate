package main

import (
	"context"
	"go-service-boilerplate/app/core"
	"go-service-boilerplate/app/server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("env/APP_ENV.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := initMysql()

	app := core.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiMux := server.NewService(r, app)

	rootMux := chi.NewRouter()
	rootMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		log.Print("TODO: /healthz actual check e.g. db and other services")
		_, _ = w.Write([]byte("OK!"))
	})

	rootMux.Mount("/api", apiMux.MountServerRoute())

	apiServer := http.Server{
		Addr:    ":8080",
		Handler: rootMux,
	}

	var shuttingDown bool
	shutdownSignal := make(chan os.Signal)
	//signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		log.Print("service is ready...")
		err := apiServer.ListenAndServe()
		if err != nil && (err != http.ErrServerClosed || !shuttingDown) {
			log.Fatalln("server stopped unexpectedly")
		}
	}()

	<-shutdownSignal
	shuttingDown = true

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_ = apiServer.Shutdown(shutdownCtx)
	log.Print("done...")

}

func initMysql() *sqlx.DB {
	host := os.Getenv("APP_DATABASE_HOST")
	port := os.Getenv("APP_DATABASE_PORT")
	dbname := os.Getenv("APP_DATABASE_NAME")
	user := os.Getenv("APP_DATABASE_USER")
	pwd := os.Getenv("APP_DATABASE_PASSWORD")

	dsn := user + ":" + pwd + "@(" + host + ":" + port + ")/" + dbname + "?parseTime=true&clientFoundRows=true"

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	log.Print("Using MySQL at ", host, ":", port, "/", dbname)

	ddl, err := ioutil.ReadFile("db/ddl.sql")
	if err != nil {
		panic(err)
	}
	sqls := string(ddl)

	// We need to split the DDL query by `;` and execute it one by one.
	// Because sql.DB.Exec() from mysql driver cannot executes multiple query at once
	// and it will give weird syntax error messages.
	splitted := strings.Split(sqls, ";")

	for _, v := range splitted {
		trimmed := strings.TrimSpace(v)

		if len(trimmed) > 0 {
			_, err = db.Exec(v)

			if err != nil {
				me, ok := err.(*mysql.MySQLError)
				if !ok {
					panic("Error executing DDL query")
				}

				// http://dev.mysql.com/doc/refman/5.7/en/error-messages-server.html
				// We will skip error duplicate key name in database (code: 1061),
				// because CREATE INDEX doesn't have IF NOT EXISTS clause,
				// otherwise we will stop the loop and print the error
				if me.Number == 1061 {

				} else {
					log.Print(err)
					return db
				}
			}
		}
	}

	log.Print("DDL file executed")

	return db
}
