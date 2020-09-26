package logger

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type LogKey string

const key LogKey = "logger"

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zap.NewExample()
		l = l.With(zap.Namespace("homatic"), zap.String("I'm", "gopher"))
		c := context.WithValue(r.Context(), key, l)
		newR := r.WithContext(c)
		next.ServeHTTP(w, newR)
	})
}

func GetLog(ctx context.Context) *zap.Logger {
	value := ctx.Value(key)
	if value == nil {
		return zap.NewExample()
	}

	logger, ok := value.(*zap.Logger)
	if !ok {
		return zap.NewExample()
	}
	return logger
}
