package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"climadash/internal/data"
	"climadash/internal/handlers"
)

func main() {
	dataPath := getenv("DATA_PATH", "data/emissions.csv")
	port := getenv("PORT", "8080")

	log.Printf("ClimaDash iniciando — dataset=%s, porta=%s", dataPath, port)

	repo, err := data.LoadFromCSV(dataPath)
	if err != nil {
		log.Fatalf("falha ao carregar dataset: %v", err)
	}
	log.Printf("dataset carregado: %d países, %d anos",
		len(repo.Countries()), len(repo.Years()))

	mux := http.NewServeMux()
	handlers.New(repo).Register(mux)
	mux.Handle("/", http.FileServer(http.Dir("web/static")))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      logging(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("escutando em http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown com erro: %v", err)
	}
	log.Printf("servidor encerrado")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
