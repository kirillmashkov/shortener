package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver"
	"github.com/kirillmashkov/shortener.git/internal/logger"
	"github.com/kirillmashkov/shortener.git/internal/model"
	pb "github.com/kirillmashkov/shortener.git/internal/proto"
	"github.com/kirillmashkov/shortener.git/internal/server"

	_ "net/http/pprof"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	err := logger.Initialize()
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("ERROR: ")
		log.Println("Can't init logger")
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	flag.Parse()
	err = app.Initialize(ctx)
	if err != nil {
		app.Log.Error("can't init app", zap.Error(err))
		panic(err)
	}

	defer app.Close()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	var restServer server.Server

	if app.ServerConf.EnableHTTPS {
		restServer = httpserver.NewHTTPS(app.ServerConf.Host)
	} else {
		restServer = httpserver.NewHTTP(app.ServerConf.Host)
	}

	grpcServer := pb.New(*app.Service, app.ServerConf.GRPCAddress)

	model.Wg.Add(2)
	go runServer(restServer, sigint, cancel)
	go runServer(grpcServer, sigint, cancel)

	if model.Wg != nil {
		model.Wg.Wait()
	}

	if model.ShortURLchan != nil {
		close(model.ShortURLchan)
	}

}

func runServer(server server.Server, sigint chan os.Signal, cancel context.CancelFunc) {
	go func() {
		<-sigint
		if errShutdown := server.Shutdown(); errShutdown != nil {
			app.Log.Error("error shutdown server")
		} else {
			app.Log.Info("server shutdown graceful")
		}

		cancel()
		model.Wg.Done()
	}()

	if err := server.Run(); err != nil {
		app.Log.Error("can't run rest server", zap.Error(err))
	}
}
