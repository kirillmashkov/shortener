package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
	"github.com/kirillmashkov/shortener.git/internal/logger"
	"github.com/kirillmashkov/shortener.git/internal/model"

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
	// cancel will be called explicitly on server shutdown

	flag.Parse()
	err = app.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	defer app.Close()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	model.Wg.Add(1)
	go runServer(sigint, cancel)

	if model.ShortURLchan != nil {
		close(model.ShortURLchan)
	}

	if model.Wg != nil {
		model.Wg.Wait()
	}
}

func runServer(sigint chan os.Signal, cancel context.CancelFunc) {
	server := &http.Server{
		Addr:    app.ServerConf.Host,
		Handler: router.Serv(),
	}

	go func() {
		<-sigint
		if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
			app.Log.Error("error shutdown server")
		} else {
			app.Log.Info("server shutdown graceful")
		}
		// After HTTP server shutdown, cancel the main context:
		cancel()
		model.Wg.Done()
	}()

	var err error

	if app.ServerConf.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}

		server.TLSConfig = manager.TLSConfig()
		if err = server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			app.Log.Error("HTTPS server ListenAndServeTLS", zap.Error(err))
		}
	} else {
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			app.Log.Error("HTTPS server ListenAnd", zap.Error(err))
		}
	}

}
