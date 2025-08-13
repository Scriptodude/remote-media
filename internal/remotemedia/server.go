package remotemedia

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"path/filepath"
	"runtime"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/scriptodude/remote-media/internal/mediahandler"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	root = filepath.Join(filepath.Dir(b), "../..")

	upgrader = websocket.Upgrader{}
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
	// kb := mediahandler.NewKeyboardMediaHandler()
	// kb := mediahandler.NewDbusMediaHandler()

	paths := &http.ServeMux{}

	static := http.FileServer(http.Dir(path.Join(root, "web/static")))
	paths.Handle("/", static)

	paths.HandleFunc("/ws", handleWebSocket)

	// Configure the web server
	var handler http.Handler = paths
	return logHandler(handler)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	mediaHandler := mediahandler.NewKeyboardMediaHandler()
	for {
		opcode, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %d, %s", opcode, message)

		// Opcode for text
		if opcode != 0x1 {
			log.Println("Not supporting message type, continuing")
			continue
		}

		switch strings.TrimSpace(string(message)) {
		case "volume_up":
			c.WriteMessage(1, []byte(strconv.Itoa(mediaHandler.VolumeUp())))

		case "volume_down":
			c.WriteMessage(1, []byte(strconv.Itoa(mediaHandler.VolumeDown())))

		case "play_next":
			mediaHandler.PlayNext()

		case "play_previous":
			mediaHandler.PlayPrevious()

		default:
			log.Printf("Unkown command: %s", message)
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
