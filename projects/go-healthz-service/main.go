// a mini project i did on webserver checks

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// this is the starting point of my mini program
// it runs the healthcheck switch and decides whether to run as a server or to run as a self check
func main() {
	healthcheck := flag.Bool("healthcheck", false, "check /healthz on the local server and exit")
	flag.Parse()

	// this part reads the door number from the PORT setting
	// if no door number is set it falls back to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// if the program in the healthcheck above started with as a self check
	// it will check itself and exist. The server will not be started
	if *healthcheck {
		os.Exit(runHealthcheck(port))
	}

	if err := serve(port); err != nil {
		log.Fatal(err)
	}
}

// this part is a router, or multiplexer. if you choose to call it so
// matches HTTP requests to the corresponding handler function based onr URL path and the HTTP method
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleRoot)        // registers a question that goes into handlerood via /
	mux.HandleFunc("GET /healthz", handleHealthz) // this registers a question that goes into handleHealthz via /healthz
	return mux
} // the code here only questions if the requests were/are accepted

// answers  the "/" handler
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "hello from go-healthz-service")
}

// answers  the "/healthz" handler
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}

// this builds the server. It listens to every entrance, uses the router and writes four timeouts
func serve(port string) error {
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,  // to receive the request's headers
		ReadTimeout:       10 * time.Second, // to read the entire request
		WriteTimeout:      10 * time.Second, // to write the answer
		IdleTimeout:       60 * time.Second, // set for an idle connection
	}

	// this part waits for a stop signal as the server runs, when it arrives, it gives a 5s notice to finish and exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on :%s", port)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	log.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}

// perfoms a self-check, asks its own /healthz with a 2-second limit.
// retuns 0 when healthy and1 when unhealthy
func runHealthcheck(port string) int {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/healthz")
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck failed:", err)
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck failed: status", resp.StatusCode)
		return 1
	}
	return 0
}

// end of program. written by a learning and practicing student, Nathan John, KTÜ, TR.
