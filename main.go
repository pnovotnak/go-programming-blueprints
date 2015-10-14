package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)


type templateHandler struct {
	once			sync.Once
	filename	string
	templ 		*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("/Users/pnovotnak/go/src/chat/templates/", t.filename)))
	})

	t.templ.Execute(w, r)
}

func getEnv(name string, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}

func main() {
	println()
	println("    `( ◔ ౪◔)´")
	println()

	listenHost := getEnv("HOST", "localhost")
	listenPort := getEnv("PORT", "8080")

	// set up the chat room
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	println(" starting chat room...")
	println()
	go r.Run()

	bindSpec := fmt.Sprintf(fmt.Sprintf("%v:%v", listenHost, listenPort))
	println("  ...starting server:", bindSpec)
	if err := http.ListenAndServe(bindSpec, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
