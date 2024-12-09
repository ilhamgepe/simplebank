package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/ilhamgepe/simplebank/token"
	"google.golang.org/grpc/metadata"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	token := md.Get(authorizationHeader)
	if len(token) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	fields := strings.Fields(token[0])
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	if strings.ToLower(fields[0]) != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", fields[0])
	}

	payload, err := s.tokenMaker.VerifyToken(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return payload, nil
}
