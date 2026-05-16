package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"todo-go/internal/db"
	"todo-go/internal/handlers"
	"todo-go/internal/middleware"
	"todo-go/internal/ports"
	"todo-go/internal/status"
)

func main() {
	db.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", handlers.TodosHandler)
	mux.HandleFunc("/todos/", handlers.TodoHandler)

	loggedMux := middleware.Logging(mux)

	defaultPort := 9090
	actualPort := ports.ConnectToPort(defaultPort, loggedMux)

	serverStatusFile, err := status.CreateStatFile(strconv.Itoa(actualPort))
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			if err := status.UpdateServerStatus(serverStatusFile); err != nil {
				panic(err)
			}

		case <-sig:
			fmt.Println("shutdown signal received")

			_ = serverStatusFile.Sync()
			_ = serverStatusFile.Close()
			_ = os.Remove(serverStatusFile.Name())

			if err := db.Close(); err != nil {
				fmt.Println("db close error: ", err)
			}

			break loop
		}
	}

	fmt.Println("clean shutdown complete")
}

