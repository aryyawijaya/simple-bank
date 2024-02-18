package utilgrpc

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func ExtractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtdt
	}

	// set UserAgent Gateway client
	if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
		mtdt.UserAgent = userAgents[0]
	}

	// set UserAgent gRPC client
	if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
		mtdt.UserAgent = userAgents[0]
	}

	// set ClientIP Gateway client
	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	// set ClientIP gRPC client
	p, ok := peer.FromContext(ctx)
	if !ok {
		return mtdt
	}
	mtdt.ClientIP = p.Addr.String()

	return mtdt
}
