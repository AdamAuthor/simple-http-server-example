package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"kbtu_go_6/internal/models"
	"kbtu_go_6/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	ctx        context.Context
	idleConsCh chan struct{}

	// не самая лучшая практика. Обычно делается на уровне 3 слоев:
	// бизнес логика, HTTP хэндлеры, база данных (в рамка курса ок)
	store   store.Store
	Address string
}

// NewServer is the function for creating new server
// Здесь мы создаём свой сервер
func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:        ctx,
		idleConsCh: make(chan struct{}),
		store:      store,

		Address: address,
	}
}

// basicHandler был создан для инкапсуляции логики настройки мультиплексера
// К тому же, вместо использования мультиплексера, используется роутер
func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	// Create
	r.Post("/laptops", func(w http.ResponseWriter, r *http.Request) {
		laptop := new(models.Laptop)
		if err := json.NewDecoder(r.Body).Decode(laptop); err != nil {
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Create(r.Context(), laptop)
		if err != nil {
			return
		}
	})

	// Read All
	r.Get("/laptops", func(w http.ResponseWriter, r *http.Request) {
		laptops, err := s.store.All(r.Context())
		if err != nil {
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		render.JSON(w, r, laptops)
	})

	// Read by id
	r.Get("/laptops/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		laptop, err := s.store.ByID(r.Context(), id)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		render.JSON(w, r, laptop)
	})

	// Update
	r.Put("/laptops", func(w http.ResponseWriter, r *http.Request) {
		laptop := new(models.Laptop)
		if err := json.NewDecoder(r.Body).Decode(laptop); err != nil {
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Update(r.Context(), laptop)
		if err != nil {
			return
		}
	})

	// Delete
	r.Delete("/laptops/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}
		_ = s.store.Delete(r.Context(), id)
	})

	return r
}

// Run is the function for running server
func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

// ListenCtxForGT is the function for Graceful Shutdown
// При запуске сервера мы также запускаем горутину, которая дожидается своего часа
// и как только контекст будет завершён, происходит Shutdown
func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // Blocked until the application context is canceled

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("[HTTP] Got err while shutting down:", err)
	}

	log.Println("[HTTP] Processed all idle connections")
	close(s.idleConsCh)
	// как только закрывается канал, функция WaitForGT завершается
	// и наш сервер полностью завершает работу
}

// WaitForGT is the function for waiting until ListenCtxForGT it will work
// С помощью канала функция позволяет дождаться исполнения Graceful Shutdown
func (s *Server) WaitForGT() {
	<-s.idleConsCh // блок до записи или закрытя канала
}
