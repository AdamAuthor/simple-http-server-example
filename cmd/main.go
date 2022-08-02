package main

import (
	"context"
	"fmt"
	"kbtu_go_6/internal/http"
	"kbtu_go_6/internal/store/inmemory"
)

func main() {
	// контекст с таймаутом
	// отправляем этот контекст в сервер.
	// Спустя время контекст взрывается, и происходит Graceful Shutdown.
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()

	// создаём базу с ноутбуками
	store := inmemory.NewDB()

	//Creating and run new server
	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

	srv.WaitForGT()
}
