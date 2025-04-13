package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/kirillmashkov/shortener.git/internal/app"
    "github.com/kirillmashkov/shortener.git/internal/httpserver/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func Initialize() {
	encoderCfg := zap.NewProductionEncoderConfig()
    encoderCfg.TimeKey = "timestamp"
    encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

    config := zap.Config{
        Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
        Development:       false,
        DisableCaller:     false,
        DisableStacktrace: false,
        Sampling:          nil,
        Encoding:          "json",
        EncoderConfig:     encoderCfg,
        OutputPaths: []string{
            "stdout",
        },
        ErrorOutputPaths: []string{
            "stderr",
        },
        InitialFields: map[string]interface{}{
            "pid": os.Getpid(),
        },
    }

	app.Log = zap.Must(config.Build())
    
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		app.Log.Info("request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )

		ww := middleware.Wrap(w)

        next.ServeHTTP(ww, r)		

		app.Log.Info("response",
			zap.Duration("elapsed", time.Since(t1)),
            zap.Int("status", ww.Status()),
            zap.Int("content length", ww.Bytes()))
			
	})
}	