package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"../../api"
	"../../controller"
	"../../maxmind"
)

// NetworkTimeoutSeconds is the number of seconds to wait until we timeout a network requests
const NetworkTimeoutSeconds = 15

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(NetworkTimeoutSeconds * time.Second))

	// Create dependencies
	maxmindSvc := maxmind.New(maxmind.NewParams{
		LicenseKey: os.Getenv("MAXMIND_LICENSE_KEY"),
	})
	if err := maxmindSvc.GenerateDB(); err != nil {
		log.Fatal(err)
	}

	controllerSvc := controller.New(controller.NewParams{
		GeoDB: maxmindSvc,
	})

	router := api.New(api.NewParams{
		ControllerSvc: controllerSvc,
	})

	// Route to use tell an orchestrator that we're ready to receive traffic
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("good to go!"))
	})

	r.Route("/", router.Route)

	http.ListenAndServe(":3000", r)
}
