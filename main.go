package main

import (
	"flag"
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
	Port       string
	Production bool
}

var fsys fs.FS

func main() {
	cfg := config{
		Port:       "8080",
		Production: true,
	}

	flag.StringVar(&cfg.Port, "port", "8080", "TCP port to listen on")
	flag.BoolVar(&cfg.Production, "prod", true, "Set to false when developing for a better experience")
	flag.Parse()

	if cfg.Production {
		fsys = html.EmbedFS
	} else {
		fsys = os.DirFS("html")
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.FS(fsys)))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/temperature", htmxTemperature)
	mux.HandleFunc("/humidity", htmxHumidity)

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

func htmxTemperature(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", 40)
}

func htmxHumidity(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", 12)
}
