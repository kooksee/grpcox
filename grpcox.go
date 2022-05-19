package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gusaul/grpcox/core"
	"github.com/gusaul/grpcox/handler"
	"github.com/gusaul/grpcox/web/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start app
	addr := "0.0.0.0:6969"
	muxRouter := chi.NewMux()
	var h = app.App()
	muxRouter.Handle("/", h)
	muxRouter.Mount("/", h)

	handler.Init(muxRouter)
	var wait time.Duration = time.Second * 15

	for _, r := range muxRouter.Routes() {
		fmt.Println(r.Pattern)
	}

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      middleware.Logger(muxRouter),
	}

	fmt.Println("Service started on", addr)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err := removeProtos()
	if err != nil {
		log.Printf("error while removing protos: %s", err.Error())
	}

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// removeProtos will remove all uploaded proto file
// this function will be called as the server shutdown gracefully
func removeProtos() error {
	log.Println("removing proto dir from /tmp")
	err := os.RemoveAll(core.BasePath)
	return err
}
