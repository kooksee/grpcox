package main

import (
	"context"
	"fmt"
	"github.com/fullstorydev/grpchan"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/gusaul/grpcox/internal/proto/demov1pb"
	"github.com/gusaul/grpcox/svc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gusaul/grpcox/handler"
	"github.com/gusaul/grpcox/web/ui"
	//	https://github.com/mlctrez/goapp-pf
	_ "github.com/fullstorydev/grpchan/httpgrpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start app
	addr := "0.0.0.0:6969"
	muxRouter := chi.NewMux()

	var uiHandle = app.App()
	muxRouter.Handle("/", uiHandle)
	muxRouter.Mount("/", uiHandle)

	handler.Init(muxRouter)

	handlers := make(grpchan.HandlerMap)
	demov1pb.RegisterTransportServer(handlers, svc.NewServer())

	httpgrpc.HandleServices(func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		fmt.Println(pattern)
		muxRouter.Options(pattern, func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
			writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			writer.WriteHeader(http.StatusOK)
		})
		muxRouter.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization")
			writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			handler(writer, request)
		})
	}, "/grpc", handlers, nil, nil)

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
