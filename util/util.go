package util

import (
	"context"
	"net/http"
)

type CtxUserName string

const (
	UserName = "userName"
	APP_PORT = ":8080" // todo use ENV variable
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
	case "http://localhost" + APP_PORT:
		return true
	default:
		return false
	}
}
