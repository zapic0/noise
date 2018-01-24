package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gorilla/mux"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "port number to listen")
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/play/{sound}", PlayHandler).Methods("Get")

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      logRequest(router),
	}

	// Run our server in a goroutine so that it doesn't block.
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

}

// PlayHandler plays a sound
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	f, _ := os.Open(vars["sound"] + ".mp3")

	s, format, _ := mp3.Decode(f)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(s)
}

// Middleware-like function to be able to log incoming requests
// https://gist.github.com/hoitomt/c0663af8c9443f2a8294#file-log_request-go-L33
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
