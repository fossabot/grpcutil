/* 
Copyright 2021 Acacio Cruz acacio@acacio.coom

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grpcutil

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/davecgh/go-spew/spew"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const prefix = "Bearer "

// TokenAuth implement gRPC PerRPCCredentials interface
type TokenAuth struct {
	Token string
}

// GetRequestMetadata implements PerRPCCredentials interface
// Return value is mapped to request headers.
func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.Token,
	}, nil
}

// RequireTransportSecurity implements PerRPCCredentials interface
func (t TokenAuth) RequireTransportSecurity() bool {
	return true
}

// WithPerRPCToken - helper function
func WithPerRPCToken(token string) grpc.DialOption {
	return grpc.WithPerRPCCredentials(TokenAuth{Token: token})
}

// TokenAuthFunc - to be used with
// https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware/auth#UnaryServerInterceptor
func TokenAuthFunc(srvToken string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) { return TokenAuthCheck(ctx, srvToken) }
}

// TokenAuthCheck ...
func TokenAuthCheck(ctx context.Context, srvToken string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("RPC auth: metadata retrieval failed")
		return ctx, status.Errorf(codes.InvalidArgument, "retrieving metadata failed")
	}

	auth, ok := md["authorization"]
	if !ok {
		log.Println("RPC auth: No auth details supplied")
		return ctx, status.Errorf(codes.InvalidArgument, "no auth details supplied")
	}

	if !strings.HasPrefix(auth[0], prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}

	if strings.TrimPrefix(auth[0], prefix) != srvToken {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token "+auth[0])
	}

	log.Printf("RPC auth: %s\n", auth)
	spew.Dump(auth)
	return ctx, nil
}
