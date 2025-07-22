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
	"github.com/kirillmashkov/shortener.git/internal/storage/memory"
	"go.uber.org/zap"
)

const originalURLPrefix = "http://www.yandex.ru/"

func Init(log *zap.Logger, b *testing.B) (*Service, error) {
	b.Helper()

	var ServerConf config.ServerConfig

	config.InitServerConf(&ServerConf, log)

	Storage, err := memory.New(&ServerConf, log, &ServerConf)
	if err != nil {
		return nil, err
	}
	return New(Storage, ServerConf, log), nil
}

func changeWorkingDir(log *zap.Logger, b *testing.B) error {
	b.Helper()

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	return err
}

func BenchmarkPostHandler(b *testing.B) {
	var logger = zap.NewNop()

	if err := changeWorkingDir(logger, b); err != nil {
		b.Error("Error change working dir", err)
	}

	ServiceShort, err := Init(logger, b)
	if err != nil {
		b.Error("Error init memory storage", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	strconv.Itoa(rand.IntN(100000))
	for i := 0; i < b.N; i++ {
		ServiceShort.ProcessURL(b.Context(), originalURLPrefix+strconv.Itoa(rand.IntN(100000)), 1)
	}
}
