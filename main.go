package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"1mk.no/saunagauge/internal/html"
	"1mk.no/saunagauge/internal/sauna"
	"1mk.no/saunagauge/internal/server"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	Port       string
	Production bool
	DBFile     string
}

var fsys fs.FS

//go:embed database/migrations
var migrations embed.FS

func main() {
	cfg := config{
		Port:       "8080",
		DBFile:     "database/database.sqlite3",
		Production: true,
	}

	var migrateOnly bool

	flag.StringVar(&cfg.Port, "port", "8080", "TCP port to listen on")
	flag.StringVar(&cfg.DBFile, "database", cfg.DBFile, "Path to SQLite database file")
	flag.BoolVar(&cfg.Production, "prod", true, "Set to false when developing for a better experience")
	flag.BoolVar(&migrateOnly, "migrate", false, "Run migrations")
	flag.Parse()

	db, err := sqlx.Connect("sqlite3", cfg.DBFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbLogger := log.New(os.Stderr, "db :: ", log.LstdFlags)

	if migrateOnly {
		if err = ensureSchema(db.DB); err != nil {
			db.Close()
			panic(err)
		}
		dbLogger.Printf("applied migrations")
		return
	}

	if cfg.Production {
		fsys = html.EmbedFS
	} else {
		fsys = os.DirFS("internal/html")
	}

	sauna := sauna.New(db)

	httpLogger := log.New(os.Stderr, "http :: ", log.LstdFlags)
	s := server.New(sauna, httpLogger)

	fmt.Println(s)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.FS(fsys)))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/temperature", s.GetTemperature())
	mux.HandleFunc("/humidity", s.GetHumidity())

	server := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: mux,
	}

	gracehttp.SetLogger(log.New(os.Stderr, "gracehttp :: ", log.LstdFlags))
	if err := gracehttp.ServeWithOptions(
		[]*http.Server{server},
		gracehttp.PreStartProcess(func() error {
			fmt.Println("Cleaning up before restart")
			return nil
		}),
	); err != nil {
		panic(err)
	}
}

func ensureSchema(db *sql.DB) error {
	sourceInstance, err := httpfs.New(http.FS(migrations), "database/migrations")
	defer sourceInstance.Close()
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance, %w", err)
	}
	m, err := migrate.NewWithInstance(
		"httpfs",
		sourceInstance, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func home(w http.ResponseWriter, r *http.Request) {
	params := html.HomeParams{
		Title:       "Sauna Gauge",
		Temperature: 39.5,
		Humidity:    14.1,
	}

	html.Home(w, fsys, params)
}
