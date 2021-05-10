package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

const (
	PLAY_NEXT = "audio_next"
	PLAY_PREV = "audio_prev"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	port := flag.String("port", "8080", "The port to run the webservlet on, defaults to 8080")
	flag.Parse()

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", *port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      configurePaths(),
	}

	log.Info("Now listening on " + s.Addr)
	log.Fatal(s.ListenAndServe())
}

func configurePaths() http.Handler {
	paths := &http.ServeMux{}

	static := http.FileServer(http.Dir("./static"))
	paths.Handle("/", static)

	paths.HandleFunc("/next", change(PLAY_NEXT))
	paths.HandleFunc("/prev", change(PLAY_PREV))

	// Configure the web server
	var handler http.Handler = paths
	return logHandler(handler)
}

func logHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method
		log.Info(fmt.Sprintf("%s %s", uri, method))
	}

	return http.HandlerFunc(fn)
}

func change(direction string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Changing song : " + direction)
		robotgo.KeyTap(direction)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
