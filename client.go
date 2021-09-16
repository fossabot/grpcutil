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
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	grpc_codes "google.golang.org/grpc/codes"

	"github.com/acacio/tlsutil"
)

func setupDialOpts(tlstype, ca, crt, key string, opts []grpc.DialOption) ([]grpc.DialOption, error) {
	if tlstype != "" && tlstype != "insecure" {
		config, err := tlsutil.SetupClientTLS(tlstype, ca, crt, key)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(config)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// NO TLS
		opts = append(opts, grpc.WithInsecure())
	}
	return opts, nil
}

// ClientOpts configure gRPC client connection
type ClientOpts struct {
	TLSType string
	CA      string
	Cert    string
	Key     string
	Token   string
	Block   bool
}

// SetupConnection handles base gRPC connection establishment
func SetupConnection(addr string, opts *ClientOpts) (*grpc.ClientConn, error) {
	// println("Setting up client...")
	callopts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(grpc_codes.DeadlineExceeded, grpc_codes.Unavailable),
	}

	dialopts := []grpc.DialOption{
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(callopts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(callopts...)),
		grpc.FailOnNonTempDialError(true),
		grpc.WithPerRPCCredentials(TokenAuth{Token: "averysecuretoken, ha!"}),
		grpc.WithBlock(),
	}

	dialopts, err := setupDialOpts(opts.TLSType, opts.CA, opts.Cert, opts.Key, dialopts)
	if err != nil {
		fmt.Printf("ERROR: Failed to setup gRPC connection:\n%v\n", err)
		return nil, err
	}
	// println("Dialing:", addr)
	conn, err := grpc.Dial(addr, dialopts...)
	// defer conn.Close()

	return conn, err
}
