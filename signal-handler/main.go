package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	returnCode := make(chan int)
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
		Addr:         ":8081",
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			returnCode <- 1
		} else {
			returnCode <- 0
		}
	}()

	for {
		select {
		case s := <-gracefulStop:
			log.Printf("SIGNAL %d received, then shutting down...\n", s)
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					log.Println(err)
					returnCode <- 1
				}
			}()
		case c := <-returnCode:
			os.Exit(c)
		}
	}
}
