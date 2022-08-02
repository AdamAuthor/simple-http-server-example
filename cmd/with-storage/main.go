package main

import (
	"context"
	"fmt"
	"kbtu_go_6/internal/http"
)

func main() {
	// контекст с таймаутом
	// отправляем этот контекст в сервер.
	// Спустя время контекст взрывается, и происходит Graceful Shutdown.

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()

	//Creating and run new server
	srv := http.NewServer(context.Background(), ":8080")
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

	srv.WaitForGT()
}
