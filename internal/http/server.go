package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	ctx        context.Context
	idleConsCh chan struct{}
	Address    string
}

// NewServer is the function for creating new server
// Здесь мы создаём свой сервер
func NewServer(ctx context.Context, address string) *Server {
	return &Server{
		ctx:        ctx,
		idleConsCh: make(chan struct{}),
		Address:    address,
	}
}

// adder is the function for addition numbers in our URL
func adder(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	sum := 0
	for _, nStr := range str {
		n, err := strconv.Atoi(string(nStr))
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error is %v", err)
		}

		sum += n
	}
	_, _ = fmt.Fprintf(w, "Sum of all nums: %d", sum)
}

// Run is the function for running server
// А здесь этот сервер мы запускаем
func (s *Server) Run() error {

	mux := http.NewServeMux()

	//Standard handler for printing "Hello, World!"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})

	//Handler for addition 1, 2 and 3
	mux.HandleFunc("/123", adder)

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      mux,
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
	close(s.idleConsCh) // как только закрывается канал, функция WaitForGT завершается
	// и наш сервер полностью завершает работу
}

// WaitForGT is the function for waiting until ListenCtxForGT it will work
// С помощью канала функция позволяет дождаться исполнения Graceful Shutdown
func (s *Server) WaitForGT() {
	<-s.idleConsCh // блок до записи или закрытя канала
}
