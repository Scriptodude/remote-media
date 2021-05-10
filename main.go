package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/micmonay/keybd_event"
	log "github.com/sirupsen/logrus"
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
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	paths := &http.ServeMux{}

	static := http.FileServer(http.Dir("./static"))
	paths.Handle("/", static)

	paths.HandleFunc("/next", change(keybd_event.VK_NEXTSONG, kb))
	paths.HandleFunc("/prev", change(keybd_event.VK_PREVIOUSSONG, kb))

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

func change(direction int, kb keybd_event.KeyBonding) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var word string = "next"

		if direction == keybd_event.VK_PREVIOUSSONG {
			word = "previous"
		}

		log.Info(fmt.Sprintf("Playing %s song", word))
		kb.SetKeys(direction)
		kb.Press()
		kb.Release()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
