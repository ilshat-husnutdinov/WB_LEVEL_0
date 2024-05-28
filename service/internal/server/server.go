package server

import (
	"log"
	"net/http"
	"service/service/internal/config"
	handlers "service/service/internal/handlers/get"

	"service/service/pkg/cache"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	srv *http.Server
	csh *cache.Cache
}

// CreateServer создаёт http-сервер с роутером
func CreateServer(serverCfg *config.HTTPServer, csh *cache.Cache) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/orders/{id}", handlers.GetOrderHandler(csh))

	return &Server{
		srv: &http.Server{
			Addr:         serverCfg.Address,
			Handler:      router,
			ReadTimeout:  serverCfg.Timeout,
			WriteTimeout: serverCfg.Timeout,
			IdleTimeout:  serverCfg.IdleTimeout,
		},
		csh: csh,
	}
}

// RunServer запускает сервер
func RunServer(srv *Server) error {
	if err := srv.srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server")
	}
	return nil
}
