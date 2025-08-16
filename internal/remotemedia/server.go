package remotemedia

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"time"

	"path/filepath"
	"runtime"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	root = filepath.Join(filepath.Dir(b), "../..")

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
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
	paths := &http.ServeMux{}

	static := http.FileServer(http.Dir(path.Join(root, "web/release")))
	paths.Handle("/", static)

	paths.HandleFunc("/ws", handleWebSocket)

	var handler http.Handler = paths
	return logHandler(handler)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	header := http.Header{}
	header.Add("Access-Control-Allow-Origin", "*")
	log.Info("Attempting to upgrade user")

	c, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	log.Info("Upgrade successful")
	messageBroker := NewMessageBroker()

	c.WriteMessage(0x1, messageBroker.GetCurrentState())

	for {
		opcode, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		response := messageBroker.HandleMessage(opcode, message)
		if response != nil {
			c.WriteMessage(0x1, response)
		}
	}
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
