package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"
	"runtime"

	"enron_corp_indexer.snowartz/internal/config"
	"enron_corp_indexer.snowartz/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var port string
	flag.StringVar(&port, "port", "8080", "Puerto del servidor HTTP")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Inicializa el router Chi.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Registrar las rutas de la API.
	apiHandler := handlers.APIHandler{Config: cfg}
	apiHandler.RegisterRoutes(r)

	// Agregar las rutas de debug/pprof
	r.Mount("/debug", middleware.Profiler())

	// Configura y arrancar el servidor HTTP.
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
