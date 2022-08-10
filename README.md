# grpcutil
gRPC Utility functions


### Example use

```go
import (
  "github.com/acacio/grpcutil"
)

func NewMyServer(port, grpcwebport string) *MyServer {

  logger := grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stderr)
  grpclog.SetLoggerV2(logger)

  // Create gRPC server
  srvOpts := grpcutil.ServerOptions()
  s := grpc.NewServer(srvOpts...)
  log.Println("Registering gRPC reflection...")
  reflection.Register(s)

  mysrv := &MyServer{}
  pb.RegisterMyServer(s, mysrv)

  // After your registrations, Prometheus metrics are initialized.
  log.Println("Registering Prometheus...")
  grpc_prometheus.Register(s)

  // Enable gRPCWeb serving. Requires goroutine
  go grpcutil.StartgRPCWeb(s, grpcwebport)

  // Start listening & serving requests
  grpcutil.Serve(s, port) // Fatal failure on error

  return mysrv
}


```
