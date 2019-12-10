package core

import (
	"context"
	"github.com/micro/go-micro/server"
)

func EnforceVersion(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Info(ctx.Value("service_version"))
		log.Info(req)
		log.Info(rsp)
		return fn(ctx, req, rsp)
	}
}
