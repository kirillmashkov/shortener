package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
	"github.com/kirillmashkov/shortener.git/internal/logger"
	"github.com/kirillmashkov/shortener.git/internal/model"

	_ "net/http/pprof" 
)

func main() {
	err := logger.Initialize()
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("ERROR: ")
		log.Println("Can't init logger")
		panic(err)
	}

	flag.Parse()
	err = app.Initialize()
	if err != nil {
		panic(err)
	}

	defer app.Close()

	err = http.ListenAndServe(app.ServerConf.Host, router.Serv())
	if err != nil {
		panic(err)
	}

	if model.ShortURLchan != nil {
		close(model.ShortURLchan)
	}

	if model.Wg != nil {
		model.Wg.Wait()
	}
}
