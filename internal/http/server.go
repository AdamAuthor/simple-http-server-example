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
	r.Post("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Categories().Create(r.Context(), category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	// Read All
	r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		filter := &models.CategoriesFilter{}

		if searchQuery := queryValues.Get("query"); searchQuery != "" {
			filter.Query = &searchQuery
		}
		categories, err := s.store.Categories().All(r.Context(), filter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, categories)
	})

	// Read by id
	r.Get("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		category, err := s.store.Categories().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, category)
	})

	// Update
	r.Put("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Categories().Update(r.Context(), category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	// Delete
	r.Delete("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}
		if err = s.store.Categories().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Post("/goods", func(w http.ResponseWriter, r *http.Request) {
		good := new(models.Good)
		if err := json.NewDecoder(r.Body).Decode(good); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Goods().Create(r.Context(), good)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	// Read All
	r.Get("/goods", func(w http.ResponseWriter, r *http.Request) {
		goods, err := s.store.Goods().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, goods)
	})

	// Read by id
	r.Get("/goods/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
		}

		good, err := s.store.Goods().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, good)
	})

	// Update
	r.Put("/goods", func(w http.ResponseWriter, r *http.Request) {
		good := new(models.Good)
		if err := json.NewDecoder(r.Body).Decode(good); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}

		err := s.store.Goods().Update(r.Context(), good)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	// Delete
	r.Delete("/goods/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Unknown error: %v", err)
			return
		}
		if err = s.store.Goods().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "DB err: %v", err)
			return
		}
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
