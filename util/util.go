package util

import "context"

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
