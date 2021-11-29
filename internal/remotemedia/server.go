package remotemedia

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"time"

	"path/filepath"
	"runtime"

	"github.com/scriptodude/remote-media/pkg/mediahandler"
	log "github.com/sirupsen/logrus"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	root = filepath.Join(filepath.Dir(b), "../..")
)

func StartServer() {
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
	kb := mediahandler.NewKeyboardMediaHandler()

	paths := &http.ServeMux{}

	static := http.FileServer(http.Dir(path.Join(root, "web/static")))
	paths.Handle("/", static)

	paths.HandleFunc("/next", wrap(kb.PlayNext))
	paths.HandleFunc("/prev", wrap(kb.PlayPrevious))
	paths.HandleFunc("/volume-up", wrap(kb.VolumeUp))
	paths.HandleFunc("/volume-down", wrap(kb.VolumeDown))

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

func wrap(fn func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
