package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	authorizationHeader        = "authorization"
	authorizationBearer        = "bearer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}

	if data, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := data.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			md.UserAgent = userAgent[0]
		}

		if userAgent := data.Get(userAgentHeader); len(userAgent) > 0 {
			md.UserAgent = userAgent[0]
		}

		if clientIP := data.Get(xForwardedForHeader); len(clientIP) > 0 {
			md.ClientIP = clientIP[0]
		}
	}

	// karna dari metadata ga dapet ip client jika hit pake grpc, kita pake peer
	if peer, ok := peer.FromContext(ctx); ok {
		if addr := peer.Addr.String(); addr != "" {
			md.ClientIP = addr
		}
	}

	return md
}
