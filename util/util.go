package util

import (
	"context"
	"net/http"
	"os"
)

type CtxUserName string

const (
	UserName = "userName"
)

func GetCtxWithUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, CtxUserName(UserName), userName)
}

func GetUserNameFromContext(ctx context.Context) string {
	s, ok := ctx.Value(CtxUserName(UserName)).(string)
	if ok {
		return s
	}
	return ""
}

func CORSCheck(r *http.Request) bool {
	switch r.Header.Get("Origin") {
	case "http://localhost:" + GetPortEnv():
		return true
	default:
		return false
	}
}
func GetPortEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}
