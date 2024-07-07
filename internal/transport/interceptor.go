package transport

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func NewApiKeyInterceptor(apiKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, GRPCError(codes.Unauthenticated, 10001, "Unable to get metadata")
		}
		keys := md.Get("x-api-key")
		if len(keys) == 0 {
			return nil, GRPCError(codes.Unauthenticated, 10001, "Unauthorized: API Key not found")
		}
		if keys[0] != apiKey {
			return nil, GRPCError(codes.Unauthenticated, 10001, "Unauthorized: API Key invalid")
		}
		return handler(ctx, req)
	}
}
