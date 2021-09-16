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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CheckRPCAuth - used by servers to check RPC auth
func CheckRPCAuth(ctx context.Context) error {
	// code from the authorize() function:
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("RPC auth: metadata retrieval failed")
		return status.Errorf(codes.InvalidArgument, "retrieving metadata failed")
	}

	elem, ok := md["authorization"]
	if !ok {
		log.Println("RPC auth: No auth details supplied")
		return status.Errorf(codes.InvalidArgument, "no auth details supplied")
	}

	log.Printf("RPC auth: %s\n", elem)
	return nil
}
