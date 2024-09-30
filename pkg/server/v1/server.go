package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/eadydb/nebulae/pkg/config"
	util "github.com/eadydb/nebulae/pkg/utils"
	"github.com/eadydb/nebulae/proto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	srv                  *server // waits for 1 second before forcing a server shutdown
	forceShutdownTimeout = 1 * time.Second
)

type server struct {
	proto.NebulaeServiceServer
}

func Initialize(opts config.NebulaeOptions) (func() error, error) {
	emptyCallback := func() error { return nil }
	if !opts.EnableRPC && opts.RPCPort == 0 && opts.RPCHTTPPort == 0 {
		slog.Debug("octopus API not starting as it's not requested")
		return emptyCallback, nil
	}

	grpcCallback, grpcPort, err := newGRPCServer(opts.RPCPort)
	if err != nil {
		return grpcCallback, fmt.Errorf("starting gRPC server: %w", err)
	}

	httpCallback := emptyCallback
	if opts.RPCHTTPPort > 0 {
		httpCallback, err = newHTTPServer(opts.RPCHTTPPort, grpcPort)
	}
	callback := func() error {
		httpErr := httpCallback()
		grpcErr := grpcCallback()
		errStr := ""
		if grpcErr != nil {
			errStr += fmt.Sprintf("grpc callback error: %s\n", grpcErr.Error())
		}
		if httpErr != nil {
			errStr += fmt.Sprintf("http callback error: %s\n", httpErr.Error())
		}

		return errors.New(errStr)
	}
	if err != nil {
		return callback, fmt.Errorf("starting HTTP server: %w", err)
	}

	if opts.EnableRPC && opts.RPCPort == 0 && opts.RPCHTTPPort == 0 {
		slog.Warn("started octopus gRPC API on random ", slog.Any("port", grpcPort))
	}

	return callback, nil
}

func newGRPCServer(preferredPort int) (func() error, int, error) {
	l, port, err := listenPort(preferredPort)
	if err != nil {
		return func() error { return nil }, 0, fmt.Errorf("creating listener: %w", err)
	}

	slog.Info("starting gRPC server", slog.Any("port", port))
	s := grpc.NewServer()

	srv = &server{}
	proto.RegisterNebulaeServiceServer(s, srv)

	go func() {
		if err := s.Serve(l); err != nil {
			slog.Error("failed to start grpc server ", slog.Int("port", port), slog.String("err", err.Error()))
		}
	}()
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), forceShutdownTimeout)
		defer cancel()
		ch := make(chan bool, 1)
		go func() {
			s.GracefulStop()
			ch <- true
		}()
		for {
			select {
			case <-ctx.Done():
				return l.Close()
			case <-ch:
				return l.Close()
			}
		}
	}, port, nil
}

func newHTTPServer(preferredPort, proxyPort int) (func() error, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := proto.RegisterNebulaeServiceHandlerFromEndpoint(context.Background(), mux, net.JoinHostPort(util.Loopback, strconv.Itoa(proxyPort)), opts)
	if err != nil {
		return func() error { return nil }, err
	}

	l, port, err := listenPort(preferredPort)
	if err != nil {
		return func() error { return nil }, fmt.Errorf("creating listener: %w", err)
	}

	slog.Info("starting gRPC HTTP server ", slog.Int("port", port), slog.Int("proxyPort", proxyPort))
	server := &http.Server{
		Handler: mux,
	}

	go server.Serve(l)

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), forceShutdownTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	}, nil
}

func listenPort(port int) (net.Listener, int, error) {
	l, err := net.Listen("tcp", net.JoinHostPort(util.Loopback, strconv.Itoa(port)))
	if err != nil {
		return nil, 0, err
	}
	return l, l.Addr().(*net.TCPAddr).Port, nil
}
