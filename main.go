package main

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"1mk.no/saunagauge/html"
	"github.com/facebookgo/grace/gracehttp"
)

type config struct {
	Port         string
	IsProduction bool
}

var fsys fs.FS

func main() {
	cfg := config{
		Port:         "8080",
		IsProduction: false,
	}

	if cfg.IsProduction {
		fsys = html.EmbedFS
	} else {
		fsys = os.DirFS("html")
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.FS(fsys)))

	mux.HandleFunc("/", home)

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

func home(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	params := html.HomeParams{
		Title:       "Sauna",
		Time:        fmt.Sprintf("%d-%d-%d %d:%d", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute()),
		Temperature: 39.5,
		Humidity:    14.1,
	}

	html.Home(w, fsys, params)
}
