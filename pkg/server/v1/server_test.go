package v1

import (
	"fmt"
	"net"
	"testing"

	"github.com/eadydb/nebulae/pkg/config"
	"github.com/eadydb/nebulae/pkg/testutil"
	"github.com/eadydb/nebulae/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rpcAddr  = 12345
	httpAddr = 23456
)

func TestServerStartUp(t *testing.T) {
	// start up servers
	shutdown, err := Initialize(config.NebulaeOptions{
		EnableRPC:   true,
		RPCPort:     rpcAddr,
		RPCHTTPPort: httpAddr,
	})
	defer shutdown()
	testutil.CheckError(t, false, err)

	// create gRPC client and ensure we can connect
	conn, err := grpc.Dial(fmt.Sprintf(":%d", rpcAddr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("unable to establish skaffold grpc connection")
	}
	defer conn.Close()

	client := proto.NewNebulaeServiceClient(conn)
	if client == nil {
		t.Errorf("unable to connect to gRPC server")
	}

	// dial httpAddr and make sure port is being served on
	httpConn, err := net.Dial("tcp", fmt.Sprintf(":%d", httpAddr))
	if err != nil {
		t.Errorf("unable to connect to gRPC HTTP server")
	}
	if httpConn == nil {
		t.Errorf("unable to connect to gRPC HTTP server")
	} else {
		httpConn.Close()
	}
}
