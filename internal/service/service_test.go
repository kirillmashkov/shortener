package service

import (
	"log"
	"math/rand/v2"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/storage/database"
	"github.com/kirillmashkov/shortener.git/internal/storage/memory"
	"go.uber.org/zap"
)

var ServerConf config.ServerConfig
var Log *zap.Logger = zap.NewNop()
var Database *database.Database
var RepositoryShortURL *database.RepositoryShortURL
var ServiceShort *Service
var Storage *memory.StoreURLMap
const originalURLPrefix = "http://www.yandex.ru/"

func Init() {
	config.InitServerConf(&ServerConf, Log)

	Storage, err := memory.New(&ServerConf, Log, &ServerConf)
	if err != nil {
		Log.Error("Error init memory storage", zap.Error(err))
		panic(err)
	}
	ServiceShort = New(Storage, ServerConf, Log)

	// Database = database.New(&ServerConf, Log)
	
	// if err := Database.Open(); err != nil {
	// 	Log.Error("Error open database", zap.Error(err))
	// 	panic(err)
	// }

	// if err := Database.Migrate(); err != nil {
	// 	Log.Error("Error migrate", zap.Error(err))
	// 	panic(err)
	// }
	// RepositoryShortURL = database.NewRepositoryShortURL(Database, Log)
	// ServiceShort = New(RepositoryShortURL, ServerConf, Log)
}

func Finish() {
	// errClose := Database.Close()
	// if errClose != nil {
	// 	Log.Error("Error close connection db", zap.Error(errClose))
	// }
}

func changeWorkingDir() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
  	if err != nil {
		Log.Error("Error change working dir", zap.Error(err))
    	panic(err)
 	}
}

func BenchmarkPostHandler(b *testing.B) {
	changeWorkingDir()
	Init()
	
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	strconv.Itoa(rand.IntN(100000))
	for i := 0; i < b.N; i++ {
		ServiceShort.ProcessURL(b.Context(), originalURLPrefix + strconv.Itoa(rand.IntN(100000)), 1)
	}

	Finish()
}

